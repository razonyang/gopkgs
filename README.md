# GOPKGs
[![Build Status](https://img.shields.io/travis/razonyang/gopkgs?style=flat-square)](https://travis-ci.org/razonyang/gopkgs)

`gopkgs` is a CLI application and HTTP service that manage custom import path of your Go packages.

The major advantages of using custom import path is that you don't need to change the import path when you are going to host your code elsewhere.

## Online preview

```shell
$ curl "https://clevergo.tech/clevergo?go-get=1"
$ curl "https://clevergo.tech/captchas?go-get=1"
$ curl "https://clevergo.tech/form?go-get=1"
```

## Installation

### Source

```shell
$ git clone https://github.com/razonyang/gopkgs.git
$ cd gopkgs
$ packr2 build --tag [database]
```

- `database`: `mysql`, `sqlite3` or `postgres`.

> `go get -u github.com/gobuffalo/packr/v2/packr2` for installing `packr2`.

### Binary

Checkout [releases](https://github.com/razonyang/gopkgs/releases) page and download.

## Configuration

Configuration is a JSON file.

```json
{
    "addr": ":8080",
    "db": {
        "driver": "sqlite3",
        "dsn": "gopkgs.db",
        "tableName": "packages"
    }
}
```

- `addr`: HTTP server address.
- `db`:
    - `dsn`: data source name, depends on what driver you use.
        - `sqlite3`: `/path/to/gopkgs.db`
        - `mysql`: `user:password@tcp(localhost:3306)/gopkgs?charset=utf8mb4&parseTime=True&loc=Local`
        - `postgres`: `postgres://user:password@localhost/gopkgs?sslmode=verify-full`
    - `tableName`: the name of packages table.

The `config.json` of the current directory will be used by default, you can specify the configuration file via `-c` or `--config` flag:

```shell
$ gopkgs -c /etc/gopkgs/config.json
```

## Usage

### Start HTTP server

```shell
$ gopkgs serve
```

> You can use [supervisord](http://supervisord.org/) to manage gopkgs HTTP service.

You may also need to set up a reverse proxy, let's take Nginx as example:

```nginx
location / {
    try_files $uri $uri/ @gopkgsproxy;
}

location @gopkgsproxy {
    proxy_set_header Host $host;
    proxy_pass http://127.0.0.1:8080;
}  
```

### Add package

```shell
$ gopkgs add <prefix> <vcs> <repo-root> [<docs-url>]
```

- `prefix`: the prefix of import path.
- `vcs`: bzr, fossil, git, hg, svn.
- `repo-root`: the location of your repository.
- `docs-url`: URL of documentations, optional.

```shell
$ gopkgs add \
    example.com/foo \
    git \
    https://github.com/example/foo
```

And then checkout the output.

```shell
$ curl https://example.com/foo
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="example.com/foo git https://github.com/example/foo">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/example.com/foo?tab=doc">
<title>Package example.com/foo</title>
</head>
<body>
Nothing to see here; <a href="https://pkg.go.dev/example.com/foo?tab=doc">move along</a>.
</body>
</html>
```

### Show package

```shell
$ gopkgs show example.com/foo
example.com/foo
vcs : git
root: https://github.com/example/foo
docs: https://pkg.go.dev/example.com/foo?tab=doc
```

### Edit package

**Modify VCS**

```shell
$ gopkgs set-vcs example.com/foo svn
```

**Repository Root**

```shell
$ gopkgs set-root example.com/foo https://gitlab.com/example/foo
```

**Documentations**

```shell
$ gopkgs set-docs example.com/foo https://docs.example.com/foo
```

**Verify that everything is OK**

```shell
$ gopkgs show example.com/foo
example.com/foo
vcs : svn
root: https://gitlab.com/example/foo
docs: https://docs.example.com/foo
```

### Remove pacakge

```shell
$ gopkgs remove example.com/foo
```

### Help

```shell
$ gopkgs help
```