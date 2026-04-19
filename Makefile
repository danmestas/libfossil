.PHONY: test test-drivers test-otel test-all vet build setup-hooks

test:
	go test ./... -count=1 -timeout=120s

test-drivers:
	@echo "=== modernc (default) ==="
	go test ./... -count=1 -timeout=120s
	@echo "=== ncruces ==="
	go test -tags test_ncruces ./... -count=1 -timeout=120s
	@echo "=== all drivers passed ==="

test-otel:
	cd observer/otel && GOWORK=off go test ./... -count=1

test-all: test-drivers test-otel build
	@echo "=== all targets passed ==="

vet:
	go vet ./...

build:
	go build ./cmd/libfossil/

setup-hooks:
	git config core.hooksPath .githooks
	@echo "Pre-commit hook installed. Runs both drivers + DST + OTel (~45s) before each commit."
	@echo "Skip with: git commit --no-verify"
