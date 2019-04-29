package raven

// project I'm emulating
// https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6

import (
  "log"
//  "fmt"
  "time"
  "net/http"
  "encoding/json"
  "html/template"
  "path/filepath"
//  "bytes"
//  "./ravenTypes"
)

var templates map[string]*template.Template

type TemplateConfig struct {
  TemplateLayoutPath  string
  TemplateIncludePath string
}
type DataType struct {
  Now string
  Data interface{}
}
var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`
var templateConfig TemplateConfig

type statusOutput struct {
  Name      string    `json:"dspyName"`
  Group     string    `json:"dspyGroup"`
  Check     string    `json:"check"`
  Options   string    `json:"options"`
  LastrunUx int64     `json:"lastRunUx"`
  NextrunUx int64     `json:"nextRunUx"`
  LastChgUx int64     `json:"lastChangeUx"`
  Lastrun   string    `json:"lastRun"`
  Nextrun   string    `json:"nextRun"`
  LastChg   string    `json:"lastChange"`
  ChgThr    string    `json:"changeThreshold"`
  Exit      int       `json:"exitCode"`
  Output    string    `json:"output"`
  Perf      string    `json:"perfData"`
  Text      string    `json:"longText"`
}

func getStatus() []statusOutput {
  var rtn []statusOutput
  for _,stat := range status {
    var t  statusOutput
    t.Name      = stat.Host.DisplayName
    t.Group     = stat.Host.Group
    t.Check     = stat.Check.DisplayName
    t.LastrunUx = stat.Last.Unix()
    t.Lastrun   = stat.Last.Format(time.UnixDate)
    t.NextrunUx = stat.Next.Unix()
    t.Nextrun   = stat.Next.Format(time.UnixDate)
    t.LastChgUx = stat.Change.Unix()
    t.LastChg   = stat.Change.Format(time.UnixDate)
    r := stat.Return
    if r == nil {
      r = stat.OldRtn
    } 
    t.Exit      = r.Exit
    t.Output    = r.Text
    t.Perf      = r.Perf
    t.Text      = r.Long
    rtn         = append( rtn, t)
  }
  return rtn
}

func jsonStatus(w http.ResponseWriter, r *http.Request) {
  enc := json.NewEncoder(w)
  enc.Encode(getStatus())
}


func loadConfiguration() {
  templateConfig.TemplateLayoutPath = "templates/layouts/"
  templateConfig.TemplateIncludePath = "templates/"
}

func loadTemplates() {
  if templates == nil {
    templates = make(map[string]*template.Template)
  }

  layoutFiles, err := filepath.Glob(templateConfig.TemplateLayoutPath + "*.tmpl")
  if err != nil {
    log.Fatal(err)
  }
  log.Printf( "layoutFiles: %v", layoutFiles)

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

  log.Println("templates loading successful")
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
  tmpl, ok := templates[name]
  if !ok {
    log.Printf("The template %s does not exist.", name)
    return
  }

  w.Header().Set("Content-Type", "text/html; charset=utf-8")

  d := DataType {
    Now: time.Now().Format(time.UnixDate),
    Data: data,
  }
  err := tmpl.Execute(w, d)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func webStatus(w http.ResponseWriter, r *http.Request) {
  data := getStatus()
  renderTemplate(w, "status.tmpl", data)
}

func StartWebserver(port string) {

  log.Printf( "Loading Templates")
  loadConfiguration()
  loadTemplates()

  log.Printf( "Loading Handler /status")
  http.HandleFunc("/status", webStatus)
  log.Printf( "Loading Handler /api/status")
  http.HandleFunc("/api/status", jsonStatus)
  log.Printf( "Webserver Starting '%s'", port)
  log.Fatal(http.ListenAndServe(port, nil))
}
