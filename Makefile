REPO ?= ghcr.io/brumhard/pix
TAG ?= $(shell svu minor)
# might be necessary to run the following to enable multiarch builds
# docker buildx create --driver docker-container --name multiarch
# docker buildx use multiarch
docker:
	docker buildx build --push -t $(REPO):$(TAG) --platform linux/amd64,linux/arm64,linux/arm/v7 .

release: docker
	@sed -i 's#image: $(REPO).*#image: $(REPO):$(TAG)#g' docker-compose.yml
	@git add docker-compose.yml && git commit --amend --no-edit
	git tag $(TAG)
	git push && git push --tags

phone:
	GOOS=android GOARCH=arm64 go build -o pix main.go
	chmod +x pix
	scp -P 8022 pix 192.168.0.239:~/
	rm pix

run:
	@pushd frontend >/dev/null && flutter build web && popd >/dev/null
	@go run main.go -images dav