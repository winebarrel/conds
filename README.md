# conds

[![CI](https://github.com/winebarrel/conds/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/conds/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/conds.svg)](https://pkg.go.dev/github.com/winebarrel/conds)

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
		C("num = @n", c.NV("n", n)).
		AND_C("str = @s", c.NV("s", s)).
		// XNV: Null value conditions are removed
		AND_C("sval = @ns", c.XNV("ns", nilstr)).
		OR_C(c.IF(nilnum == nil, "sval IS NULL", "sval = @nn"), c.NV("nn", nilnum))

	stmt, params := w.StmtParams()

	fmt.Println(stmt)   //=> "num = @n AND str = @s OR sval IS NULL"
	fmt.Println(params) //=> map[n:100 nn:<nil> s:foo]
}
```
