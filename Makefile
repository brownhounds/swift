lint:
	@golangci-lint run --fast

fix:
	@golangci-lint run --fix

install-hooks:
	@pre-commit install
