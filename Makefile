test:
	go list -f '{{.Dir}}/...' -m | xargs go test

submodule:
	@git submodule update --remote --recursive

gen:
	@./scripts/gen_proto.sh

gen-win:
	@powershell -ExecutionPolicy Bypass -File ./scripts/gen_proto.ps1

# make create-migration name="name"
create-migration:
	@go run github.com/pressly/goose/v3/cmd/goose postgres "user=taskem password=taskem dbname=taskem sslmode=disable" create $(name) sql