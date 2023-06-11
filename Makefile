gen:
	protoc -I ./proto --go_out=proto/gen \
	--go-grpc_out=proto/gen \
	--grpc-gateway_out=proto/gen \
	shorten/shorten.proto
test:
	go test -cover ./...

mocks:
	mockgen -source=internal/service/service.go Repository > internal/mocks/repository_mock.go

run_db:
	docker compose up shorten_db --build

run_inmemory:
	docker compose up shorten_inmemory --build

