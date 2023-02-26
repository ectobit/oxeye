# oxeye

[![Build Status](https://github.com/ectobit/oxeye/workflows/build/badge.svg)](https://github.com/ectobit/oxeye/actions)
[![Go Reference](https://pkg.go.dev/badge/go.ectobit.com/oxeye.svg)](https://pkg.go.dev/go.ectobit.com/oxeye)
[![Go Report](https://goreportcard.com/badge/go.ectobit.com/oxeye)](https://goreportcard.com/report/go.ectobit.com/oxeye)

Concurrent microservice in Go using message broker to receive jobs and send results.

## Contribution

- `make lint` lints the project
- `make start` starts docker-compose stack
- `make stop` stops docker-compose stack
- `make test` runs unit tests
- `make test-cov` displays test coverage (requires docker-stack to be up)

## License

Licensed under either of

- Apache License, Version 2.0
  ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT license
  ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in the work by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
