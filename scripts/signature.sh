#!/bin/bash

RSA_DIRECTORY=/data/rsa

mkdir -p ${RSA_DIRECTORY}

echo "Generating an RSA private key..."
openssl genrsa -out ${RSA_DIRECTORY}/kripto.rsa 1024

echo "Generating an RSA public key ..."
openssl rsa -in ${RSA_DIRECTORY}/kripto.rsa -pubout > ${RSA_DIRECTORY}/kripto.rsa.pub
