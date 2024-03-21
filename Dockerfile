FROM ubuntu:20.04
RUN rm /bin/sh && ln -s /bin/bash /bin/sh
RUN mkdir /app
WORKDIR /app

RUN apt-get update && apt-get upgrade -y && apt install -y git make gcc gcc-arm-linux-gnueabi
COPY --from=golang:1.21-bullseye /usr/local/go/ /usr/local/go/
RUN echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
