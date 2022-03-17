# oxeye

[![Build Status](https://github.com/ectobit/oxeye/workflows/build/badge.svg)](https://github.com/ectobit/oxeye/actions)
[![Go Reference](https://pkg.go.dev/badge/go.ectobit.com/oxeye.svg)](https://pkg.go.dev/go.ectobit.com/oxeye)
[![Go Report](https://goreportcard.com/badge/go.ectobit.com/oxeye)](https://goreportcard.com/report/go.ectobit.com/oxeye)
[![License](https://img.shields.io/badge/license-BSD--2--Clause--Patent-orange.svg)](https://github.com/ectobit/oxeye/blob/main/LICENSE)

Concurrent microservice in Go using message broker to receive jobs and send results.

## Contribution

- `make lint` lints the project
- `make start` starts docker-compose stack
- `make stop` stops docker-compose stack
- `make test` runs unit tests
- `make test-cov` displays test coverage (requires docker-stack to be up)
