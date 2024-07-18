# Changelog

## [0.2.2] - 2024-07-18

### Added

- Ability to add `CORS Middleware` with `AddCorsMiddleware` method.

### Changed

- Open Api schema middleware now will return `401` status code when making a request with missing `bearer token`.

## [0.2.1] - 2024-07-18

### Added

- **Development Proxy for Frontend Apps**: Introduced the `DEV_FRONTEND_PROXY` environment variable. When this variable is set, it will overwrite the `RootStaticServer` configuration.

## [0.2.0] - 2024-07-15

### Added

- Root Static Server
- Root Static Server: Ability to customize 404 pages with a `404.html` file.
- Root Static Server: Serving static content in `SPA` mode with an entry point of `index.html`.
- Api specific prefix when combined with `RootStaticServer`.
- `changelog-lint` pre commit hook.
- `changelog-lint` and `golang-ci` linting stages in github workflow.
- Helper `HtmlTemplate` in `response` package to respond with `html` templates using `html/template` package.

### Changed

- Renamed `response` package to `res`.

### Removed

- Custom responses for `405` status code. Bad requests will fallback to `404` handlers.

## [0.1.2] - 2024-07-10

### Added

- Changelog.
- Release Notes.

## [0.1.1] - 2024-07-10

### Fixed

- nil pointer causing segmentation fault when `SwaggerStaticServer` was unused.
