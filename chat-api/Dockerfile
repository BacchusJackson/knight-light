FROM artifactory.cloud.cms.gov/docker/golang:alpine as builder

ARG GO_VERSION=1.20
WORKDIR /src

COPY . .

RUN mkdir bin && \
  CGO_ENABLED=0 go build -o bin/ ./...

FROM artifactory.cloud.cms.gov/docker/alpine:latest

WORKDIR /app

EXPOSE 3000

ARG UID=10001

COPY --from=builder /src/bin/chat-api .

RUN adduser --disabled-password --gecos "" --home "/nonexistent" \
    --shell "/sbin/nologin" --no-create-home --uid "${UID}" appuser && \
    chown -R appuser:appuser /app

USER appuser


CMD ["/app/chat-api", "-addr", ":3000"]
