package handler

import (
	"github.com/edward/scp-294/internal"
	"github.com/edward/scp-294/internal/constants"
	"github.com/edward/scp-294/internal/dbaccess"
	"github.com/edward/scp-294/pkg/processor"
	"html/template"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/", loadMainPage)
	registerConverterRoutes()
	registerGroupRoutes()
	registerRecordRoutes()
}

func loadMainPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("layout")
	t, err := t.ParseFiles("./web/template/layout.html", "./web/template/header.html")
	if err != nil {
		internal.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	groups, err := dbaccess.ListGroups()
	if err != nil {
		internal.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	records, err := dbaccess.ListRecords()
	if err != nil {
		internal.Logger.Log("Main", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := struct {
		ProjName            string
		InoutTypes          []constants.Type
		InoutFormatMappings []constants.FormatMapping
		Groups              []dbaccess.Group
		Records             []dbaccess.Record
		Processors          []string
	}{
		"SCP-294",
		constants.CreateTypes(),
		constants.CreateFormatMappings(),
		groups,
		records,
		append([]string{"None"}, processor.DefaultProcRegistry.GetProcessorNames()...),
	}
	t.ExecuteTemplate(w, "layout", data)
}

func reloadHeader(w http.ResponseWriter) {
	t := template.New("reloadHeader")
	t, err := t.ParseFiles("./web/template/header.html")
	if err != nil {
		respondError(w, "Failed to parse files, error: "+err.Error())
		return
	}
	groups, err := dbaccess.ListGroups()
	if err != nil {
		respondError(w, "Failed to list groups, error: "+err.Error())
		return
	}
	records, err := dbaccess.ListRecords()
	if err != nil {
		respondError(w, "Failed to list records, error: "+err.Error())
		return
	}
	data := struct {
		Groups  []dbaccess.Group
		Records []dbaccess.Record
	}{groups, records}
	t.ExecuteTemplate(w, "header", data)
}
