package resource

type Resource struct {
	ApiVersion string
	Category   string
	Kind       string
	Labels     map[string]string
	Name       string
	Selectors  map[string]string
}

type Event struct {
	Id   int
	Name string
}
