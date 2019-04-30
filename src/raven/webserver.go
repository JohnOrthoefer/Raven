package raven

// project I'm emulating
// https://hackernoon.com/golang-template-2-template-composition-and-how-to-organize-template-files-4cb40bcdf8f6

import (
  "log"
  "fmt"
  "time"
  "net/http"
  "encoding/json"
  "html/template"
  "path/filepath"
  "./ravenLog"
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

func getStatusByGroup() map[string][]statusOutput {
  rtn := make( map[string][]statusOutput)
  for _,ent := range getStatus() {
    rtn[ent.Group] = append(rtn[ent.Group], ent)
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
  ravenLog.SendError( 10, "loadTemplates", fmt.Sprintf( "layoutFiles: %v", layoutFiles))

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

  ravenLog.SendError( 10, "loadTemplates", fmt.Sprintf("templates loading successful"))
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
  tmpl, ok := templates[name]
  if !ok {
    ravenLog.SendMessage( 10, "renderTemplate", fmt.Sprintf("The template %s does not exist.", name))
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

func tabStatus(w http.ResponseWriter, r *http.Request) {
  data := getStatusByGroup()
  renderTemplate(w, "tabstatus.tmpl", data)
}

func webStatus(w http.ResponseWriter, r *http.Request) {
  data := getStatus()
  for i := len(data)/2-1; i>=0; i-- {
    opp := len(data)-1-i
    data[i], data[opp] = data[opp], data[i]
  }
  renderTemplate(w, "status.tmpl", data)
}

func logMessages(w http.ResponseWriter, r *http.Request) {
  data := ravenLog.GetLog()
  renderTemplate(w, "logs.tmpl", data)
}

func lastMessage(w http.ResponseWriter, r *http.Request) {
  data := ravenLog.GetLastMessage()
  renderTemplate(w, "logs.tmpl", data)
}

func StartWebserver(port string) {

  ravenLog.SendError( 10, "StartWebserver", "Loading Templates")
  loadConfiguration()
  loadTemplates()

  ravenLog.SendError( 10, "StartWebserver", "Loading Handler functions")
  http.HandleFunc("/status", webStatus)
  http.HandleFunc("/tabstatus", tabStatus)
  http.HandleFunc("/log", logMessages)
  http.HandleFunc("/thread", lastMessage)
  http.HandleFunc("/api/status", jsonStatus)
  ravenLog.SendError( 10, "StartWebServer", fmt.Sprintf( "Webserver Starting '%s'", port))
  log.Fatal(http.ListenAndServe(port, nil))
}
