package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"neosim_go/internal/modules/auth/contracts"
	"neosim_go/internal/modules/auth/dto"
	"neosim_go/internal/modules/auth/models"
	rbacContracts "neosim_go/internal/modules/rbac/contracts"
	rbacDto "neosim_go/internal/modules/rbac/dto"
	userContracts "neosim_go/internal/modules/users/contracts"
	userDto "neosim_go/internal/modules/users/dto"
	userModels "neosim_go/internal/modules/users/models"
	"neosim_go/internal/shared/types"
	"neosim_go/internal/shared/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ─── Config ────────────────────────────────────────────────────────────────────

type AuthServiceConfig struct {
	JWTManager                   *utils.JWTManager
	LoginMaxAttempts             int
	LoginLockDurationMinutes     int
	MaxConcurrentSessions        int
	RateLimitLoginPerIPPerMinute int
	PasswordPolicy               *utils.PasswordPolicy
	PasswordHistoryCount         int
	IsRegistrationActive         bool
	AutoActiveUser               bool
	MailResetTokenExpMinutes     int
	AppFrontendURL               string
	Mailer                       *utils.Mailer
}

// ─── Init ──────────────────────────────────────────────────────────────────────
// ─── Init ──────────────────────────────────────────────────────────────────────
type authService struct {
	repo     contracts.AuthRepository
	userRepo userContracts.Repository
	rbacRepo rbacContracts.RBACRepository // ← Tambahkan ini
	redis    *redis.Client
	cfg      AuthServiceConfig
}

func NewAuthService(
	repo contracts.AuthRepository,
	userRepo userContracts.Repository,
	rbacRepo rbacContracts.RBACRepository, // ← Tambahkan ini
	redisClient *redis.Client,
	cfg AuthServiceConfig,
) contracts.AuthService {
	return &authService{
		repo:     repo,
		userRepo: userRepo,
		rbacRepo: rbacRepo, // ← Set di sini
		redis:    redisClient,
		cfg:      cfg,
	}
}

// ─── End Init ──────────────────────────────────────────────────────────────────

// ─── Login ─────────────────────────────────────────────────────────────────────

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, ip, userAgent string) (*dto.TokenResponse, error) {
	ctxs := context.Background()
	identifier := strings.TrimSpace(req.Identifier)

	// 1. Cari user by username atau email
	user, err := s.findUserByIdentifier(ctx, identifier)
	if err != nil || user == nil {
		s.saveLoginHistory(ctx, nil, identifier, ip, userAgent, models.LoginStatusFailed, "identifier_not_found")
		return nil, NewAuthError(401, "Username atau email tidak terdaftar.")
	}

	// 2. Cek is_active
	if !user.IsActive {
		s.saveLoginHistory(ctx, &user.ID, identifier, ip, userAgent, models.LoginStatusFailed, "account_inactive")
		return nil, NewAuthError(403, "Akun Anda tidak aktif. Hubungi administrator.")
	}

	// 3. Cek IP blacklist
	if exists, _ := s.redis.Exists(ctxs, utils.KeyIPBlacklist(ip)).Result(); exists > 0 {
		return nil, NewAuthError(429, fmt.Sprintf("IP Anda (%s) telah diblokir sementara karena terlalu banyak percobaan login.", ip))
	}

	// 4. Cek login lock
	if ttl, err := s.redis.TTL(ctxs, utils.KeyLoginLock(identifier)).Result(); err == nil && ttl > 0 {
		minutesLeft := int(math.Ceil(ttl.Minutes()))
		return nil, NewAuthError(429, fmt.Sprintf(
			"Terlalu banyak percobaan login. Coba lagi dalam %d menit.", minutesLeft,
		))
	}

	// 5. Verifikasi password
	if !utils.VerifyPassword(req.Password, user.Password) {
		s.handleFailedLogin(ctxs, identifier, ip, &user.ID, userAgent)
		remaining := s.getRemainingAttempts(ctxs, identifier)
		return nil, NewAuthError(401, fmt.Sprintf("Password salah. Tersisa %d percobaan.", remaining))
	}

	// 6. Cek concurrent sessions
	activeCount, err := s.repo.CountActiveTokens(ctx, user.ID)
	if err != nil {
		return nil, NewAuthError(500, "Terjadi kesalahan sistem.")
	}
	if activeCount >= int64(s.cfg.MaxConcurrentSessions) {
		return nil, NewAuthError(403, fmt.Sprintf("Batas sesi aktif tercapai (max: %s device). Silakan logout dari perangkat lain terlebih dahulu.", strconv.Itoa(s.cfg.MaxConcurrentSessions)))
	}

	// 7. Generate tokens
	accessToken, err := s.cfg.JWTManager.GenerateAccessToken(
		user.ID, user.Username, user.IsSuperadmin, user.IsStaff,
	)
	if err != nil {
		return nil, NewAuthError(500, "Gagal membuat token.")
	}

	refreshToken, jti, expiresAt, err := s.cfg.JWTManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, NewAuthError(500, "Gagal membuat refresh token.")
	}

	deviceInfo := userAgent
	ipCopy := ip
	if err := s.repo.SaveToken(ctx, &models.AuthToken{
		UserID:      user.ID,
		JTI:         jti,
		TokenType:   "refresh",
		DeviceInfo:  &deviceInfo,
		IPAddress:   &ipCopy,
		IsBlacklist: false,
		ExpiresAt:   expiresAt,
	}); err != nil {
		return nil, NewAuthError(500, "Gagal menyimpan token.")
	}

	// 8. Cleanup & update
	s.redis.Del(ctxs, utils.KeyLoginAttempts(identifier))
	s.updateLastLogin(ctx, user)
	s.saveLoginHistory(ctx, &user.ID, identifier, ip, userAgent, models.LoginStatusSuccess, "")

	return s.buildTokenResponse(ctx, accessToken, refreshToken, user), nil
}

// ─── Register ──────────────────────────────────────────────────────────────────

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// 1. Cek setting registrasi
	if !s.cfg.IsRegistrationActive {
		return nil, NewAuthError(403, "Registrasi sedang tidak tersedia.")
	}

	// 2. Validasi password policy
	if errs := s.cfg.PasswordPolicy.Validate(req.Password); len(errs) > 0 {
		return nil, NewAuthError(422, strings.Join(errs, ", "))
	}

	// 3. Cek duplikat bersamaan via users repository
	var fieldErrors []string
	if existing, _ := s.userRepo.GetByUsername(ctx, req.Username); existing != nil {
		fieldErrors = append(fieldErrors, "username sudah digunakan")
	}
	if existing, _ := s.userRepo.GetByEmail(ctx, req.Email); existing != nil {
		fieldErrors = append(fieldErrors, "email sudah digunakan")
	}
	if len(fieldErrors) > 0 {
		return nil, NewAuthError(422, strings.Join(fieldErrors, ", "))
	}

	// 4. Tentukan status berdasarkan setting
	isActive := s.cfg.AutoActiveUser
	isVerified := s.cfg.AutoActiveUser

	// 5. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, NewAuthError(500, "Gagal memproses password.")
	}

	// 6. Buat user via users repository
	name := ""
	if req.Name != nil {
		name = *req.Name
	}

	user := &userModels.User{
		Username:   req.Username,
		Email:      req.Email,
		Name:       name,
		Password:   hashedPassword,
		IsActive:   isActive,
		IsVerified: isVerified,
	}

	// Set default settings dari users dto
	if err := user.SetSettings(userDto.DefaultUserSettings()); err != nil {
		return nil, NewAuthError(500, "Gagal menyimpan settings default.")
	}

	// Simpan via users repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, NewAuthError(500, "Gagal membuat akun.")
	}

	// 7. Simpan password history
	s.repo.SavePasswordHistory(ctx, &models.PasswordHistory{
		UserID:       user.ID,
		PasswordHash: hashedPassword,
	})

	// 8. Kirim email verifikasi jika perlu
	if !isVerified && s.cfg.Mailer != nil {
		// TODO: generate verification token dan kirim email
	}

	return &dto.RegisterResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}, nil
}

// ─── Refresh Token ─────────────────────────────────────────────────────────────

func (s *authService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// 1. Parse & validasi JWT
	claims, err := s.cfg.JWTManager.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, NewAuthError(401, "Refresh token tidak valid.")
	}
	if claims.TokenType != "refresh" {
		return nil, NewAuthError(401, "Token bukan refresh token.")
	}

	// 2. Cek JTI di database
	authToken, err := s.repo.GetTokenByJTI(ctx, claims.ID)
	if err != nil || authToken == nil {
		return nil, NewAuthError(401, "Token tidak ditemukan.")
	}
	if authToken.IsBlacklist {
		return nil, NewAuthError(401, "Token sudah tidak berlaku.")
	}
	if authToken.IsExpired() {
		return nil, NewAuthError(401, "Token sudah kedaluwarsa.")
	}

	// 3. Ambil user via users repository
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return nil, NewAuthError(401, "User tidak ditemukan.")
	}

	// 4. Token rotation: blacklist lama, buat baru
	if err := s.repo.BlacklistToken(ctx, claims.ID); err != nil {
		return nil, NewAuthError(500, "Gagal memproses token.")
	}

	accessToken, err := s.cfg.JWTManager.GenerateAccessToken(
		user.ID, user.Username, user.IsSuperadmin, user.IsStaff,
	)
	if err != nil {
		return nil, NewAuthError(500, "Gagal membuat token baru.")
	}

	refreshToken, jti, expiresAt, err := s.cfg.JWTManager.GenerateRefreshToken(user.ID, user.Username)
	if err != nil {
		return nil, NewAuthError(500, "Gagal membuat refresh token baru.")
	}

	deviceInfo := ""
	if err := s.repo.SaveToken(ctx, &models.AuthToken{
		UserID:      user.ID,
		JTI:         jti,
		TokenType:   "refresh",
		DeviceInfo:  &deviceInfo,
		IsBlacklist: false,
		ExpiresAt:   expiresAt,
	}); err != nil {
		return nil, NewAuthError(500, "Gagal menyimpan token baru.")
	}

	return s.buildTokenResponse(ctx, accessToken, refreshToken, user), nil
}

// ─── Forgot Password ───────────────────────────────────────────────────────────

func (s *authService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	ctxs := context.Background()

	// Selalu return nil — jangan bocorkan apakah identifier terdaftar
	user, _ := s.findUserByIdentifier(ctx, req.Identifier)
	if user == nil {
		return nil
	}

	// Generate crypto-safe token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil
	}
	token := hex.EncodeToString(tokenBytes)

	ttl := time.Duration(s.cfg.MailResetTokenExpMinutes) * time.Minute
	s.redis.Set(ctxs, utils.KeyResetPassword(token), user.ID, ttl)

	if s.cfg.Mailer != nil {
		resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.cfg.AppFrontendURL, token)
		go s.cfg.Mailer.SendResetPasswordEmail(user.Email, user.Name, resetURL)
	}

	return nil
}

// ─── Reset Password ────────────────────────────────────────────────────────────

func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	ctxs := context.Background()

	// 1. Validasi password policy
	if errs := s.cfg.PasswordPolicy.Validate(req.NewPassword); len(errs) > 0 {
		return NewAuthError(422, strings.Join(errs, ", "))
	}

	// 2. Cek confirm password
	if req.NewPassword != req.ConfirmPassword {
		return NewAuthError(422, "Password konfirmasi tidak sesuai.")
	}

	// 3. Ambil user_id dari Redis
	key := utils.KeyResetPassword(req.Token)
	userIDStr, err := s.redis.Get(ctxs, key).Result()
	if err != nil {
		return NewAuthError(400, "Token tidak valid atau sudah kedaluwarsa.")
	}

	var userID int64
	fmt.Sscanf(userIDStr, "%d", &userID)

	// 4. Ambil user via users repository
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return NewAuthError(400, "Token tidak valid.")
	}

	// 5. Cek password history
	histories, _ := s.repo.GetPasswordHistories(ctx, userID, s.cfg.PasswordHistoryCount)
	hashes := make([]string, len(histories))
	for i, h := range histories {
		hashes[i] = h.PasswordHash
	}
	if err := utils.CheckPasswordHistory(req.NewPassword, hashes); err != nil {
		return NewAuthError(422, err.Error())
	}

	// 6. Hash & update password via users repository
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return NewAuthError(500, "Gagal memproses password.")
	}

	now := time.Now()
	user.Password = hashed
	user.PasswordChangedAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return NewAuthError(500, "Gagal mengupdate password.")
	}

	// 7. Simpan password history
	s.repo.SavePasswordHistory(ctx, &models.PasswordHistory{
		UserID:       userID,
		PasswordHash: hashed,
	})

	// 8. Blacklist semua token user & hapus Redis key
	s.repo.BlacklistAllUserTokens(ctx, userID)
	s.redis.Del(ctxs, key)

	return nil
}

// ─── Logout ────────────────────────────────────────────────────────────────────

// Logout memblacklist refresh token dari device saat ini
func (s *authService) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	// Parse token untuk ambil JTI
	claims, err := s.cfg.JWTManager.ParseToken(req.RefreshToken)
	if err != nil {
		// Token tidak valid / expired — anggap sudah logout, return success
		return nil
	}

	// Cek token ada di DB
	authToken, err := s.repo.GetTokenByJTI(ctx, claims.ID)
	if err != nil || authToken == nil {
		// Token tidak ditemukan — anggap sudah logout
		return nil
	}

	// Blacklist token
	if err := s.repo.BlacklistToken(ctx, claims.ID); err != nil {
		return NewAuthError(500, "Gagal logout. Coba lagi.")
	}

	return nil
}

// LogoutAll memblacklist semua refresh token user (logout dari semua device)
func (s *authService) LogoutAll(ctx context.Context, userID int64) error {
	if err := s.repo.BlacklistAllUserTokens(ctx, userID); err != nil {
		return NewAuthError(500, "Gagal logout dari semua perangkat.")
	}
	return nil
}

// ─── Private Helpers ───────────────────────────────────────────────────────────

// findUserByIdentifier mencari user by username atau email
func (s *authService) findUserByIdentifier(ctx context.Context, identifier string) (*userModels.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return s.userRepo.GetByEmail(ctx, identifier)
}

func (s *authService) handleFailedLogin(ctx context.Context, identifier, ip string, userID *int64, userAgent string) {
	lockDuration := time.Duration(s.cfg.LoginLockDurationMinutes) * time.Minute
	attemptsKey := utils.KeyLoginAttempts(identifier)

	attempts, _ := s.redis.Incr(ctx, attemptsKey).Result()
	s.redis.Expire(ctx, attemptsKey, lockDuration)

	if int(attempts) >= s.cfg.LoginMaxAttempts {
		s.redis.Set(ctx, utils.KeyLoginLock(identifier), 1, lockDuration)
		s.redis.Del(ctx, attemptsKey)
	}

	ipKey := utils.KeyLoginAttemptsIP(ip)
	ipAttempts, _ := s.redis.Incr(ctx, ipKey).Result()
	s.redis.Expire(ctx, ipKey, time.Minute)
	if int(ipAttempts) >= s.cfg.RateLimitLoginPerIPPerMinute {
		s.redis.Set(ctx, utils.KeyIPBlacklist(ip), 1, lockDuration)
	}

	s.saveLoginHistory(ctx, userID, identifier, ip, userAgent, models.LoginStatusFailed, "wrong_password")
}

func (s *authService) getRemainingAttempts(ctx context.Context, identifier string) int {
	attempts, _ := s.redis.Get(ctx, utils.KeyLoginAttempts(identifier)).Int()
	remaining := s.cfg.LoginMaxAttempts - attempts
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (s *authService) saveLoginHistory(ctx context.Context, userID *int64, identifier, ip, userAgent, status, failureReason string) {
	ua := userAgent
	history := &models.LoginHistory{
		UserID:     userID,
		Identifier: identifier,
		IPAddress:  ip,
		UserAgent:  &ua,
		Status:     status,
	}
	if failureReason != "" {
		history.FailureReason = &failureReason
	}
	s.repo.SaveLoginHistory(ctx, history)
}

// updateLastLogin mengupdate last_login_at via users repository
func (s *authService) updateLastLogin(ctx context.Context, user *userModels.User) {
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(ctx, user)
}

// buildTokenResponse membangun response token dari user model
func (s *authService) buildTokenResponse(ctx context.Context, accessToken, refreshToken string, user *userModels.User) *dto.TokenResponse {
	settings, _ := user.GetSettings()
	histories, _ := s.repo.GetUserLoginHistories(ctx, user.ID, 10)

	// 1. Ambil roles dari DB — tanpa preload permissions agar ringan
	roles, err := s.rbacRepo.GetUserRoles(user.ID)
	var roleSimple []rbacDto.RoleSimpleResponse
	if err == nil {
		roleSimple = rbacDto.ToRoleSimpleListResponse(roles)
	} else {
		roleSimple = []rbacDto.RoleSimpleResponse{}
	}

	// 2. Ambil semua permissions (dari role + direct) sebagai object lengkap
	// Gunakan map untuk deduplication berdasarkan permission ID
	permMap := make(map[int64]rbacDto.PermissionResponse)

	// 2a. Permissions dari role
	for _, role := range roles {
		for _, p := range role.Permissions {
			if _, exists := permMap[p.ID]; !exists {
				permMap[p.ID] = *rbacDto.ToPermissionResponse(&p)
			}
		}
	}

	// 2b. Direct permissions yang di-grant — override/tambah ke map
	directPerms, err := s.rbacRepo.GetUserDirectPermissions(user.ID)
	if err == nil {
		for _, up := range directPerms {
			if !up.IsGranted {
				// Direct deny — hapus dari map jika ada
				delete(permMap, up.PermissionID)
				continue
			}
			// Direct grant — tambah jika belum ada
			if _, exists := permMap[up.PermissionID]; !exists {
				perm, err := s.rbacRepo.GetPermissionByID(up.PermissionID)
				if err == nil && perm != nil {
					permMap[perm.ID] = *rbacDto.ToPermissionResponse(perm)
				}
			}
		}
	}

	// 2. Ambil data creator jika CreatedBy tidak nil
	var creatorDTO *userModels.UserCreator
	if user.CreatedBy != nil {
		creatorUser, err := s.userRepo.GetByID(ctx, *user.CreatedBy)
		if err == nil && creatorUser != nil {
			creatorDTO = &userModels.UserCreator{
				ID:       creatorUser.ID,
				Username: creatorUser.Username,
				Name:     creatorUser.Name,
			}
		}
	}

	// 3. Convert map ke slice
	permList := make([]rbacDto.PermissionResponse, 0, len(permMap))
	for _, p := range permMap {
		permList = append(permList, p)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.cfg.JWTManager.AccessExpSeconds(),
		User: userDto.UserResponse{
			ID:             user.ID,
			Photo:          user.Photo,
			PhotoThumbnail: user.PhotoThumbnail,
			Username:       user.Username,
			Email:          user.Email,
			Name:           user.Name,
			IsSuperadmin:   user.IsSuperadmin,
			IsStaff:        user.IsStaff,
			IsVerified:     user.IsVerified,
			Roles:          roleSimple, // ← Inject Roles ke DTO Response
			Permissions:    permList,   // ← Inject Permissions ke DTO Response
			Settings:       settings,
			Histories:      histories,
			LastLoginAt:    user.LastLoginAt,
			Creator:        creatorDTO,
			CreatedAt:      types.CustomTime(user.CreatedAt),
			UpdatedAt:      types.CustomTime(user.UpdatedAt),
		},
	}
}

// ─── Auth Error ────────────────────────────────────────────────────────────────

type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string { return e.Message }

func NewAuthError(code int, message string) *AuthError {
	return &AuthError{Code: code, Message: message}
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
