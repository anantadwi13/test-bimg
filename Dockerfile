FROM ubuntu:18.04

RUN apt update && \
    apt upgrade -y
RUN apt install -y wget git libtool build-essential pkg-config libglib2.0-dev libexpat1-dev \
      libpng-dev libtiff5-dev libwebp-dev libgif-dev libjpeg-turbo8-dev libexif-dev libgsf-1-dev

WORKDIR /tmp

RUN wget https://github.com/libvips/libvips/releases/download/v8.12.2/vips-8.12.2.tar.gz && \
    tar -xvzf vips-8.12.2.tar.gz && \
    cd vips-8.12.2 && \
    ./configure && \
    make -j 8 && \
    make install && \
    ldconfig

RUN wget https://go.dev/dl/go1.18.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

COPY . /root/test-bimg

WORKDIR /root/test-bimg

RUN go mod tidy

CMD go run ./