#!/usr/bin/env bash

ORG_NAME="alvadorncorp"
PROJECT_NAME="bunny-cli"
VERSION=$1
IMAGE_NAME="${ORG_NAME}/${PROJECT_NAME}:${VERSION}"
PLATFORMS="linux/amd64,linux/arm64,linux/arm/v8,linux/arm/v7"

declare -A XCOMPILE=(
    ["linux"]="arm64 amd64 arm"
    ["darwin"]="arm64"
    ["windows"]="arm64 arm 386 amd64"
)

go-test() {
  go test ./...
}

compile() {
    for os in "${!XCOMPILE[@]}"
    do
        archs="${XCOMPILE[$os]}"
        for arch in ${archs}
        do
            filename="${PROJECT_NAME}-${VERSION}-${os}-${arch}"
            if [ "$os" == "windows" ];then filename+=".exe"; fi
            GOOS=$os GOARCH=$arch go build -o "./build/$filename" ./cmd
            chmod +x "./build/$filename"
        done
    done
}

docker-build() {
    archs="${XCOMPILE["linux"]}"
    local images=()
    for arch in ${archs}
    do
        # docker build --platform=linux/${arch} --build-arg=OS=linux --build-arg=ARCH=$arch --build-arg=VERSION=$VERSION -t ${IMAGE_NAME}-${arch} .
        docker build --no-cache --build-arg=OS=linux --build-arg=ARCH=$arch --build-arg=VERSION=$VERSION -t ${IMAGE_NAME}-${arch} .
        docker push ${IMAGE_NAME}-${arch}
        images+=("${IMAGE_NAME}-${arch}")
    done

    echo ${images[*]}
    echo "docker manifest create ${IMAGE_NAME} ${images[*]}"
    docker manifest create ${IMAGE_NAME} ${images[*]} 
}

docker-push() {
    docker manifest push ${IMAGE_NAME}
}

case "$1" in
  test ) go-test;;
  compile ) "$@";;
  all ) compile; docker-build; docker-push;;
  * ) echo "command does not exist";;
esac
