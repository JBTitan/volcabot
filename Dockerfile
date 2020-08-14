FROM golang:1.15-alpine AS builder

ADD . /src
WORKDIR /src
RUN go get -v .
RUN go build -o /volcabot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /volcabot /root/volcabot

WORKDIR /root/
ENV VOLCABOT_LOG_LEVEL=debug
CMD ["./volcabot"]
