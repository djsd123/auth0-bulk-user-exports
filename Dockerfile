FROM golang:1.13-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

RUN apt-get update -y && \
    apt-get install ca-certificates

WORKDIR /build

COPY . .

RUN go mod download && \
    go mod verify

# Build fails if tests fail
RUN go test ./...

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./auth0-bulk-user-exports

WORKDIR /dist

RUN cp /build/auth0-bulk-user-exports ./auth0-bulk-user-exports

# This app uses dynamic linking due to the indirect use of CGO libraries:
# - crypto/x509
# - net
# - runtime/cgo

# These RUNs ensure dependencies are copied to the resulting image
RUN ldd auth0-bulk-user-exports | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'

RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --chown=65534:0 --from=builder /dist /

WORKDIR /dist

# Run as nobody
USER 65534

ENTRYPOINT ["/auth0-bulk-user-exports"]
