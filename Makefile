proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./protobuf/capture.proto


image-thor:
	docker build -t thor -f .docker/capture.Dockerfile .