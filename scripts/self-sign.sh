#!/bin/bash

mkdir -p ./ssl

echo "Generating an SSL private key to sign your certificate..."
openssl genrsa -des3 -out kripto-ssl.key 1024

echo "Generating a Certificate Signing Request..."
openssl req -new -key kripto-ssl.key -out kripto-ssl.csr

echo "Removing passphrase from key (for nginx)..."
cp kripto-ssl.key kripto-ssl.key.org
openssl rsa -in kripto-ssl.key.org -out kripto-ssl.key
rm kripto-ssl.key.org

echo "Generating certificate..."
openssl x509 -req -days 365 -in kripto-ssl.csr -signkey kripto-ssl.key -out kripto-ssl.crt

mv -f kripto-ssl.* ./ssl/
