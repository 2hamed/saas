IMAGE_REGISTRY=gcr.io/cloud-voyage

proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./protobuf/capture.proto


image-thor: 
	docker build -t $(IMAGE_REGISTRY)/thor -f .docker/thor.Dockerfile .
image-huginn:
	docker build -t $(IMAGE_REGISTRY)/huginn -f .docker/huginn.Dockerfile .
image-muninn:
	docker build -t $(IMAGE_REGISTRY)/muninn -f .docker/muninn.Dockerfile .
image-heimdall:
	docker build -t $(IMAGE_REGISTRY)/heimdall -f .docker/heimdall.Dockerfile .

images: image-thor image-huginn image-muninn image-heimdall

up:
	docker-compose up -d

push:
	docker push $(IMAGE_REGISTRY)/thor
	docker push $(IMAGE_REGISTRY)/huginn
	docker push $(IMAGE_REGISTRY)/muninn
	docker push $(IMAGE_REGISTRY)/heimdall

