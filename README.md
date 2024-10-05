[![Go Reference](https://pkg.go.dev/badge/go.abhg.dev/testing/stub.svg)](https://pkg.go.dev/go.abhg.dev/testing/stub)
[![CI](https://github.com/abhinav/stub-go/actions/workflows/ci.yml/badge.svg)](https://github.com/abhinav/stub-go/actions/workflows/ci.yml)

stub is a simple package for stubbing values in tests.
It reassigns a variable or a function reference for the duration of a test,
and restores the original value when the test completes.

Idiomatic usage looks like this:

```go
func TestFoo(t *testing.T) {
    defer stub.Value(&someGlobal, newValue)()

    // test code that uses someGlobal
}
```

See [API Reference](https://abhinav.github.io/stub-go/) for more information.

## History

This package was originally developed as part of [git-spice](https://abhinav.github.io/git-spice/).
It has been extracted into a separate library for reuse in other projects.

The package is inspired by [prashantv/gostub](https://github.com/prashantv/gostub).
The chief difference is that `stub` provides a smaller API surface.

## License

This software is made available under the BSD-3 license.
See LICENSE for the full license text.
