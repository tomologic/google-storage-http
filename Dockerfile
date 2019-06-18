FROM golang:1.12.6-alpine3.9 as builder
WORKDIR /go/src/google-storage-http
RUN apk add git --no-cache
COPY *.go .
RUN go get -v
RUN go build -v

FROM alpine:3.9
ENV PORT=8080
ENV LOGGING=true
RUN apk add ca-certificates --no-cache
COPY --from=builder /go/src/google-storage-http/google-storage-http /usr/bin/
ENTRYPOINT ["google-storage-http"]
