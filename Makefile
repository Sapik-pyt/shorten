gen:
	protoc -I ./proto --go_out=proto/gen \
	--go-grpc_out=proto/gen \
	--grpc-gateway_out=proto/gen \
	shorten/shorten.proto