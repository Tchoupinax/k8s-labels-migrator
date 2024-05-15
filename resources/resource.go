package resource

type Resource struct {
	ApiVersion string
	Category   string
	Kind       string
	Selectors  map[string]string
	Name       string
}

type Event struct {
	Id   int
	Name string
}
