FROM golang:1.16 as build

ENV CGO_ENABLED=0

RUN go get -u github.com/mitchellh/gox

WORKDIR /go/src/

COPY . .

RUN echo "==> Building..." && \
    gox -output="pkg/{{.OS}}_{{.Arch}}/kubectl-druid" \
        -os="darwin linux" \
        -arch="amd64"
