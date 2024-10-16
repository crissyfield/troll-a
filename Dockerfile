# Build
FROM golang:1.23-alpine3.20 AS build

RUN apk --no-cache add build-base \
                       git \
                       re2-dev

WORKDIR /app
COPY . .

RUN make build-re2

# Deploy
FROM alpine:3.19 AS deploy

RUN apk --no-cache add re2 \
                       tini \
                       tzdata

COPY --from=build /app/troll-a /troll-a
COPY --from=build /usr/lib/libabsl_flags_* /usr/lib

ENTRYPOINT [ "/sbin/tini", "--", "/troll-a" ]
