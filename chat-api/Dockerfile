FROM artifactory.cloud.cms.gov/docker/golang:alpine as builder

WORKDIR /app

COPY . .

RUN mkdir bin && \
  go build -o bin/ ./...

FROM artifactory.cloud.cms.gov/docker/alpine:latest

WORKDIR /app

EXPOSE 3000

COPY --from=builder /app/bin/chat-api .

CMD ["/app/chat-api", "-addr", ":3000"]
