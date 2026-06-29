module github.com/apalala/jb

go 1.26.4

replace github.com/apalala/jb => .

replace github.com/apalala/jb/pkg => ./pkg

require (
	github.com/alecthomas/kong v1.15.0
	github.com/mitchellh/go-ps v1.0.0
	github.com/tortxof/z85 v0.1.0
	golang.org/x/sys v0.46.0
	gonum.org/v1/gonum v0.17.0
)
