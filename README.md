# Go Packages
[![Build Status](https://img.shields.io/travis/razonyang/gopkgs?style=flat-square)](https://travis-ci.org/razonyang/gopkgs)

A self-host HTTP service that allow customizing your Go package import paths.

## Live

I launch up an HTTP service(https://pkg.clevergo.tech) for meeting my own needs. With this service, you only need a domain name to customize the import path of your package.

The rest of sections introduces how to host on your own server.

## Requirements

- MySQL/MariaDB.
- Redis.
- [Auth0](https://auth0.com/) Application.

## Configuration

```shell
$ cp .env.example .env
```

Checkout the [.env.example] for details.

## Migration

```shell
$ migrate --database mysql://user:password@tcp(host:port)/dbname?query -path=/migrations up
```

Checkout [go-migrate](https://github.com/golang-migrate/migrate) for details.

## Start service

```shell
$ go run main.go serve
```
