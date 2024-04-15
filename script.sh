#!/bin/bash

version="2.1.1"

build() {
    CGO_ENABLED=0  GOOS=darwin   GOARCH=amd64 go build -ldflags "-s -w" -o ./deploy/MagicBox_amd64_darwin .
    CGO_ENABLED=0  GOOS=windows  GOARCH=amd64 go build -ldflags "-s -w" -o ./deploy/MagicBox_amd64_win.exe .
    CGO_ENABLED=0  GOOS=linux    GOARCH=amd64 go build -ldflags "-s -w" -o ./deploy/MagicBox_amd64_linux .
    CGO_ENABLED=0  GOOS=darwin   GOARCH=arm64 go build -ldflags "-s -w" -o ./deploy/MagicBox_arm64_darwin .
    CGO_ENABLED=0  GOOS=windows  GOARCH=arm64 go build -ldflags "-s -w" -o ./deploy/MagicBox_arm64_win.exe .
    CGO_ENABLED=0  GOOS=linux    GOARCH=arm64 go build -ldflags "-s -w" -o ./deploy/MagicBox_arm64_linux .
}

bpush() {
    docker buildx build -t srcrs/magicbox:$version --platform=linux/arm64,linux/amd64 . --push
    docker buildx build -t srcrs/magicbox:latest --platform=linux/arm64,linux/amd64 . --push
}

test() {
    docker build -t srcrs/magicbox:$version .
    docker tag srcrs/magicbox:$version srcrs/magicbox:latest
}

# 调用指定的函数
for func in "$@"
do
    $func
done

exit 0