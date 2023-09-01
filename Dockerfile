FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

WORKDIR /app
COPY ./src ./
#############JUST FOR KAFKA####################################
RUN apk add --no-cache git
RUN apk add librdkafka-dev pkgconf build-base
#############JUST FOR KAFKA####################################

RUN go mod download
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags musl -o /out/audit-server .

FROM alpine as runner
WORKDIR /app
COPY --from=builder /out/audit-server .

EXPOSE 8080

CMD ["/app/audit-server"]