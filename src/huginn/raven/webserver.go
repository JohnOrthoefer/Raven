/*
   Raven Network Discovery and Monitoring
   Copyright (C) 2019 John{at}Orthoefer{dot}org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package raven

import (
	"./ravenLog"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var templates map[string]*template.Template

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}
type DataType struct {
	Now       string
  HostCnt   int
  CheckCnt  int
  StatusCnt int
  States    [4]int
	Data    interface{}
}
type statusOutputList []statusOutput

var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`
var templateConfig TemplateConfig

type statusOutput struct {
	Name      string `json:"dspyName"`
	Group     string `json:"dspyGroup"`
	Check     string `json:"check"`
	Options   string `json:"options"`
	LastrunUx int64  `json:"lastRunUx"`
	NextrunUx int64  `json:"nextRunUx"`
	LastChgUx int64  `json:"lastChangeUx"`
	Lastrun   string `json:"lastRun"`
	Nextrun   string `json:"nextRun"`
	LastChg   string `json:"lastChange"`
	ChgThr    string `json:"changeThreshold"`
	Exit      int    `json:"exitCode"`
	Output    string `json:"output"`
	Perf      string `json:"perfData"`
	Text      string `json:"longText"`
}

func formatTime(t time.Time) string {
	if t.Unix() < 864000 {
		return "Never"
	}
	return t.Format(time.Stamp)
}

func getStatus() statusOutputList {
	var rtn statusOutputList
	for _, stat := range status {
		var t statusOutput
		t.Name = stat.Host.DisplayName
		t.Group = stat.Host.Group
		t.Check = stat.Check.DisplayName
		t.LastrunUx = stat.Last.Unix()
		t.Lastrun = formatTime(stat.Last)
		t.NextrunUx = stat.Next.Unix()
		t.Nextrun = formatTime(stat.Next)
		t.LastChgUx = stat.Change.Unix()
		t.LastChg = formatTime(stat.Change)
		t.ChgThr = fmt.Sprintf("%d/%d", stat.Count, stat.Check.Threshold)
		t.Exit = stat.CurExit
		r := stat.Return
		t.Output = r.Text
		t.Perf = r.Perf
		t.Text = r.Long
		rtn = append(rtn, t)
	}
	return rtn
}

func getStatusByGroup() map[string][]statusOutput {
	rtn := make(map[string][]statusOutput)
	for _, ent := range getStatus() {
		rtn[ent.Group] = append(rtn[ent.Group], ent)
	}
	return rtn
}

func jsonStatus(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	enc.Encode(getStatus())
}

func loadConfiguration() {
  // Todo: add a configuration file that allows you to move these
	templateConfig.TemplateLayoutPath = "templates/layouts/"
	templateConfig.TemplateIncludePath = "templates/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

// These are the overall layout files (headers, menus etc)
	layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	ravenLog.SendError(10, "loadTemplates", fmt.Sprintf("layoutFiles: %v", layoutFiles))

// These are the specific parts included to make up each page
	includeFiles, err := filepath.Glob(templateConfig.TemplateIncludePath + "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	mainTemplate := template.New("main")

	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	ravenLog.SendError(10, "loadTemplates", fmt.Sprintf("templates loading successful"))
}

func countStatus() [4]int {
  var cnts [4]int

  for _,v := range status {
    cnts[v.CurExit] += 1
  }
  return cnts
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		ravenLog.SendMessage(10, "renderTemplate", fmt.Sprintf("The template %s does not exist.", name))
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	d := DataType{
		Now:        time.Now().Format(time.UnixDate),
    HostCnt:    len(hosts),
    CheckCnt:   len(checks),
    StatusCnt:  len(status),
    States:     countStatus(),
		Data: data,
	}
	err := tmpl.Execute(w, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func tabStatus(w http.ResponseWriter, r *http.Request) {
	data := getStatusByGroup()
	renderTemplate(w, "tabstatus.tmpl", data)
}

func errStatus(w http.ResponseWriter, r *http.Request) {
	data := statusOutputList{}

	for _, v := range getStatus() {
		if v.Exit > 0 {
			data = append(data, v)
		}
	}

	sort.Slice(data, func(i, j int) bool {
		if data[i].Exit == data[j].Exit {
			return strings.ToLower(data[i].Name) < strings.ToLower(data[j].Name)
		}
		return data[i].Exit < data[j].Exit
	})

	renderTemplate(w, "status.tmpl", data)
}

func webStatus(w http.ResponseWriter, r *http.Request) {
	data := getStatus()
	sort.Slice(data, func(i, j int) bool {
		if data[i].Group == data[j].Group {
			return strings.ToLower(data[i].Name) < strings.ToLower(data[j].Name)
		}
		return data[i].Group < data[j].Group
	})
	renderTemplate(w, "status.tmpl", data)
}

func logMessages(w http.ResponseWriter, r *http.Request) {
	data := ravenLog.GetLog()
	for i := len(data)/2 - 1; i >= 0; i-- {
		opp := len(data) - 1 - i
		data[i], data[opp] = data[opp], data[i]
	}
	renderTemplate(w, "logs.tmpl", data)
}

func lastMessage(w http.ResponseWriter, r *http.Request) {
	data := ravenLog.GetLastMessage()
	renderTemplate(w, "thread.tmpl", data)
}

func errMessage(w http.ResponseWriter, r *http.Request) {
	data := ravenLog.GetErrors()
	renderTemplate(w, "errors.tmpl", data)
}

func aboutMessage(w http.ResponseWriter, r *http.Request) {
  data := ""
	renderTemplate(w, "about.tmpl", data)
}

func StartWebserver(port string) {

	ravenLog.SendError(10, "StartWebserver", "Loading Templates")
	loadConfiguration()
	loadTemplates()

	ravenLog.SendError(10, "StartWebserver", "Loading Handler functions")
	http.HandleFunc("/", aboutMessage)
	http.HandleFunc("/errors", errStatus)
	http.HandleFunc("/status", webStatus)
	http.HandleFunc("/tabstatus", tabStatus)
	http.HandleFunc("/log", logMessages)
	http.HandleFunc("/thread", lastMessage)
	http.HandleFunc("/startup", errMessage)
	http.HandleFunc("/api/status", jsonStatus)
	ravenLog.SendError(12, "StartWebServer", fmt.Sprintf("Webserver Starting '%s'", port))
	go func() {
    log.Fatal(http.ListenAndServe(port, nil))
  }()
}
