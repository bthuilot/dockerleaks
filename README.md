# Dockerleaks

[![Go Report Card](https://goreportcard.com/badge/github.com/bthuilot/dockerleaks)](https://goreportcard.com/report/github.com/bthuilot/dockerleaks)
[![License](https://img.shields.io/github/license/bthuilot/dockerleaks)](./LICENSE.md)

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

## Usage

To run DockerSecrets, use the following command:

Here are the available flags:

* `-c, --config string`: Path to the configuration file.
* `-h, --help`: Show help for `dockerleaks`.
* `-i, --image string`: The name of the Docker image to scan.
* `-l, --log-level string`: Set the log level. Options are `off`, `debug`, `info`, `warn`, `error`, and `fatal`. The default is `off`.
* `-p, --pull`: The image should be pulled from a remote source before scanning.

For example, to scan a remote image named `my-image`, you could use the following command:

```commandline
dockerleaks -i my-image:latest -p
```

This command would pull `my-image:latest` from its remote source and scan it for leaked secrets.