.PHONY: install-devtools
install-devtools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: format
format:
	gofmt -s -w .

.PHONY: imports
imports:
	goimports -l -w .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint: format imports tidy vet staticcheck

.PHONY: precommit
precommit: lint test

.PHONY: test
test:
	go test -v ./momento

.PHONY: vendor
vendor:
	go mod vendor