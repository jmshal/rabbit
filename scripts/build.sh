#!/usr/bin/env bash

cd "$(dirname "$0")/.."

if [ -d ./bin ]; then
    echo "Cleaning bin directory..."
    rm -rf ./bin/*
fi

if [ -d ./vendor ]; then
    echo "Cleaning vendor directory..."
    rm -rf ./vendor/*/
fi

echo "Installing dependencies..."
govendor sync

echo "Detecting version..."
VERSION=v`go run ./cmd/rabbit/*.go version`

for GOOS in darwin windows linux; do
  for GOARCH in 386 amd64; do
    echo "Building $GOOS/$GOARCH..."
    FILENAME="rabbit"
    if [ "windows" == $GOOS ]; then
        FILENAME="rabbit.exe"
    fi
    GOOS=$GOOS GOARCH=$GOARCH go build -o ./bin/$GOOS/$GOARCH/$FILENAME ./cmd/rabbit
    if [ "windows" == $GOOS ]; then
        zip -rjX ./bin/rabbit-$VERSION-$GOOS-$GOARCH.zip ./bin/$GOOS/$GOARCH/
    else
        tar -C ./bin/$GOOS/$GOARCH/ -cvzf ./bin/rabbit-$VERSION-$GOOS-$GOARCH.tar.gz .
    fi
  done
done

echo "Building linux/arm..."
GOOS=linux GOARCH=arm go build -o ./bin/linux/arm/rabbit ./cmd/rabbit
tar -C ./bin/linux/arm/ -cvzf ./bin/rabbit-$VERSION-linux-arm.tar.gz .
