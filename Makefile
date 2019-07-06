VERSION=$(shell cat VERSION)
CASSANDRA_VERSION=$(shell cat CASSANDRA_VERSION)
COMMIT=$(shell git rev-parse HEAD)
SHORT_COMMIT=$(shell git rev-parse --short HEAD)
REV=$(COMMIT)
CHANGED=$(shell git status -s)
ifneq ($(strip $(CHANGED)),)
    REV="$(COMMIT)-wip"
endif

.PHONY: build
build:
	@mkdir -p build
	@echo "building cassandra $(CASSANDRA_VERSION) $(VERSION)-$(REV)..."
	go build -ldflags "-X github.com/linuxuser586/cassandra/release.Version=$(VERSION) \
	-X github.com/linuxuser586/cassandra/release.CassandraVersion=$(CASSANDRA_VERSION) \
	-X github.com/linuxuser586/cassandra/release.Commit=$(REV) \
	-X github.com/linuxuser586/cassandra/release.ShortCommit=$(SHORT_COMMIT)" \
	-o build/cassandra cmd/cassandra-tools/*.go

.PHONY: dist
dist: build test

.PHONY: release
release: image
	docker push linuxuser586/cassandra:$(CASSANDRA_VERSION)
	docker push linuxuser586/cassandra:$(CASSANDRA_VERSION)-$(VERSION)
	docker push linuxuser586/cassandra:latest

.PHONY: image
image:
	docker build -t linuxuser586/cassandra:$(CASSANDRA_VERSION) .
	docker tag linuxuser586/cassandra:$(CASSANDRA_VERSION) linuxuser586/cassandra:$(CASSANDRA_VERSION)-$(VERSION)
	docker tag linuxuser586/cassandra:$(CASSANDRA_VERSION) linuxuser586/cassandra:latest

.PHONY: test
test:
	go test -race ./...

.PHONY: test-fast
test-fast:
	go test ./...

.PHONY: cover
cover:
	go test -cover ./...

.PHONY: cover-html
cover-html:
	@mkdir -p build
	go test -coverprofile=build/coverage.out ./...
	@go tool cover -html=build/coverage.out -o build/coverage.html
	@firefox build/coverage.html

.PHONY: clean
clean:
	rm -rf build