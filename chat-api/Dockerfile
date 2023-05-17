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

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

USER appuser

COPY --from=builder --chown=appuser:appuser /src/bin/chat-api .

CMD ["/app/chat-api", "-addr", ":3000"]
