FROM docker.io/golang:1.22.5 as builder
ARG VERSION=devel

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.version=${VERSION}'" ./cmd/guildops

FROM docker.io/alpine:3.20.1
# renovate: datasource=repology depName=alpine_3_18/ca-certificates versioning=loose
ARG CA_CERTIFICATES_VERSION=20230506-r0

COPY --from=builder /build/guildops /guildops

RUN apk add --no-cache ca-certificates=${CA_CERTIFICATES_VERSION}

EXPOSE 9252
USER 65534

ENTRYPOINT ["/guildops"]
