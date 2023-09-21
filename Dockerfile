FROM docker.io/golang:1.21.1 as builder
ARG VERSION=devel

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.version=${VERSION}'" ./cmd/guildops

FROM scratch

COPY --from=builder /build/guildops /guildops

EXPOSE 9252
USER 65534

ENTRYPOINT ["/guildops"]
