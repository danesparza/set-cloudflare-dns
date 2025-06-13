# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Start from the latest golang base image
FROM --platform=$BUILDPLATFORM golang:1.24.0 AS builder

ARG packagePath
ARG buildNum
ARG circleSha
ARG TARGETOS
ARG TARGETARCH

# Set the Current Working Directory inside the container
WORKDIR /app

# Set git credentials for private repo access
# RUN git config --global credential.helper store && echo "${DOCKER_GIT_CREDENTIALS}" > ~/.git-credentials

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Emit debug info
RUN echo "Building for OS: $TARGETOS, Arch: $TARGETARCH"
RUN echo "Package Path: $packagePath"
RUN echo "Build Number: $buildNum"
RUN echo "Commit SHA: $circleSha"

# Build the Go app for the correct architecture
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -ldflags "-X ${packagePath}/version.BuildNumber=${buildNum} -X ${packagePath}/version.CommitID=${circleSha} -X '${packagePath}/version.Prerelease=-'" -installsuffix cgo -o set-cloudflare-dns ./cmd/set-cloudflare-dns

######## Start a new stage #######
FROM --platform=$TARGETPLATFORM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/set-cloudflare-dns /usr/local/bin/set-cloudflare-dns
ENTRYPOINT ["set-cloudflare-dns"]
