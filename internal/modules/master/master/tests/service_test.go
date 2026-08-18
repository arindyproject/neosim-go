package tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"neosim_go/config"
	"neosim_go/internal/modules/master/master/services"
	"neosim_go/internal/modules/master/master/tests/mocks"

	masterContracts "neosim_go/internal/modules/master/master/contracts"
	"neosim_go/internal/shared/cache"

	he "neosim_go/internal/shared/httputil"
)

func TestMain(m *testing.M) {
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")
	fmt.Println("\033[35m  Master Service Test Suite\033[0m")
	fmt.Println("\033[34m" + strings.Repeat("─", 55) + "\033[0m")

	code := m.Run()

	if code == 0 {
		fmt.Println("\n\033[32m✓  PASS\033[0m  neosim_go/internal/modules/master/master")
	} else {
		fmt.Println("\n\033[31m✗  FAIL\033[0m  neosim_go/internal/modules/master/master")
	}

	os.Exit(code)
}

type MasterServiceTestSuite struct {
	suite.Suite
	repo     *mocks.MasterRepositoryMock
	rbacRepo *mocks.RBACRepositoryMock
	authRepo *mocks.AuthRepositoryMock
	svc      masterContracts.Service
	cfg      *config.Config
}

func (s *MasterServiceTestSuite) SetupTest() {
	s.repo = new(mocks.MasterRepositoryMock)
	s.rbacRepo = new(mocks.RBACRepositoryMock)
	s.authRepo = new(mocks.AuthRepositoryMock)

	// 1. INSIALISASI CONFIG (Mencegah nil pointer dereference pada s.cfg)
	s.cfg = &config.Config{
		DefaultPageSize:    10,
		DefaultPageSizeMax: 100,
	}

	// 2. INSIALISASI CACHE MANAGER
	// Memakai cacheManager dummy (disabled) agar tidak bergantung Redis
	cacheManager := cache.NewManager(nil, false, int(5*time.Minute))

	// 3. PASANG KE SERVICE
	s.svc = services.NewMasterService(s.repo, s.rbacRepo, s.authRepo, cacheManager, s.cfg)
}

func TestMasterService(t *testing.T) {
	suite.Run(t, new(MasterServiceTestSuite))
}

func superadminActor() he.AuthContext {
	return he.AuthContext{UserID: 1, IsSuperadmin: true}
}

func regularActor() he.AuthContext {
	return he.AuthContext{UserID: 2, IsSuperadmin: false}
}

func (s *MasterServiceTestSuite) mockHasPermission(perm string, result bool) {
	s.rbacRepo.On("HasPermission", regularActor().UserID, perm, mock.Anything).Return(result, nil).Maybe()
}

func (s *MasterServiceTestSuite) mockNoPermissions() {
	s.rbacRepo.On("HasPermission", regularActor().UserID, mock.Anything, mock.Anything).Return(false, nil)
}
