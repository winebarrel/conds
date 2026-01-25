# conds

[![CI](https://github.com/winebarrel/conds/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/conds/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/conds.svg)](https://pkg.go.dev/github.com/winebarrel/conds)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/winebarrel/conds)](https://pkg.go.dev/github.com/winebarrel/conds?tab=versions)

cond is a tiny builder of where clause conditions.

## Installation

```sh
go get github.com/winebarrel/conds
```

## Usage

```go
package main

import (
	"fmt"

	c "github.com/winebarrel/conds"
)

func main() {
	n := 100
	s := "foo"
	var nilstr *string
	var nilnum *int

	w := c.
		C("num = @n", c.V("n", n)).
		AND_C("str = @s", c.V("s", s)).
		// XV: Null value conditions are removed
		AND_C("sval = @ns", c.XV("ns", nilstr)).

	stmt, params := w.StmtParams()

	fmt.Println(stmt)   //=> "num = @n AND str = @s"
	fmt.Println(params) //=> map[n:100 nn:<nil> s:foo]
}
```
