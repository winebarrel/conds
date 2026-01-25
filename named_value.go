package conds

type NamedValue struct {
	name  string
	value any
}

func (nv NamedValue) null() bool {
	return nv.name == ""
}

func V(name string, value any) NamedValue {
	return NamedValue{name: name, value: value}
}

func XV[T any](name string, value *T) NamedValue {
	if value == nil {
		return NamedValue{}
	}

	return V(name, *value)
}

func VMap(m map[string]any) []NamedValue {
	nvs := []NamedValue{}

	for k, v := range m {
		nvs = append(nvs, V(k, v))
	}

	return nvs
}
