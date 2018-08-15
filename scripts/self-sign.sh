#!/bin/bash

rm -rf ./ssl
mkdir -p ./ssl

echo "Generating self signed certificate..."

openssl req \
    -new \
    -newkey rsa:4096 \
    -days 365 \
    -nodes \
    -x509 \
    -subj "/C=BR/ST=None/L=Kripto/O=None/CN=www.kripto.local" \
    -keyout kripto-ssl.key \
    -out kripto-ssl.crt

mv -f kripto-ssl.* ./ssl/

echo "Done!"
