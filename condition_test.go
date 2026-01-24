package conds_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/conds"
)

func TestConditionEmpty(t *testing.T) {
	c := conds.Empty()
	assert.Equal(t, conds.Condition{}, c)
	assert.True(t, c.Empty())
}

func TestConditionEnclose(t *testing.T) {
	c := conds.C("foo = @bar", conds.NV("bar", 100))
	stmt, params := c.Enclose().StmtParams()
	assert.Equal(t, "(foo = @bar)", stmt)
	assert.Equal(t, map[string]any{"bar": 100}, params)
}

func TestIF(t *testing.T) {
	cthen := conds.C("foo = @bar", conds.NV("bar", 100))
	celse := conds.C("zoo = @baz", conds.NV("baz", "FOO"))
	assert.Equal(t, cthen, conds.IF(true, cthen, celse))
	assert.Equal(t, celse, conds.IF(false, cthen, celse))
}

func TestIFTHEN(t *testing.T) {
	cthen := conds.C("foo = @bar", conds.NV("bar", 100))
	celse := conds.Empty()
	assert.Equal(t, cthen, conds.IFTHEN(true, cthen))
	assert.Equal(t, celse, conds.IFTHEN(false, cthen))
}

func TestIFF(t *testing.T) {
	cthen := conds.C("foo = @bar", conds.NV("bar", 100))
	celse := conds.Empty()
	assert.Equal(t, cthen, conds.IFF(true, func() conds.Condition { return cthen }))
	assert.Equal(t, celse, conds.IFF(false, func() conds.Condition { return cthen }))
}

func TestNonNil(t *testing.T) {
	var n = 100

	tests := []struct {
		c      conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			c: conds.NonNil(&n, func(v int) conds.Condition {
				return conds.C("foo = @n", conds.NV("n", v))
			}),
			stmt:   `foo = @n`,
			params: map[string]any{"n": 100},
		},
		{
			c: conds.NonNil(nilint, func(v int) conds.Condition {
				return conds.C("bar = @nn", conds.NV("nn", v))
			}),
			empty: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.c), func(t *testing.T) {
			stmt, params := tt.c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, tt.c.Empty())
		})
	}
}

func TestC(t *testing.T) {
	tests := []struct {
		c      conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			c:      conds.C(`foo = "bar"`),
			stmt:   `foo = "bar"`,
			params: map[string]any{},
		},
		{
			c:      conds.C(`foo = @bar`, conds.NV("bar", "zoo")),
			stmt:   `foo = @bar`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			c:      conds.C(`foo IN (@bar, @zoo`, conds.NV("bar", 100), conds.NV("zoo", true)),
			stmt:   `foo IN (@bar, @zoo`,
			params: map[string]any{"bar": 100, "zoo": true},
		},
		{
			c:      conds.C(`foo = @bar OR foo = @zoo`, conds.NVMap(map[string]any{"bar": 100, "zoo": true})...),
			stmt:   `foo = @bar OR foo = @zoo`,
			params: map[string]any{"bar": 100, "zoo": true},
		},
		{
			c:      conds.C(`foo IN (@bar, @zoo`, conds.NV("bar", 100), conds.NV("zoo", true)),
			stmt:   `foo IN (@bar, @zoo`,
			params: map[string]any{"bar": 100, "zoo": true},
		},
		{
			c:      conds.C(`foo = @bar`, conds.XNV("bar", ptr("zoo"))),
			stmt:   `foo = @bar`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			c:      conds.C(`foo = @bar`, conds.XNV("bar", nilstr)),
			stmt:   ``,
			params: nil,
			empty:  true,
		},
		{
			c:      conds.C(`foo = @bar AND zoo = @baz`, conds.NV("bar", 100), conds.XNV("baz", nilint)),
			stmt:   ``,
			params: nil,
			empty:  true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.c), func(t *testing.T) {
			stmt, params := tt.c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, tt.c.Empty())
		})
	}
}

func TestAND(t *testing.T) {
	tests := []struct {
		cs     []conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			cs:     []conds.Condition{},
			stmt:   ``,
			params: nil,
			empty:  true,
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND zoo = @baz AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.XNV("fuga", nilint)),
			},
			stmt:   `foo = "bar" AND zoo = @baz`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.AND(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND (zoo = @baz AND hello = @world) AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.OR(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND (zoo = @baz OR hello = @world) AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.cs), func(t *testing.T) {
			c := conds.AND(tt.cs...)
			stmt, params := c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, c.Empty())
		})
	}
}

func TestOR(t *testing.T) {
	tests := []struct {
		cs     []conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			cs:     []conds.Condition{},
			stmt:   ``,
			params: nil,
			empty:  true,
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR zoo = @baz OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.XNV("fuga", nilint)),
			},
			stmt:   `foo = "bar" OR zoo = @baz`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.AND(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR (zoo = @baz AND hello = @world) OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.OR(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR (zoo = @baz OR hello = @world) OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.cs), func(t *testing.T) {
			c := conds.OR(tt.cs...)
			stmt, params := c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, c.Empty())
		})
	}
}

func TestConditionAND(t *testing.T) {
	tests := []struct {
		cs     []conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			cs:     []conds.Condition{},
			stmt:   ``,
			params: nil,
			empty:  true,
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND zoo = @baz AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.XNV("fuga", nilint)),
			},
			stmt:   `foo = "bar" AND zoo = @baz`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.AND(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND (zoo = @baz AND hello = @world) AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.OR(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" AND (zoo = @baz OR hello = @world) AND hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.cs), func(t *testing.T) {
			c := conds.Empty()
			for _, o := range tt.cs {
				c = c.AND(o)
			}
			stmt, params := c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, c.Empty())
		})
	}
}

func TestConditionOR(t *testing.T) {
	tests := []struct {
		cs     []conds.Condition
		stmt   string
		params map[string]any
		empty  bool
	}{
		{
			cs:     []conds.Condition{},
			stmt:   ``,
			params: nil,
			empty:  true,
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR zoo = @baz OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
				conds.C(`hoge = @fuga`, conds.XNV("fuga", nilint)),
			},
			stmt:   `foo = "bar" OR zoo = @baz`,
			params: map[string]any{"bar": "zoo"},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.AND(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR (zoo = @baz AND hello = @world) OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
		{
			cs: []conds.Condition{
				conds.C(`foo = "bar"`),
				conds.OR(
					conds.C(`zoo = @baz`, conds.NV("bar", "zoo")),
					conds.C(`hello = @world`, conds.NV("hello", "world")),
				).Enclose(),
				conds.C(`hoge = @fuga`, conds.NV("fuga", 100)),
			},
			stmt:   `foo = "bar" OR (zoo = @baz OR hello = @world) OR hoge = @fuga`,
			params: map[string]any{"bar": "zoo", "hello": "world", "fuga": 100},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d: %#v", i, tt.cs), func(t *testing.T) {
			c := conds.Empty()
			for _, o := range tt.cs {
				c = c.OR(o)
			}
			stmt, params := c.StmtParams()
			assert.Equal(t, tt.stmt, stmt)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.empty, c.Empty())
		})
	}
}

func TestConditionAND_C(t *testing.T) {
	c := conds.C(`foo = "bar"`).
		AND_C(`zoo = @baz`, conds.NV("bar", "zoo")).
		AND_C(`hoge = @fuga`, conds.NV("fuga", 100)).
		AND_C(`hello = @world`, conds.XNV("world", nilstr))

	stmt, params := c.StmtParams()
	assert.Equal(t, `foo = "bar" AND zoo = @baz AND hoge = @fuga`, stmt)
	assert.Equal(t, map[string]any{"bar": "zoo", "fuga": 100}, params)
}

func TestConditionOR_C(t *testing.T) {
	c := conds.C(`foo = "bar"`).
		OR_C(`zoo = @baz`, conds.NV("bar", "zoo")).
		OR_C(`hoge = @fuga`, conds.NV("fuga", 100)).
		OR_C(`hello = @world`, conds.XNV("world", nilint))

	stmt, params := c.StmtParams()
	assert.Equal(t, `foo = "bar" OR zoo = @baz OR hoge = @fuga`, stmt)
	assert.Equal(t, map[string]any{"bar": "zoo", "fuga": 100}, params)
}
