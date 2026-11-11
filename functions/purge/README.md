# Purge function

Go code of the Scaleway Function that deletes old tags from the registry namespace. It is deployed by [`functions.tf`](../../functions.tf), which packages this folder into a zip file. This README is not included in the zip.

## Go version

Scaleway builds the function with the Go version of its runtime (`runtime` in [`functions.tf`](../../functions.tf), e.g. `go126` for Go 1.26), and can't download a newer one. Hence:

- **The runtime must provide the Go version required by the `go` directive of `go.mod`.** When raising the `go` directive, e.g. after a dependency update, move the runtime to the same Go version in the same pull request. Only use a runtime Scaleway offers: `scw function runtime list`.
- **The `go` directive must not require a patch release:** write `go 1.26.0`, not `go 1.26.3`. Scaleway's runtime may run an older patch release, and the build would then fail (see [Scaleway's documentation](https://www.scaleway.com/en/docs/serverless-functions/troubleshooting/go-function-version-error/)). Scaleway suggests the short form `go 1.26`, but Go refuses it here: the Scaleway SDK requires `go 1.26.0`, which Go considers newer than `go 1.26`. `go mod init` writes the full version of your local Go, so check the directive after running it.

The CI enforces both rules: it builds with the runtime's Go version and `GOTOOLCHAIN=local`, and rejects a `go` directive requiring a patch release.

## Local checks

Run the same checks as the CI from this folder, with the runtime's Go version:

```sh
export GOTOOLCHAIN=local
gofmt -l .  # must print nothing
go vet ./...
go build ./...
```

`init()` panics without the `REGISTRY_*` environment variables, so run the tests with format-valid dummy values:

```sh
REGISTRY_ACCESS_KEY=SCW0123456789ABCDEF1 \
REGISTRY_NAMESPACE=test \
REGISTRY_PROJECT_ID=01234567-89ab-cdef-0123-456789abcdef \
REGISTRY_REGION=fr-par \
REGISTRY_SECRET_KEY=01234567-89ab-cdef-0123-456789abcdef \
go test ./...
```

Without Go installed locally, use the official image matching the runtime, e.g.:

```sh
docker run --rm -v "$PWD":/src -w /src golang:1.26 go vet ./...
```
