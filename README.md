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

### Go

```shell
go install github.com/antony-ramos/guildops/cmd/guildops@latest
guildops ...
```

### Docker

```shell
docker pull ghcr.io/antony-ramos/guildops
docker run --publish 8080 ghcr.io/antony-ramos/guildops
```

## Usage

### Configuration

A configuration file is required to run GuildOps. 

This configuration file should be located at `"config/config.yml"`

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

For exemple :
>
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

- <https://github.com/go-co-op/gocron> - A Golang Job Scheduling Package
- <https://github.com/gorilla/mux> - A powerful HTTP router and URL matcher for building Go web servers with 🦍
- <https://github.com/sirupsen/logrus> - Structured, pluggable logging for Go.
- <https://github.com/spf13/viper> - Go configuration with fangs
  github.com/Masterminds/squirrel v1.5.4
  github.com/bwmarrin/discordgo v0.27.1
  github.com/ilyakaznacheev/cleanenv v1.5.0
  github.com/jackc/pgx/v4 v4.18.1
  github.com/lib/pq v1.10.9
  github.com/prometheus/client_golang v1.16.0
  go.opentelemetry.io/otel v1.18.0
  go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.18.0
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.18.0
  go.opentelemetry.io/otel/sdk v1.18.0
  go.opentelemetry.io/otel/trace v1.18.0
  go.uber.org/zap v1.26.0