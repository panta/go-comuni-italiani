# go-comuni-italiani

## Description

This is a simple Go library and command-line utility to help working with Italian
administrative territorial entities (towns, cities, regions, ...).

It provides:

* a struct `Comune` to represent the most relevant properties
* a pre-generated array `Comuni` (it's a slice of `Comune` structs) containing up-to-date records already converted to Go
* and an utility to download [the official CSV](https://www.istat.it/it/archivio/6789) with ISTAT codes
and convert to an easily (de-)serializable JSON format (or to Go source directly)


## Building

```shell
$ go build -o go-comuni-italiani cmd/go-comuni/main.go
```

Otherwise you can also use the provided `Makefile` (which will create the binary `bin/github.com/panta/go-comuni-italiani`):

```shell
$ make
```

### Downloading and converting the official CSV to JSON 

```shell
$ ./bin/github.com/panta/go-comuni-italiani convert -allow-insecure -output comuni.json
```

(please note that at this time it seems the `-allow-insecure` flag seems to be needed on macOS)

You can then unmarshal the generated json into an array of `Comune`. See the code in `example/example.go`.
