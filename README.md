<!-- markdownlint-disable MD033 -->
<h1 align="center">
  <img src="assets/logo.png" alt="guildOps logo" width="150" height="150" style="border-radius: 25%">
</h1>

<h4 align="center">GuildOps - Manage your WoW guild with a Discord Bot</h4>

<div align="center">
  <a href="https://github.com/antony-ramos/guildops/issues/new">Report a Bug</a> ·
  <a href="https://github.com/antony-ramos/guildops/issues/new">Request a Feature</a> ·
  <a href="https://github.com/antony-ramos/guildops/discussions">Ask a Question</a>
  <br/>
  <br/>

[![GoReportCard](https://goreportcard.com/badge/github.com/antony-ramos/guildops)](https://goreportcard.com/report/github.com/antony-ramos/guildops)
[![Codecov branch](https://img.shields.io/codecov/c/github/antony-ramos/guildops/main?label=code%20coverage)](https://app.codecov.io/gh/antony-ramos/guildops/tree/main)
[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/antony-ramos/guildops)
<br/>
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/antony-ramos/guildops?logo=go&logoColor=white&logoWidth=20)
[![License](https://img.shields.io/badge/license-CeCILL%202.1-blue?logo=git&logoColor=white&logoWidth=20)](LICENSE)

<a href="#about">About</a> ·
<a href="#install">How to Install?</a> ·
<a href="#exported-metrics">Metrics</a> ·
<a href="#support">Support</a> ·
<a href="#contributing">Contributing</a> ·
<a href="#security">Security</a>

</div>

---
<!-- markdownlint-enable MD033 -->

## About

**GuildOps** provides a way to manage your WoW Guild with a discord Bot.
* Create raids ;
* Assign Loots ;
* Calculate Loot counter ;
* Add notes on players ;
* And more !

## Install

Before anything, you need to edit the config at config/config.yml to add your discord bot token, guild id and postgres credentials
You can also provides them via environment variables. dotenv is required for docker-compose.

### Go

```shell
go install github.com/antony-ramos/guildops/cmd/guildops@latest
guildops ...
```

### Docker


```shell
docker run  ghcr.io/antony-ramos/guildops
```

### Docker-compose

```shell
docker-compose up
```

### Debian
Before anything, you need to edit the config at config/config.yml to add your discord bot token, guild id and postgres credentials.
```shell
make dpkg.build 
make dpkg.install
```

You can remove it with `make dpkg.remove`

You can also install it via ssh. You need to setup SSH_HOST and SSH_USER in your environment.

```shell
export SSH_HOST=yourhost
export SSH_USER=youruser
make ssh.install
```


## Usage

### Configuration

A configuration file is required to run GuildOps. 

This configuration file should be located at `"config/config.yml"` or any file specified by `CONFIG_PATH`

```yaml
app:
  name: 'guildops'
  version: '0.0.1'
  environment: "development"

logger:
  level: debug

metrics:
  port: 2213

postgres:
  pool_max: 10
  conn_attempts: 10
  conn_timeout: 5s
```

> For exemple :
> ```yaml
> app:
>   name: 'guildops'
>   version: '0.0.1'
>   environment: "development"
>
> logger:
>   level: debug
>
> metrics:
>   port: 2213
>
> postgres:
>   pool_max: 10
>   conn_attempts: 10
>   conn_timeout: 5s
>```

### Secret 

Secrets are needed to be set in the environment variables
    
  ```shell
  export DISCORD_TOKEN=yourtoken
  export DISCORD_GUILD_ID=yourguildid
  export PG_URL=yourpostgresurl
  ```

For Debian installation, you must set them in config file.
  ```yaml
app:
    name: 'guildops'
    version: '0.0.1'
    environment: "development"

logger:
    level: debug

metrics:
    port: 2213

discord:
    token: <yourtoken>
    guild_id: <yourguildid>

postgres:
    pool_max: 10
    conn_attempts: 10
    conn_timeout: 5s
    url: <yourpostgresurl>
  ```

## Use Discord Commands

Please read [our usage guide](docs/USAGE.md)


## Support

Reach out to the maintainer at one of the following places:

- [GitHub Discussions](https://github.com/antony-ramos/guildops/discussions)
- Open an issue on [GitHub](https://github.com/antony-ramos/guildops/issues/new)

## Contributing

First off, thanks for taking the time to contribute! Contributions are what make the
open-source community such an amazing place to learn, inspire, and create. Any contributions
you make will benefit everybody else and are **greatly appreciated**.

Please read [our contribution guidelines](docs/CONTRIBUTING.md), and thank you for being involved!

## Security

`guildops` follows good practices of security, but 100% security cannot be assured.
`guildops` is provided **"as is"** without any **warranty**. Use at your own risk.

*For more information and to report security issues, please refer to our [security documentation](docs/SECURITY.md).*

## License

This project is licensed under the **CeCILL License 2.1**.

See [LICENSE](LICENSE) for more information.

## Acknowledgements

Thanks for these awesome resources and projects that were used during development:

- <https://github.com/Masterminds/squirrel> Squirrel - fluent SQL generator for Go
- <https://github.com/bwmarrin/discordgo>  Go bindings for Discord
- <https://github.com/ilyakaznacheev/cleanenv> Clean and minimalistic environment configuration reader for Golang
- <https://github.com/jackc/pgx> PostgreSQL driver and toolkit for Go
- <https://github.com/lib/pq> Pure Go Postgres driver for database/sql
- <https://github.com/prometheus/client_golang> Prometheus instrumentation library for Go applications
- <https://github.com/stretchr/testify> A toolkit with common assertions and mocks that plays nicely with the standard library
- <https://go.opentelemetry.io/otel> OpenTelemetry-Go is the Go implementation of OpenTelemetry
- <https://go.uber.org/zap> Blazing fast, structured, leveled logging in Go
