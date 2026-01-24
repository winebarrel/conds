package conds_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/conds"
)

func TestIF(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		{expr: conds.IF(true, "vthen", "velse"), expected: "vthen"},
		{expr: conds.IF(false, "vthen", "velse"), expected: "velse"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.expr)
		})
	}
}
