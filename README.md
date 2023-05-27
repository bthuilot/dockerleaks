# Dockerleaks

[![Go Report Card](https://goreportcard.com/badge/github.com/bthuilot/dockerleaks)](https://goreportcard.com/report/github.com/bthuilot/dockerleaks)
[![License](https://img.shields.io/github/license/bthuilot/dockerleaks)](./LICENSE)

Dockerleaks is a command-line tool designed to uncover secrets within Docker images.
Secrets, which can include API keys, passwords, and access tokens, pose potential security risks if left exposed. 
Dockerleaks will comb through Docker images to help identify these secrets.

You can scan Docker images either located in a remote registry or stored locally.
It uses different methods to investigate environment variables and build arguments,
which are common places where secrets might inadvertently be embedded during image creation.
Furthermore, dockerSecrets can dig deep into the filesystem within a Docker image,
scanning through files for potential secret leaks.

By using this tool, you can ensure that your Docker images maintain their integrity,
adhering to the best practices of sensitive information management.
It's a valuable addition to any security-conscious developer or organization's toolset,
assisting in preventing unauthorized access to critical services, databases, and other resources.

© 2023 Bryce Thuilot. Dockerleaks is an open-source project and comes with ABSOLUTELY NO WARRANTY.
It is free software, and you are welcome to redistribute it under specific conditions.

## Installation

### via Homebrew

```shell
brew install bthuilot/tap/dockerleaks
```

### via GitHub release

Navigate to the Releases tab to download the compile binary for your specific platform.
Be sure to then add it to your shell's `PATH`

### from source

```shell
git clone github.com/bthuilot/dockerleaks && cd dockerleaks
go build -o dockerleaks .
# add dockerleaks to your PATH or execute via ./dockerleaks
```

### via Docker

A docker image containing the script is distributed via `thuilot/dockerleaks`. To run via docker be sure to mount
the docker socket into the container such that the binary can connect to the daemon to perform scans, an example
is shown below.

```shell
docker run -it -v /var/run/docker.sock:/var/run/docker.sock thuilot/dockerleaks:latest -i [IMAGE TO SCAN] 
```

## Configuration

The application can be configured via a file named `dockerleaks.yml` located in the same directory the
tool is run from, the directory `$HOME/.dockerleaks`, or the folder `/etc/dockerleaks`.

Checkout the file [`dockerleaks.example.yml`](/dockerleaks.example.yml) located in the root of this
repository for more information

## Usage

The tool can be used to scan both remote and local built docker images.
For example, to scan a remote image named `my-image`, you could use the following command:

```commandline
dockerleaks -i my-image:latest -p
```

This command would pull `my-image:latest` from its remote source and scan it for leaked secrets.


## Support the project

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/thuilot)