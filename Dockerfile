FROM --platform=linux/amd64 golang:1.18.1-alpine as BUILDER
RUN apk add --no-cache make git ca-certificates
WORKDIR /app
COPY . .
RUN make download
RUN make build-linux

FROM --platform=linux/amd64 scratch as RUNNER
WORKDIR /app
COPY --from=BUILDER /app/main /app/main
COPY --from=BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/app/main"]
