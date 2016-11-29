# timeout
[![Build Status](https://travis-ci.org/go-http-utils/timeout.svg?branch=master)](https://travis-ci.org/go-http-utils/timeout)
[![Coverage Status](https://coveralls.io/repos/github/go-http-utils/timeout/badge.svg?branch=master)](https://coveralls.io/github/go-http-utils/timeout?branch=master)

HTTP timeout middleware for Go.

## Installation

```go
go get -u github.com/go-http-utils/timeout
```

## Documentation

API documentation can be found here: https://godoc.org/github.com/go-http-utils/timeout

## Usage

```go
import (
  "github.com/go-http-utils/timeout"
)
```

```go
mux := http.NewServeMux()
mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
  res.Write([]byte("Hello World"))
})

http.ListenAndServe(":8080", timeout.Handler(mux, time.Second*10, timeout.DefaultTimeoutHandler))
```
