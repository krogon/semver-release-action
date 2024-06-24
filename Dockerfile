# Build
FROM golang:1.19-alpine AS build

RUN apk add --no-cache --update make \
    && rm -f /var/cache/apk/*

WORKDIR /go/src/app
ADD ./ /go/src/app

RUN make build

# Run
FROM alpine:3.20

COPY entrypoint.sh /entrypoint.sh
COPY --from=build /go/src/app/bumper /bumper

ENTRYPOINT ["/entrypoint.sh"]
