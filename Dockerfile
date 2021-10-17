# syntax=docker/dockerfile:1

FROM golang:latest as build

WORKDIR /go/src/github.com/nicolastakashi/gitana

RUN apt-get update
RUN useradd -ms /bin/bash gitana

COPY --chown=gitana:gitana . .

RUN make all

FROM golang:latest as runtime

RUN useradd -ms /bin/bash gitana

WORKDIR /gitana

COPY --from=build /go/src/github.com/nicolastakashi/gitana/bin/* /bin/

USER gitana

ENTRYPOINT [ "/bin/gitana" ]