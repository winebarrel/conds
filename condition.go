package conds

import (
	"maps"
	"strings"
)

type Condition struct {
	stmt   string
	params map[string]any
}

func (c Condition) Empty() bool {
	return c.stmt == ""
}

func (c Condition) Enclose() Condition {
	return Condition{
		stmt:   "(" + c.stmt + ")",
		params: c.params,
	}
}

func (c Condition) StmtParams() (string, map[string]any) {
	return c.stmt, c.params
}

func (c Condition) AND(other Condition) Condition {
	return AND(c, other)
}

func (c Condition) AND_C(stmt string, nvs ...NamedValue) Condition {
	return AND(c, C(stmt, nvs...))
}

func (c Condition) OR(other Condition) Condition {
	return OR(c, other)
}

func (c Condition) OR_C(stmt string, nvs ...NamedValue) Condition {
	return OR(c, C(stmt, nvs...))
}

/////////////////////////////////////////////////////////////////////

func Empty() Condition {
	return Condition{}
}

func C(stmt string, nvs ...NamedValue) Condition {
	params := map[string]any{}

	for _, nv := range nvs {
		if nv.null() {
			return Condition{}
		}

		params[nv.name] = nv.value
	}

	return Condition{
		stmt:   stmt,
		params: params,
	}
}

func AND(conditions ...Condition) Condition {
	return join("AND", conditions...)
}

func OR(others ...Condition) Condition {
	return join("OR", others...)
}

func join(op string, others ...Condition) Condition {
	stmts := []string{}
	params := map[string]any{}

	for _, o := range others {
		if !o.Empty() {
			stmts = append(stmts, o.stmt)
			maps.Copy(params, o.params)
		}
	}

	if len(stmts) == 0 {
		return Condition{}
	}

	return Condition{
		stmt:   strings.Join(stmts, " "+op+" "),
		params: params,
	}
}

func IF[T any](expr bool, cthen, celse T) T {
	if expr {
		return cthen
	} else {
		return celse
	}
}

func NonNil[T any](param *T, f func(T) Condition) Condition {
	if param == nil {
		return Condition{}
	}

	return f(*param)
}
