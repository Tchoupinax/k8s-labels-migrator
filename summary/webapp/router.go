package webapp

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	"github.com/thedevsaddam/renderer"
)

var ViewRenderer *renderer.Render

func init() {
	ViewRenderer = renderer.New()
}

func filter[T any](array []T, test func(T) bool) (ret []T) {
	for _, s := range array {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func StartWebServer(
	resources []resource.Resource,
) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpls := []string{"summary/webapp/views/summary.html"}
		ViewRenderer.FuncMap(template.FuncMap{
			"contains": strings.Contains,
		})

		data := struct {
			NativeResources []resource.Resource
			IstioResources  []resource.Resource
			KedaResources   []resource.Resource
		}{
			NativeResources: filter(resources, func(r resource.Resource) bool { return r.Category == "Native" }),
			IstioResources:  filter(resources, func(r resource.Resource) bool { return r.Category == "Istio" }),
			KedaResources:   filter(resources, func(r resource.Resource) bool { return r.Category == "Keda" }),
		}
		ViewRenderer.Template(w, http.StatusOK, tpls, data)
	})

	fmt.Printf("Starting server at port 8080\n")
	go http.ListenAndServe(":8080", nil)

	fmt.Println("start server")
}
