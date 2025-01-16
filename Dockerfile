ARG VERSION=0.0.0
ARG COMMIT=''
ARG DATE=''

FROM golang:1-alpine as builder
WORKDIR /go/src/headscale_exporter
COPY . .
RUN apk --no-cache add git openssh build-base
RUN go build -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o app .

FROM alpine as production
LABEL maintainer="msterhu" \
  org.opencontainers.image.url="https://github.com/msterhu/headscale_exporter" \
  org.opencontainers.image.source="https://github.com/msterhu/headscale_exporter" \
  org.opencontainers.image.vendor="msterhu" \
  org.opencontainers.image.title="headscale_exporter" \
  org.opencontainers.image.description="Prometheus exporter for Headscale." \
  org.opencontainers.image.licenses=""
RUN <<EOF
    apk add --no-cache ca-certificates libc6-compat \
    rm -rf /var/cache/apk/*
EOF
COPY --from=builder /go/src/headscale_exporter/app /app
ENTRYPOINT ["/app"]