#!/usr/bin/env sh
set -eu

platform="${1:-linux/arm64}"
image="${BENZHI_IMAGE:-url-short-link:benzhi}"

docker buildx build --load --platform "$platform" -f benzhi.Dockerfile -t "$image" .
docker run --rm --platform "$platform" --entrypoint go "$image" build ./...
docker run --rm --platform "$platform" "$image"
