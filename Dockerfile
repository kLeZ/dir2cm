FROM library/golang:1.20.4-alpine

LABEL org.opencontainers.image.schema-version = "2.0.0"
LABEL org.opencontainers.image.title = "dir2cm"
LABEL org.opencontainers.image.description = "Creates a ConfigMap from a Directory"
LABEL org.opencontainers.image.vendor = "me.klez"
LABEL org.opencontainers.image.url = "https://github.com/kLeZ/dir2cm"
LABEL org.opencontainers.image.source = "https://github.com/kLeZ/dir2cm"
LABEL org.opencontainers.image.licenses = "MIT"

WORKDIR /go/src/github.com/kLeZ/dir2cm

RUN apk add --no-cache git

COPY . .

RUN go install -v ./...

VOLUME ["/data"]

WORKDIR /data

ENTRYPOINT ["dir2cm"]
