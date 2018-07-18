#!/bin/sh

set -e

export PHRASE=$(cat .krpt)

./kserver
