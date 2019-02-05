# Ech0

Ech0 (pronounced "echo zero") is a logging adapter for `echo.Logger` that uses [github.com/rs/zerolog](https://github.com/rs/zerolog) as the logging backend instead of the default [github.com/labstack/gommon/log](https://github.com/labstack/gommon/tree/master/log)

# Why?

1. I like [Echo](https://echo.labstack.com/).
1. I like [zerolog](https://github.com/rs/zerolog).
1. I wanted to have *one* logging backend in my Echo apps.

# Installing
`go get -u github.com/dvelitchkov/ech0`

# Benchmarks

```
goos: darwin
goarch: amd64
pkg: github.com/dvelitchkov/ech0
BenchmarkZeroFormat-8     	 1000000	      1915 ns/op	     533 B/op	       3 allocs/op
BenchmarkZeroJSON-8       	 1000000	      2249 ns/op	    1032 B/op	       3 allocs/op
BenchmarkZero-8           	 1000000	      1893 ns/op	     533 B/op	       3 allocs/op
BenchmarkGommonFormat-8   	 1000000	      2203 ns/op	     256 B/op	      11 allocs/op
BenchmarkGommonJSON-8     	  500000	      2695 ns/op	     472 B/op	      14 allocs/op
BenchmarkGommon-8         	 1000000	      2136 ns/op	     256 B/op	      11 allocs/op
PASS
ok  	github.com/dvelitchkov/ech0	11.903s

```
