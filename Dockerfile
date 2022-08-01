FROM ubuntu:20.04
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN mkdir /app
WORKDIR /app

RUN apt-get update && apt-get dist-upgrade -y
RUN apt-get install -y make wget gcc gcc-arm-linux-gnueabi
RUN wget -O go.tgz http://golang.org/dl/go1.17.9.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go.tgz \
    && echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc \
