```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/gnvk/gorm_bench
BenchmarkSqlQuery-8          124           9630147 ns/op
BenchmarkPgxQuery-8          136           9302955 ns/op
BenchmarkGormQuery-8          48          25418400 ns/op
PASS
```
