submodule:
	@git submodule update --remote --recursive

gen:
	@./scripts/gen_proto.sh

gen-win:
	@powershell -ExecutionPolicy Bypass -File ./scripts/gen_proto.ps1

create-migration:
	@GOOSE_MIGRATION_DIR=./migrations go run github.com/pressly/goose/v3/cmd/goose postgres "user=stream password=stream dbname=stream sslmode=disable" create $(name) sql