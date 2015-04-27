FROM google/golang

RUN mkdir -p /gopath/src/github.com/etcinit/imageproxy
WORKDIR /gopath/src/github.com/etcinit/imageproxy
ADD . /gopath/src/github.com/etcinit/imageproxy
RUN go get github.com/etcinit/imageproxy/...

RUN go install github.com/etcinit/imageproxy/cmd/imageproxy

CMD []
ENTRYPOINT ["/gopath/bin/imageproxy"]
