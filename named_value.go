package conds

type NamedValue struct {
	name  string
	value any
}

func (nv NamedValue) null() bool {
	return nv.name == ""
}

func NV(name string, value any) NamedValue {
	return NamedValue{name: name, value: value}
}

func XNV[T any](name string, value *T) NamedValue {
	if value == nil {
		return NamedValue{}
	}

	return NV(name, *value)
}

func NVMap(m map[string]any) []NamedValue {
	nvs := []NamedValue{}

	for k, v := range m {
		nvs = append(nvs, NV(k, v))
	}

	return nvs
}
