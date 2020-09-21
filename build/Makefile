OWNER=Janrajc
IMAGE_NAME=kube-chorus-router
VERSION_FILE="../version/VERSION"
version=1.0.0
error=0.0.0


.PHONY: build
build:
	docker build -f Dockerfile -t $(IMAGE_NAME) .. -m 4g
	docker save -o $(IMAGE_NAME).tar $(IMAGE_NAME)
test:	
	@echo 'Commit Message'
	@echo $(TRAVIS_COMMIT_MESSAGE)
	go test -v ../... -coverprofile cp.out
	go tool cover -html=cp.out -o cover.html
doc:
	cd ../cmd;godoc -http=:9090
run:
	cd ../cmd;go run .
coverage:
	go test -v ../... -coverprofile=coverage.txt -covermode=atomic
	go tool cover -html=coverage.txt -o cover.html

publish:
	@echo 'Updating Version, and push image to quay'
	$(eval version = $(shell ./gitpush.sh | tail -1))
	@echo 'New versions $(version)'

clean:
	docker rmi -f $$(docker images -q -f dangling=true) || true
	docker rmi -f $$(docker images | awk '$$1 ~ /kube-chorus-router/ { print $$3}') || true
