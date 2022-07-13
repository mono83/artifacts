package data

// Result represents artifact data result
type Result struct {
	Metric string
	Group  map[string]string
	Value  int64
}
