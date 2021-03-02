IMAGE_REGISTRY=gcr.io/cloud-voyage
GCP_AUTH_FILE=$(shell base64 ./gcloud-config.json)

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
push-thor:
	docker push $(IMAGE_REGISTRY)/thor
push-heimdall:
	docker push $(IMAGE_REGISTRY)/heimdall
push-muninn:
	docker push $(IMAGE_REGISTRY)/muninn
push-huginn:
	docker push $(IMAGE_REGISTRY)/huginn
push: push-thor push-heimdall push-muninn push-huginn

install-thor:
	helm install thor .helm/thor --set env.gcp.credentials="$(GCP_AUTH_FILE)"

install-odin:
	helm install odin .helm/odin --set rabbitmq.auth.username=rabbit,rabbitmq.auth.password=rabbitpass

install-heimdall:
	helm install heimdall .helm/heimdall

install: install-thor install-odin install-heimdall

