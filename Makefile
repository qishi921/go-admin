.PHONY: all build run test clean tidy lint migrate seed cli db-check api-test

# 默认目标
all: tidy build

# 整理依赖
tidy:
	go mod tidy

# 构建项目
build:
	go build -o bin/admin-server main.go

# 运行服务
run:
	go run main.go

# 运行测试
test:
	go test ./... -v -count=1

# 运行测试（带覆盖率）
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 清理构建产物
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 代码检查
lint:
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run ./...

# 数据库迁移（运行服务时自动执行）
migrate:
	@echo "数据库迁移在服务启动时自动执行"
	go run main.go &

# CLI 工具：生成密码哈希
cli:
	go run cmd/admin-cli/main.go hash

# 数据库检查工具
db-check:
	go run cmd/db-check/main.go storage/database.db

# API 测试
api-test:
	./scripts/test-api.sh

# 开发模式（热重载，需要安装 air）
dev:
	@which air > /dev/null || go install github.com/air-verse/air@latest
	air

# 格式化代码
fmt:
	go fmt ./...

# 检查代码
vet:
	go vet ./...

# 静态检查
staticcheck:
	@which staticcheck > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# 安装开发工具
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# 创建管理员用户（需要先运行服务）
create-admin:
	@echo "使用以下 SQL 创建管理员用户（在 SQLite 中执行）："
	@echo "INSERT INTO users (username, password, email, status) VALUES ('admin', '\$2a\$10\$...', 'admin@example.com', 'active');"

# 备份数据库
backup-db:
	mkdir -p backups
	cp storage/database.db backups/database_$$(date +%Y%m%d_%H%M%S).db
	@echo "数据库已备份到 backups/ 目录"

# 帮助
help:
	@echo "可用命令："
	@echo "  make build      - 构建项目"
	@echo "  make run        - 运行服务"
	@echo "  make test       - 运行测试"
	@echo "  make test-cover - 运行测试并生成覆盖率报告"
	@echo "  make tidy       - 整理依赖"
	@echo "  make lint       - 代码检查"
	@echo "  make fmt        - 格式化代码"
	@echo "  make cli        - 运行 CLI 工具"
	@echo "  make db-check   - 检查数据库"
	@echo "  make api-test   - 运行 API 测试"
	@echo "  make dev        - 开发模式（热重载）"
	@echo "  make backup-db  - 备份数据库"
