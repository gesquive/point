# [gesquive/reflect](https://github.com/gesquive/reflect)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/reflect/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/reflect)
[![Pipeline Status](https://img.shields.io/gitlab/pipeline/gesquive/reflect?style=flat-square)](https://gitlab.com/gesquive/reflect/pipelines)
[![Coverage Report](https://gitlab.com/gesquive/reflect/badges/master/coverage.svg?style=flat-square)](https://gesquive.gitlab.io/reflect/coverage.html)
[![Github Release](https://img.shields.io/github/v/tag/gesquive/reflect?style=flat-square)](https://github.com/gesquive/reflect)

# Supported Architectures

This image supports multiple architectures:

- `amd64`, `x86-64`
- `armv7`, `armhf`
- `arm64`, `aarch64`

Docker images are uploaded with using Docker manifest lists to make multi-platform deployments easer. More info can be found from [Docker](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md#manifest-list)

You can simply pull the image using `gesquive/reflect` and docker should retreive the correct image for your architecture.

# Supported Tags
If you want a specific version of `reflect` you can pull it by specifying a version tag.

## Version Tags
This image provides versions that are available via tags. 

| Tag    | Description |
| ------ | ----------- |
| `latest` | Latest stable release |
| `0.9.0`  | Stable release v0.9.0 |
| `0.9.0-<git_hash>` | Development preview of version v0.9.0 at the given git hash |

# Usage

Here are some example snippets to help you get started creating a container.
docker

## Docker CLI

```shell
docker create \
  --name=reflect \
  -p 2626:2626 \
  -v path/to/config:/config \
  --restart unless-stopped \
  gesquive/reflect
```

## Docker Compose
Compatible with docker-compose v2 schemas.

```docker
---
version: "2"
services:
  reflect:
    image: gesquive/reflect
    container_name: reflect
    volumes:
      - path/to/config:/config
    ports:
      - 2626:2626
    restart: unless-stopped
```
# Parameters
The container defines the following parameters that you can set:

| Parameter | Function |
| --------- | -------- |
| `-p 2626`          | The port for the reflect REST API |
| `-v /etc/reflect`  | The reflect config goes here |
