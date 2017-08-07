SHELL:=/bin/bash
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

# For local testing, run `docker-compose up -d`
SERVER ?=http://127.0.0.1:8153/go/
export GOCD_URL=$(SERVER)
export GOCD_SKIP_SSL_CHECK=1


travis: before_install script after_success deploy_on_develop

before_install:
	go get -t -v ./...
	go get github.com/golang/lint/golint

script: test
	git diff-index HEAD --
	diff -u <(echo -n) <(gofmt -d -s .)
	bash ./scripts/clean-workspace.sh
	ls -lah ./godata/server/
	chmod -R 777 ./godata/
	make testacc

after_failure:
	docker-compose down

after_success:
	docker-compose down
	bash <(curl -s https://codecov.io/bash)
	go get github.com/goreleaser/goreleaser

deploy_on_tag:
	go get github.com/goreleaser/goreleaser
	gem install --no-ri --no-rdoc fpm
	go get
	goreleaser

deploy_on_develop:
	go get github.com/goreleaser/goreleaser
	gem install --no-ri --no-rdoc fpm
	go get
	goreleaser --snapshot

default: build

build: format
	go build

test: format
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
	$(MAKE) -C gocd test

testacc: format provision-test-gocd
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

format: lint
	gofmt -w $(GOFMT_FILES)
	$(MAKE) -C ./gocd fmt

lint:
	golint . gocd

provision-test-gocd:
	docker-compose up -d
	bash scripts/wait-for-test-server.sh

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
