package main

import (
  "fmt"
  "flag"
  "./raven"
  "./raven/ravenLog"
)

// this will never have true in it.. it's jsut a simple way to hold the main thread open
var done = make( chan bool)

func main() {
  ravenLog.SendError(10, "main", "Starting...")

  configFile := flag.String("config", "../etc/raven.ini", "Configuration File")
  webPort := flag.String("port", ":8000", "Webserver Port")
  workers := flag.Int("workers", 3, "Worker Process")
  flag.Parse()

  ravenLog.SendError(10, "main", fmt.Sprintf( "Config File: %s", *configFile))
  ravenLog.SendError(10, "main", fmt.Sprintf( "Listen Port: %s", *webPort))
  ravenLog.SendError(10, "main", fmt.Sprintf( "Workers: %d", *workers))
  raven.ReadConfig( *configFile)
  raven.BuildSchedule()
  raven.StartSchedule( *workers)
  raven.StartWebserver(*webPort)
  <-done
}
