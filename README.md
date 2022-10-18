# happendb

[![madeby](https://img.shields.io/badge/made%20by-%40drgomesp-blue)](https://github.com/drgomesp/)
[![Go Report Card](https://goreportcard.com/badge/github.com/drgomesp/happendb)](https://goreportcard.com/report/github.com/drgomesp/happendb)
[![build](https://github.com/drgomesp/happendb/actions/workflows/go-test.yml/badge.svg?style=squared)](https://github.com/drgomesp/happendb/actions)
[![codecov](https://codecov.io/gh/drgomesp/happendb/branch/main/graph/badge.svg?token=BRMFJRJV2X)](https://codecov.io/gh/drgomesp/happendb)

> A decentralized event-sourcing platform.

## Table of Contents

- [Install](#install)
- [Features](#features)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

```bash
go get github.com/drgomesp/happendb
```

## Building & Testing

This project uses [Taskfile](https://taskfile.dev/), which is easier than something like GNU Make.

To list available tasks, run:

```bash
$ task --list-all
```

## Running

- The easiest way is to use docker-compose:

```bash
$ docker-compose up
```

- Or, if you have a Tendermint instance running locally, you can also run the binary or the Go command directly:

```bash
# assuming tendermint rpc is avaialble at tcp://127.0.0.1:26658
$ go run cmd/happendb -socket-addr=tcp://127.0.0.1:26658
```

## Contributing

PRs accepted.

## License

MIT © [Daniel Ribeiro](https://github.com/drgomesp)
