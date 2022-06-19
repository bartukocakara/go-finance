FROM golang:1.16 as builder
ENV DATA_DIRECTORY /go/src/financial-app
WORKDIR $DATA_DIRECTORY
ARG APP_VERSION
ARG CGO_ENABLED=0

COPY . .
RUN go build -ldflags="-X financial-app/internal/config.Version=$APP_VERSION" financial-app/cmd/server

FROM alpine:3.10
ENV DATA_DIRECTORY /go/src/financial-app
RUN apk add --update --no-cache \
    ca-certificates
COPY ./internal/database/migrations ${DATA_DIRECTORY}/internal/database/migrations
COPY --from=builder ${DATA_DIRECTORY}/server /financial-app

ENTRYPOINT ["/financial-app"]