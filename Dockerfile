FROM alpine:latest
MAINTAINER Fabio Favero Henkes <fabio.favero@gmail.com>

ADD ./ssl/kripto-ssl.crt /kripto-ssl.crt
ADD ./ssl/kripto-ssl.key /kripto-ssl.key
ADD ./cmd/kserver/kserver /kserver
ADD ./docker-entrypoint.sh /entrypoint.sh

RUN apk add --no-cache util-linux \
    && mkdir -p /data/secrets \
    && uuidgen > .krpt \
    && chmod +x entrypoint.sh

ENTRYPOINT /entrypoint.sh
