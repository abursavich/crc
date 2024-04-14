# CRC
[![License][license-img]][license]
[![GoDev Reference][godev-img]][godev]
[![Go Report Card][goreportcard-img]][goreportcard]

This module wraps standard library CRC packages and provides additional functionality to efficiently
combine independently calculated checksums of sequential blocks of data.

The algorithm for combining checksums is adapted from [zlib] by Mark Adler.


[license]: https://raw.githubusercontent.com/abursavich/crcx/main/LICENSE
[license-img]: https://img.shields.io/badge/license-mit-blue.svg?style=for-the-badge

[godev]: https://pkg.go.dev/bursavich.dev/crc
[godev-img]: https://img.shields.io/static/v1?logo=go&logoColor=white&color=00ADD8&label=dev&message=reference&style=for-the-badge

[goreportcard]: https://goreportcard.com/report/bursavich.dev/crc
[goreportcard-img]: https://goreportcard.com/badge/bursavich.dev/crc?style=for-the-badge

[zlib]: https://github.com/madler/zlib
