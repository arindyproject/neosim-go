# Makefile for Neosim Go Project

.PHONY: build run migrate migrate-dev migrate-prod clean test

# Build the API
build:
	go build -o bin/api ./cmd/api

# Run the API
run:
	go run ./cmd/api/main.go

# Run migrations using GORM auto-migration (development)
migrate-dev:
	@echo "Running GORM auto-migration (development)..."
	go run ./cmd/migrate/main.go -env=DEV -type=gorm

# Run migrations using GORM auto-migration (production)
migrate-prod:
	@echo "Running GORM auto-migration (production)..."
	go run ./cmd/migrate/main.go -env=PROD -type=gorm

# Run manual SQL migrations (development)
migrate-sql:
	@echo "Running SQL-based migrations (development)..."
	go run ./cmd/migrate/main.go -env=DEV -type=sql

# Run manual SQL migrations (production)
migrate-sql-prod:
	@echo "Running SQL-based migrations (production)..."
	go run ./cmd/migrate/main.go -env=PROD -type=sql

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test ./...

# Create new migration file (example)
create-migration:
	@echo "Creating new migration file..."
	@read -p "Enter migration name: " name; \
	touch internal/module/$$name/migrations/$$(date +%Y%m%d%H%M%S)_$$name.go

# Show database stats
db-stats:
	go run -c 'package main; import ("fmt"; "neosim_go/config"); func main() { cfg := config.LoadConfig("DEV"); db, _ := cfg.ConnectDB(); defer config.CloseDB(db); stats, _ := config.GetDBStats(db); fmt.Printf("%+v\n", stats) }'



# ============================================================
# Fresh Migrations (Drop All + Re-migrate)
# ============================================================

migrate-fresh-dev:
	@echo "⚠️  WARNING: This will DROP ALL TABLES on DEV and re-migrate!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🔄 Running fresh migration (DEV)..."; \
		go run ./cmd/migrate/main.go -env=DEV -type=gorm -fresh=true; \
	else \
		echo "❌ Aborted."; \
	fi

migrate-fresh-dev-sql:
	@echo "⚠️  WARNING: This will DROP ALL TABLES on DEV and re-migrate!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🔄 Running fresh SQL migration (DEV)..."; \
		go run ./cmd/migrate/main.go -env=DEV -type=sql -fresh=true; \
	else \
		echo "❌ Aborted."; \
	fi

migrate-fresh-prod:
	@echo "🚨 DANGER: This will DROP ALL TABLES on PRODUCTION!"
	@echo "🚨 This action is IRREVERSIBLE!"
	@read -p "Type 'PRODUCTION' to confirm: " confirm; \
	if [ "$$confirm" = "PRODUCTION" ]; then \
		read -p "Are you absolutely sure? (yes/no): " confirm2; \
		if [ "$$confirm2" = "yes" ]; then \
			echo "🔄 Running fresh migration (PROD)..."; \
			go run ./cmd/migrate/main.go -env=PROD -type=gorm -fresh=true -force=true; \
		else \
			echo "❌ Aborted."; \
		fi \
	else \
		echo "❌ Aborted. You must type 'PRODUCTION' exactly."; \
	fi

migrate-fresh-prod-sql:
	@echo "🚨 DANGER: This will DROP ALL TABLES on PRODUCTION!"
	@echo "🚨 This action is IRREVERSIBLE!"
	@read -p "Type 'PRODUCTION' to confirm: " confirm; \
	if [ "$$confirm" = "PRODUCTION" ]; then \
		read -p "Are you absolutely sure? (yes/no): " confirm2; \
		if [ "$$confirm2" = "yes" ]; then \
			echo "🔄 Running fresh SQL migration (PROD)..."; \
			go run ./cmd/migrate/main.go -env=PROD -type=sql -fresh=true -force=true; \
		else \
			echo "❌ Aborted."; \
		fi \
	else \
		echo "❌ Aborted. You must type 'PRODUCTION' exactly."; \
	fi


# ─── Seeder ────────────────────────────────────────────────────────────────────

# Jalankan seeder (skip jika data sudah ada)
seed:
	@echo "🌱 Menjalankan seeder..."
	go run ./cmd/seed/main.go -env=DEV

# Jalankan seeder di production
seed-prod:
	@echo "🌱 Menjalankan seeder (PROD)..."
	go run ./cmd/seed/main.go -env=PROD

# Fresh seed: hapus semua data lalu seed ulang (DEV only)
seed-fresh:
	@echo "⚠️  WARNING: Ini akan menghapus semua data dan seed ulang!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🔄 Running fresh seed..."; \
		go run ./cmd/seed/main.go -env=DEV -fresh=true; \
	else \
		echo "❌ Aborted."; \
	fi

# Migrate + seed sekaligus (DEV)
migrate-seed:
	@echo "🔄 Migrate + seed (DEV)..."
	go run ./cmd/migrate/main.go -env=DEV -type=gorm
	go run ./cmd/seed/main.go -env=DEV

# Fresh migrate + seed (DEV)
migrate-fresh-seed:
	@echo "⚠️  WARNING: Drop semua tabel, migrate ulang, dan seed!"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "🔄 Fresh migrate + seed..."; \
		go run ./cmd/migrate/main.go -env=DEV -type=gorm -fresh=true; \
		go run ./cmd/seed/main.go -env=DEV; \
	else \
		echo "❌ Aborted."; \
	fi