REPO ?= ghcr.io/brumhard/pix
TAG ?= $(shell svu minor)
# might be necessary to run the following to enable multiarch builds
# docker buildx create --driver docker-container --name multiarch --use
release:
	@sed -i 's#image:.*#image: $(REPO):$(TAG)#g' docker-compose.yml
	@git add docker-compose.yml && git commit --amend --no-edit
	git tag $(TAG)
	docker buildx build --push -t $(REPO):$(TAG) --platform linux/amd64,linux/arm64,linux/arm/v7 .
	git push && git push --tags

run:
	@pushd frontend >/dev/null && flutter build web && popd >/dev/null
	@go run main.go -images dav