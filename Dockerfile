FROM golang:1.16-alpine as builder

ARG REVISION

RUN mkdir -p /feelguuds_platform/

WORKDIR /feelguuds_platform

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
    -X github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version.REVISION=${REVISION}" \
    -a -o bin/feelguuds_platform cmd/feelguuds_platform/*

RUN CGO_ENABLED=0 go build -ldflags "-s -w \
    -X github.com/yoanyombapro1234/FeelguudsPlatform/pkg/version.REVISION=${REVISION}" \
    -a -o bin/feelguuds_platform cmd/feelguuds_platform/*

FROM alpine:3.14

ARG BUILD_DATE
ARG VERSION
ARG REVISION

LABEL maintainer="yoanyomba"

RUN addgroup -S app \
    && adduser -S -G app app \
    && apk --no-cache add \
    ca-certificates curl netcat-openbsd

WORKDIR /home/app

COPY --from=builder /feelguuds_platform/bin/feelguuds_platform .
COPY --from=builder /feelguuds_platform/bin/feelguuds_platform_cli /usr/local/bin/feelguuds_platform_cli
COPY ./ui ./ui
RUN chown -R app:app ./

USER app

CMD ["./feelguuds_platform"]
