FROM artifactory.cloud.cms.gov/docker/golang:alpine as builder

WORKDIR /app

COPY . .

RUN mkdir bin && \
  go build -o bin/ ./...

FROM artifactory.cloud.cms.gov/docker/alpine:latest

WORKDIR /app

ENV PORT 5000

EXPOSE ${PORT}

COPY --from=builder /app/bin/counter-api .

CMD ["/app/counter-api"]
