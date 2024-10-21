BIN=nexus

VERSION=$(shell git describe --tags)
DATE=$(shell date -u +%Y-%m-%d)
COMMIT=$(shell git rev-parse --short HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X github.com/MizuchiLabs/ssh-nexus/tools/updater.Version=${VERSION} -X github.com/MizuchiLabs/ssh-nexus/tools/updater.BuildDate=${DATE} -X github.com/MizuchiLabs/ssh-nexus/tools/updater.Commit=${COMMIT}"

all: clean build

.PHONY: clean
clean:
	rm -rf $(PWD)/bin/ $(PWD)/web/build

.PHONY: audit
audit:
	go fmt ./...
	go vet ./...
	go mod tidy
	go mod verify
	go test -v ./...
	- gosec --exclude=G104 ./...
	- govulncheck -show=color ./...
	- staticcheck -checks=all -f=stylish ./...

.PHONY: build
build: audit proto
	cd web && pnpm install && pnpm run build
	go build $(LDFLAGS) -o $(BIN) cmd/server/main.go && upx $(BIN)
	go build $(LDFLAGS) -o $(BIN)-agent cmd/agent/main.go && upx $(BIN)-agent

build-fast:
	go build $(LDFLAGS) -o bin/$(BIN) cmd/server/main.go
	go build $(LDFLAGS) -o bin/$(BIN)-agent cmd/agent/main.go

.PHONY: proto
proto:
	cd api/proto && buf generate && buf lint

.PHONY: docker
docker:
	cd web && pnpm install && pnpm run build
	go run test/cmd/mock.go # generate mock data
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o nexus-linux-amd64 cmd/server/main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o nexus-linux-arm64 cmd/server/main.go
	docker build \
		--label "org.opencontainers.image.source=https://github.com/MizuchiLabs/ssh-nexus" \
		--label "org.opencontainers.image.description=SSH Nexus" \
		--label "org.opencontainers.image.version=${VERSION}" \
		--label "org.opencontainers.image.revision=${COMMIT}" \
		--label "org.opencontainers.image.created=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		--label "org.opencontainers.image.licenses=Apache-2.0" \
		-t ghcr.io/mizuchilabs/ssh-nexus:latest .
	rm nexus-linux-amd64 nexus-linux-arm64

.PHONY: upgrade
upgrade:
	go get -u ./cmd/server && go mod tidy
	go get -u ./cmd/agent && go mod tidy
	cd web && pnpm update

.PHONY: test
test:
	go test -v -coverprofile=coverage.out ./...

.PHONY: release
release:
	goreleaser release --clean

.PHONY: test-db
test-db:
	go run test/cmd/mock.go
	rm -rf pb_data
	mv test_pb_data pb_data

.PHONY: snapshot
snapshot:
	goreleaser release --snapshot --clean

.PHONY: run
run-agent:
	go run cmd/agent/main.go

run-server:
	go run cmd/server/main.go

run-web:
	cd web && npm run dev
