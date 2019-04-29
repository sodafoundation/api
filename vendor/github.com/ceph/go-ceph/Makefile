DOCKER_CI_IMAGE = go-ceph-ci
build:
	go build -v
fmt:
	go fmt ./...
test:
	go test -v ./...

test-docker: .build-docker
	docker run --device /dev/fuse --cap-add SYS_ADMIN --security-opt apparmor:unconfined --rm -it -v $(CURDIR):/go/src/github.com/ceph/go-ceph $(DOCKER_CI_IMAGE)

.build-docker:
	docker build -t $(DOCKER_CI_IMAGE) .
	@docker inspect -f '{{.Id}}' $(DOCKER_CI_IMAGE) > .build-docker

check:
	# TODO: add this when golint is fixed	@for d in $$(go list ./... | grep -v /vendor/); do golint -set_exit_status $${d}; done
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done
