FROM golang:1.15.3 as builder
ENV DATA_DIRECTORY /go/src/finance-app
WORKDIR $DATA_DIRECTORY
ARG APP_VERSION
ARG CGO_ENABLED=0

COPY . .
RUN go build -ldflags="-X finance-app/internal/config.Version=$APP_VERSION" finance-app/cmd/server

FROM alpine:3.10
ENV DATA_DIRECTORY /go/src/finance-app
RUN apk add --update --no-cache \
    ca-certificates
COPY ./internal/database/migrations ${DATA_DIRECTORY}/internal/database/migrations
COPY --from=builder ${DATA_DIRECTORY}/server /finance-app

ENTRYPOINT ["/finance-app"]