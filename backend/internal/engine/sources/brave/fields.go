package brave

type fields map[string]any

func (f fields) obj(key string) fields {
	m, _ := f[key].(map[string]any)
	return m
}

func (f fields) list(key string) []any {
	l, _ := f[key].([]any)
	return l
}

func (f fields) str(key string) string {
	s, _ := f[key].(string)
	return s
}

func (f fields) num(key string) float64 {
	v, _ := richFloat(f[key])
	return v
}

func (f fields) i64(key string) int64 {
	v, _ := intVal(f[key])
	return v
}
