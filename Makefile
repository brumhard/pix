run:
	@pushd -q frontend && flutter build web && popd -q
	@go run main.go -images dav


REPO ?= ghcr.io/brumhard/pix
TAG ?= $(shell svu n)
# might be necessary to run the following to enable multiarch builds
# docker buildx create --driver docker-container --name multiarch --use
docker:
	docker buildx build --push -t $(REPO):$(TAG) --platform linux/amd64,linux/arm64,linux/arm/v7 .