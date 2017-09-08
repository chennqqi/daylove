FROM golang:1.9

ADD . /go/src/daylove/


WORKDIR /go/src/daylove

RUN go get -u github.com/golang/dep/cmd/dep ; \
    dep ensure -v ; \ 
    go build -buildmode=pie -v 

FROM centos
RUN mkdir /www
COPY --from=0 /go/src/daylove/daylove /www/daylove
WORKDIR /www
CMD ["/www/daylove"]


