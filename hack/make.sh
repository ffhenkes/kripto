#!/bin/bash

echo 'hack 4 good... _\,,/'
echo $GOPATH

SRC=src/$REPOSITORY/$GROUP/$PROJECT/

rm -rf $GOPATH/$SRC/
echo 'removed: '$GOPATH/$SRC

shopt -s extglob dotglob

mkdir -p $GOPATH/$SRC && mv !(hack*) $GOPATH/$SRC

shopt -u extglob dotglob

cd $GOPATH/$SRC

echo 'calling make: '$COMMAND
make $COMMAND
