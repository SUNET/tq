FROM golang:latest as build
RUN mkdir -p /go/src/github.com/SUNET/tq
ADD . /go/src/github.com/SUNET/tq/
WORKDIR /go/src/github.com/SUNET/tq
RUN go build ./cmd/tq
RUN env GOBIN=/usr/bin go install ./cmd/tq

# Now copy it into our base image.
FROM gcr.io/distroless/base:debug
COPY --from=build /usr/bin/tq /usr/bin/tq

ENTRYPOINT ["/usr/bin/tq"]
