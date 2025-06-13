# build stage
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o set-cloudflare-dns

# final stage
FROM alpine:3.19
COPY --from=build /app/set-cloudflare-dns /usr/local/bin/set-cloudflare-dns
ENTRYPOINT ["set-cloudflare-dns"]
