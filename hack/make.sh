#!/bin/bash

echo 'hack 4 good... _\,,/'
echo $GOPATH

SRC=src/$REPOSITORY/$GROUP/$PROJECT/

rm -rf $GOPATH/$SRC/
echo 'removed: '$GOPATH/$SRC

shopt -s extglob

mkdir -p $GOPATH/$SRC && mv !(hack*) $GOPATH/$SRC

shopt -u extglob

cd $GOPATH/$SRC

make deps

echo 'calling make: '$COMMAND
make $COMMAND
