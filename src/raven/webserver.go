package raven

import (
  "log"
//  "fmt"
//  "time"
  "net/http"
  "encoding/json"
//  "bytes"
//  "./ravenTypes"
)

type statusOutput struct {
  Name    string    `json:"dspyName"`
  Group   string    `json:"dspyGroup"`
  Check   string    `json:"check"`
  Options string    `json:"options"`
  Lastrun int64     `json:"lastRunUx"`
  Nextrun int64     `json:"nextRunUx"`
  LastChg int64     `json:"lastChangeUx"`
  ChgThr  string    `json:"changeThreshold"`
  Exit    int       `json:"exitCode"`
  Output  string    `json:"output"`
  Perf    string    `json:"perfData"`
  Text    string    `json:"longText"`
}

func jsonStatus(w http.ResponseWriter, r *http.Request) {
  var jsonOut []statusOutput

  for _,stat := range status {
    var t  statusOutput
    t.Name    = stat.Host.DisplayName
    t.Group   = stat.Host.Group
    t.Check   = stat.Check.DisplayName
    t.Lastrun = stat.Last.Unix()
    t.Nextrun = stat.Next.Unix()
    t.Exit    = stat.Return.Exit
    t.Output  = stat.Return.Text
    t.Perf    = stat.Return.Perf
    t.Text    = stat.Return.Long
    jsonOut = append( jsonOut, t)
  }
  enc := json.NewEncoder(w)
  enc.Encode(jsonOut)
}

func StartWebserver(port string) {
  http.HandleFunc("/api/status", jsonStatus)
  http.HandleFunc("/status", webStatus)
  log.Printf( "Webserver Starting '%s'", port)
  log.Fatal(http.ListenAndServe(port, nil))
}
