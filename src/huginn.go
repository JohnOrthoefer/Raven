package main

import (
  "log"
  "flag"
  "./raven"
)

// this will never have true in it.. it's jsut a simple way to hold the main thread open
var done = make( chan bool)

func main() {
  log.Printf( "Starting...")

  configFile := flag.String("config", "../etc/raven.ini", "Configuration File")
  webPort := flag.String("port", ":8000", "Webserver Port")
  workers := flag.Int("workers", 3, "Worker Process")
  flag.Parse()

  log.Printf( "Config File: %s", *configFile)
  log.Printf( "Listen Port: %s", *webPort)
  log.Printf( "Workers: %d", *workers)
  raven.ReadConfig( *configFile)
  raven.BuildSchedule()
  raven.DumpStorage()
  raven.DumpSchedule()
  raven.StartSchedule( *workers)
  raven.StartWebserver(*webPort)
  <-done
}
