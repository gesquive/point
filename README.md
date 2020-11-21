# reflect
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/reflect/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/gesquive/reflect)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/reflect?style=flat-square)](https://circleci.com/gh/gesquive/reflect)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/reflect?style=flat-square)](https://codecov.io/gh/gesquive/reflect)
[![Docker Pulls](https://img.shields.io/docker/pulls/gesquive/reflect?style=flat-square)](https://hub.docker.com/r/gesquive/reflect)

A web client info server.

It provides endpoints to see web client information like: public IP, proxy list, user agent, headers

## Installing

### Compile
This project has only been tested with go1.11+. To compile just run `go get -u github.com/gesquive/reflect` and the executable should be built for you automatically in your `$GOPATH`. This project uses go mods, so you might need to set `GO111MODULE=on` in order for `go get` to complete properly.

Optionally you can clone the repo and run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/reflect/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/reflect`.

### Docker
You can also run reflect from the provided [Docker image](https://hub.docker.com/r/gesquive/reflect) by providing a configuration file:

```shell
docker run -d -p 2626:2626 -v $PWD/docker:/config gesquive/reflect:latest
```

For more details read the [Docker image documentation](https://hub.docker.com/r/gesquive/reflect).

## Configuration

### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./config.yml`
 - `~/.config/reflect/config.yml`
 - `/etc/reflect/config.yml`

Copy `pkg/config.example.yml` to one of these locations and populate the values with your own. Since the config contains a writable API token, make sure to set permissions on the config file appropriately so others cannot read it. A good suggestion is `chmod 600 /path/to/config.yml`.

If you are planning to run this app as a service, it is recommended that you place the config in `/etc/reflect/config.yml`.

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix "REFLECT_" in front of the uppercased variable name. For example, the config variable `log-file` would be the environment variable `REFLECT_LOG_FILE`.

### Service
This application was developed to run as a service behind a webserver such as nginx, apache, or caddy.

You can use upstart, init, runit or any other service manager to run the `reflect` executable. Example scripts for systemd and upstart can be found in the `pkg/services` directory. A logrotate script can also be found in the `pkg/services` directory. All of the configs assume the user to run as is named `reflect`, make sure to change this if needed.

## Usage

```console
A web API for client browser information

Usage:
  reflect [flags]

Flags:
  -a, --address string    The IP address to bind the web server too (default "0.0.0.0")
      --config string     Path to a specific config file (default "./config.yml")
  -l, --log-file string   Path to log file (default "/var/log/reflect.log")
  -p, --port int          The port to bind the webserver too (default 2626)
      --version           Display the version info and exit
```

Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

## Endpoints
There are currently three endpoints:
- `/ip` - show public IP and proxy information
- `/headers` - show all headers
- `/agent` - show the user agent string

Each endpoint accepts requests with a `Content-Type` of json or text.

## Documentation

This documentation can be found at github.com/gesquive/reflect

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
