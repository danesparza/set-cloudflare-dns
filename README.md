# set-cloudflare-dns

A simple command line tool written in Go for creating DNS records in Cloudflare.

## Usage

```
set-cloudflare-dns -zone example.com -name host.example.com -type A -content 1.2.3.4 -ttl 120 -token <api-token>
```

Arguments can also be provided via flags. The API token may be passed with the `-token` flag or via the `CF_API_TOKEN` environment variable.

Verbose logging can be enabled with `-v`.

## Building

```
go build
```

## Running tests

```
go test ./...
```

## Docker

A multi-stage `Dockerfile` is provided to build a minimal image:

```
docker build -t set-cloudflare-dns .
```

Once built (or after pulling the published image) you can run the tool via
Docker. Pass the required flags as arguments and provide your API token either
with `-token` or the `CF_API_TOKEN` environment variable:

```bash
docker run --rm -e CF_API_TOKEN=<api-token> \
  set-cloudflare-dns -zone example.com -name host.example.com -type A \
  -content 1.2.3.4 -ttl 120
```

## GitHub Actions

The included workflow builds and tests the project and, on pushes to `main`, publishes a Docker image to Docker Hub. The workflow expects the following secrets:

- `DOCKERHUB_USER`
- `DOCKERHUB_PASS`

## A Haiku

Cloudflare records bloom
Crafted with ChatGPT Codex
Automation flows


