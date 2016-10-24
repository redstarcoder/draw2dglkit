draw2dglkit
======
[![GoDoc](https://godoc.org/github.com/redstarcoder/draw2dglkit?status.svg)](https://godoc.org/github.com/redstarcoder/draw2dglkit)

Package draw2dglkit offers useful tools for using [draw2d](https://github.com/llgcode/draw2d) with OpenGL.

Benchmarks
---------------

```
$ go test -cpu 1 -bench .
BenchmarkIsPointInShape 	   10000	    151902 ns/op
PASS
ok  	github.com/redstarcoder/draw2dglkit	2.349s
```

Installation
---------------

Install [golang](http://golang.org/doc/install). To install or update the package draw2dglkit on your system, run:

```
go get -u github.com/redstarcoder/draw2dglkit
```

Acknowledgments
---------------

[redstarcoder](https://github.com/redstarcoder) wrote this library.
[Laurent Le Goff](https://github.com/llgcode) wrote draw2d.

