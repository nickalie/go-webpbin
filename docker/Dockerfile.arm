FROM resin/raspberry-pi3-alpine-golang:slim

RUN apk add --no-cache --update libpng-dev libjpeg-turbo-dev giflib-dev tiff-dev autoconf automake make g++ wget git

WORKDIR /

RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
	tar -xvzf libwebp-0.6.0.tar.gz && \
	mv libwebp-0.6.0 libwebp && \
	rm libwebp-0.6.0.tar.gz && \
    cd /libwebp && \
	./configure && \
	make && \
	make install && \
	cd / && rm -rf libwebp

RUN go get -v github.com/nickalie/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v golang.org/x/image/webp

RUN mkdir -p $GOPATH/src/github.com/nickalie/go-webpbin
COPY . $GOPATH/src/github.com/nickalie/go-webpbin
WORKDIR $GOPATH/src/github.com/nickalie/go-webpbin
RUN go test -v ./...