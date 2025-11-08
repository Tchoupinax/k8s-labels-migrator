package webapp

import (
	_ "embed"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	resource "github.com/Tchoupinax/k8s-labels-migrator/resources"
	utils "github.com/Tchoupinax/k8s-labels-migrator/utils"
	"github.com/thedevsaddam/renderer"
)

//nolint:gochecknoglobals
var ViewRenderer *renderer.Render

//go:embed views/summary.html
var summaryHTMLPage string

func Setup() {
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("index").Parse(summaryHTMLPage)

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
		templateError := tmpl.ExecuteTemplate(w, "index", data)
		utils.Check(templateError)
	})

	utils.LogInfo("View resumed here: http://localhost:8080")
	err := utils.OpenURL("http://localhost:8080")
	utils.Check(err)

	go func() {
		srv := &http.Server{
			Addr:         ":8080",
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
}
