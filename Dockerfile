FROM golang:1.16.2-alpine3.13 AS builder
WORKDIR /go/src/github.com/buneyev/network-access-checker/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app/network-access-checker .

FROM alpine:3.13  
RUN addgroup -g 2000 app && \
    adduser -u 2000 -G app -D app && \
    apk --no-cache add ca-certificates && \
    apk --no-cache add bash && \
    apk --no-cache add busybox-extras && \
    apk --no-cache add curl
WORKDIR /app/
COPY --chown=app:app --from=builder /go/src/github.com/buneyev/network-access-checker/app/network-access-checker .
USER app
ENTRYPOINT [ "./network-access-checker" ]