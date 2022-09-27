# syntax=docker/dockerfile:1

FROM golang:latest as build

WORKDIR /go/src/github.com/nicolastakashi/gitana

RUN apt-get update
RUN useradd -ms /bin/bash gitana

COPY --chown=gitana:gitana . .

RUN make all

FROM gcr.io/distroless/static:latest-amd64

WORKDIR /gitana

COPY --from=build /go/src/github.com/nicolastakashi/gitana/bin/* /bin/

USER nobody

ENTRYPOINT [ "/bin/gitana" ]