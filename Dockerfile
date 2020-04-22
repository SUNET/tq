FROM golang:latest as build
RUN mkdir -p /go/src/github.com/sunet/tq
ADD . /go/src/github.com/sunet/tq/
WORKDIR /go/src/github.com/sunet/tq
RUN make
RUN env GOBIN=/usr/bin go install ./cmd/tq

# Now copy it into our base image.
FROM gcr.io/distroless/base:debug
COPY --from=build /usr/bin/tq /usr/bin/tq

ENTRYPOINT ["/usr/bin/tq"]
