package conds_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/conds"
)

func TestV(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{name: "foo", value: "bar", expected: `conds.NamedValue{name:"foo", value:"bar"}`},
		{name: "zoo", value: 100, expected: `conds.NamedValue{name:"zoo", value:100}`},
		{name: "baz", value: nil, expected: `conds.NamedValue{name:"baz", value:interface {}(nil)}`},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, fmt.Sprintf("%#v", conds.V(tt.name, tt.value)))
		})
	}
}

func TestXV(t *testing.T) {
	tests := []struct {
		name     string
		value    *any
		expected string
	}{
		{name: "foo", value: ptr("bar"), expected: `conds.NamedValue{name:"foo", value:"bar"}`},
		{name: "zoo", value: ptr(100), expected: `conds.NamedValue{name:"zoo", value:100}`},
		{name: "baz", value: nil, expected: `conds.NamedValue{name:"", value:interface {}(nil)}`},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, fmt.Sprintf("%#v", conds.XV(tt.name, tt.value)))
		})
	}
}

func TestVMap(t *testing.T) {
	tests := []struct {
		m     map[string]any
		items []string
	}{
		{
			m: map[string]any{"foo": "bar", "zoo": 100, "baz": nil},
			items: []string{
				`conds.NamedValue{name:"foo", value:"bar"}`,
				`conds.NamedValue{name:"zoo", value:100}`,
				`conds.NamedValue{name:"baz", value:interface {}(nil)}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.m), func(t *testing.T) {
			nvs := conds.VMap(tt.m)
			assert.Len(t, nvs, len(tt.m))
			s := fmt.Sprintf("%#v", nvs)
			for _, i := range tt.items {
				assert.Contains(t, s, i)
			}
		})
	}
}
