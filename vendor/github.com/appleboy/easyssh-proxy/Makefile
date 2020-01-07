GOFMT ?= gofmt "-s"
GO ?= go
PACKAGES ?= $(shell $(GO) list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)

all: lint

fmt:
	$(GOFMT) -w $(SOURCES)

vet:
	$(GO) vet $(PACKAGES)

lint:
	@hash revive > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mgechev/revive; \
	fi
	revive -config .revive.toml ./... || exit 1

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(SOURCES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(SOURCES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(SOURCES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test: fmt-check
	@$(GO) test -v -cover -coverprofile coverage.txt $(PACKAGES) && echo "\n==>\033[32m Ok\033[m\n" || exit 1

embedmd:
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/campoy/embedmd; \
	fi
	embedmd -d *.md

clean:
	go clean -x -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST) vendor

ssh-server:
	adduser -h /home/drone-scp -s /bin/bash -D -S drone-scp
	echo drone-scp:1234 | chpasswd
	mkdir -p /home/drone-scp/.ssh
	chmod 700 /home/drone-scp/.ssh
	cp tests/.ssh/id_rsa.pub /home/drone-scp/.ssh/authorized_keys
	chown -R drone-scp /home/drone-scp/.ssh
	# install ssh and start server
	apk add --update openssh openrc
	rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key
	sed -i 's/AllowTcpForwarding no/AllowTcpForwarding yes/g' /etc/ssh/sshd_config
	./tests/entrypoint.sh /usr/sbin/sshd -D &

version:
	@echo $(VERSION)
