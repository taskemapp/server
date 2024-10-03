TEST_DIR := ./...
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

test:
	go test -v -coverprofile=$(COVERAGE_FILE) $(TEST_DIR)

coverage: test
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

submodule:
	@git submodule update --remote --recursive

server-compose:
	docker compose build server && docker compose up server

run-server:
	go run apps/server/cmd/server/main.go

gen:
	@./scripts/gen_proto.sh

gen-win:
	@powershell -ExecutionPolicy Bypass -File ./scripts/gen_proto.ps1

# make create-migration name="name"
create-migration:
	@go run github.com/pressly/goose/v3/cmd/goose postgres "user=taskem password=taskem dbname=taskem sslmode=disable" create $(name) sql