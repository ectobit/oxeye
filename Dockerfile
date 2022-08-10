# syntax=docker/dockerfile:1.3
FROM golang:1.19.0-alpine AS builder

RUN --mount=type=cache,target=/var/cache/apk if [ "${TARGETPLATFORM}" = "linux/amd64" ]; \
    then apk add --no-cache git tzdata upx; \
    else apk add --no-cache git tzdata; fi

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod tidy

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags='-s -w -extldflags "-static"' -o /app/app .
RUN if [ "${TARGETPLATFORM}" = "linux/amd64" ]; then upx /app/app; fi

FROM alpine:3.16.2

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/app /usr/bin/app

USER 65534:65534

ENTRYPOINT ["app"]
