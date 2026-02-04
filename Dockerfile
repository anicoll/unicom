FROM golang:1.25.7-alpine as BUILDER
RUN apk add --no-cache make git ca-certificates

FROM scratch as RUNNER
WORKDIR /app
COPY /main /app/main
COPY --from=BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/app/main"]
