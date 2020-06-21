# GOPKGs
[![Build Status](https://img.shields.io/travis/razonyang/gopkgs?style=flat-square)](https://travis-ci.org/razonyang/gopkgs)

`gopkgs` is a CLI application and HTTP service that manage custom import path of your Go packages.

## Installation

### Source

```shell
$ git clone https://github.com/razonyang/gopkgs.git
$ cd gopkgs
$ go install
```

Rebuild plugins

```shell
$ cd plugins
$ go build --buildmode=plugin mysql/mysql.go
$ go build --buildmode=plugin sqlite3/sqlite3.go
$ go build --buildmode=plugin postgres/postgres.go
```

### Binary

Checkout [releases](https://github.com/razonyang/gopkgs/releases) page and download.

## Usage

### Configuration

Configuration is a JSON file.

```json
{
    "addr": ":8080",
    "plugins": "plugins",
    "db": {
        "driver": "sqlite3",
        "dsn": "gopkgs.db",
        "tableName": "packages"
    }
}
```

- `addr`: HTTP server address.
- `plugins`: the location of plugins.
- `db`:
    - `driver`: database driver: sqlite3, mysql, postgres. You need to download corresponding plugin from [releases](https://github.com/razonyang/gopkgs/releases), and put it in the `plugins` directory.
    - `dsn`: data source name, depends on what driver you use.
    - `tableName`: the name of packages table.

The `config.json` of the current directory will be used by default, you can specify the configuration file via `-c` or `--config` flag:

```shell
$ gopkgs -c /etc/gopkgs/config.json
```

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
$ gopkgs <prefix> <vcs> <repo-root> [<docs-url>]
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