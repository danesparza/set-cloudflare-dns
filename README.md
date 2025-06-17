# set-cloudflare-dns

A simple command line tool written in Go for creating or modifying DNS records in Cloudflare.

## Usage

```
set-cloudflare-dns -zone example.com -name host.example.com -type A -content 1.2.3.4 -ttl 120 -token <api-token>
```

Arguments can also be provided via flags. The API token may be passed with the `-token` flag or via the `CF_API_TOKEN` environment variable.

`-type` defaults to `A`

`-ttl` defaults to `1` (which in Cloudflare is 'automatic')

Verbose logging can be enabled with `-v`.

## Docker
Pass the required flags as arguments and provide your API token either
with `-token` or the `CF_API_TOKEN` environment variable:

```bash
docker run --rm -e CF_API_TOKEN=<api-token> \
  danesparza/set-cloudflare-dns -zone example.com -name host.example.com -type A \
  -content 1.2.3.4 -ttl 120
```


