lint:
	@golangci-lint run --fast

changelog-lint:
	@changelog-lint

fix:
	@golangci-lint run --fix

install-hooks:
	@pre-commit install

install-changelog-lint:
	@go install github.com/chavacava/changelog-lint@master

dev-version:
	@./scripts/dev-version.sh

git-tag:
	git tag --sign v$(v) -m v$(v)
	git push origin v$(v)
