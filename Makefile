proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./protobuf/capture.proto


image-thor:
	docker build -t thor -f .docker/thor.Dockerfile .
image-huginn:
	docker build -t huginn -f .docker/huginn.Dockerfile .
image-muninn:
	docker build -t muninn -f .docker/muninn.Dockerfile .
image-heimdall:
	docker build -t heimdall -f .docker/heimdall.Dockerfile .

images: image-thor image-huginn image-muninn image-heimdall

up:
	docker-compose up -d