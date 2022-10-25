FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS builder

WORKDIR /app
COPY ./src ./
#RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf
RUN go mod download
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/audit-server .

FROM alpine as runner
WORKDIR /app
COPY --from=builder /out/audit-server .

EXPOSE 8080

CMD ["/app/audit-server"]