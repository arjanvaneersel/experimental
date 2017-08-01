FROM golang:1.8 as builder
RUN apt-get update && apt-get install -y flite-dev && mkdir -p /go/src/app
COPY main.go /go/src/app
COPY vendor /go/src/vendor
RUN go install app

FROM debian:jessie-slim
RUN apt-get update && apt-get install -y flite-dev
COPY --from=builder /go/bin/app /bin/app
ENTRYPOINT ["/bin/app"]