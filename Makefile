GOMARKDOC_VERSION := v1.1.0

.PHONY: test test-drivers test-otel test-all vet build setup-hooks docs-gen-sdk docs-gen-llms

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

docs-gen-sdk:
	@command -v gomarkdoc >/dev/null 2>&1 || go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@$(GOMARKDOC_VERSION)
	gomarkdoc --output docs/site/content/docs/reference/sdk/libfossil/api.md     ./
	gomarkdoc --output docs/site/content/docs/reference/sdk/cli/api.md           ./cli/
	gomarkdoc --output docs/site/content/docs/reference/sdk/db/api.md            ./db/
	gomarkdoc --output docs/site/content/docs/reference/sdk/observer/otel/api.md ./observer/otel/
	gomarkdoc --output docs/site/content/docs/reference/sdk/dst/api.md           ./dst/

docs-gen-llms:
	bash scripts/gen-llms-txt.sh
