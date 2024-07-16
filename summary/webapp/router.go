package webapp

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
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
	deploymentName string,
	resources []resource.Resource,
	podLabels map[string]string,
) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpls := []string{"summary/webapp/views/summary.html"}
		ViewRenderer.FuncMap(template.FuncMap{
			"contains": strings.Contains,
		})

		data := struct {
			DeploymentName      string
			IstioResources      []resource.Resource
			IstioResourcesFound bool
			KedaResources       []resource.Resource
			KedaResourcesFound  bool
			NativeResources     []resource.Resource
			PodLabels           map[string]string
		}{
			DeploymentName:      deploymentName,
			IstioResources:      filter(resources, func(r resource.Resource) bool { return r.Category == "Istio" }),
			IstioResourcesFound: len(filter(resources, func(r resource.Resource) bool { return r.Category == "Istio" })) > 0,
			KedaResources:       filter(resources, func(r resource.Resource) bool { return r.Category == "Keda" }),
			KedaResourcesFound:  len(filter(resources, func(r resource.Resource) bool { return r.Category == "Keda" })) > 0,
			NativeResources:     filter(resources, func(r resource.Resource) bool { return r.Category == "Native" }),
			PodLabels:           podLabels,
		}
		err := ViewRenderer.Template(w, http.StatusOK, tpls, data)
		utils.Check(err)
	})

	utils.LogInfo("View resumed here: http://localhost:8080")
	err := utils.OpenURL("http://localhost:8080")
	utils.Check(err)

	httpServerError := make(chan error, 1)
	go func() {
		httpServerError <- http.ListenAndServe(":8080", nil)
	}()

	if err := <-httpServerError; err != nil {
		fmt.Println(err)
	}
}
