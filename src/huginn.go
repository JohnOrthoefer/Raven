package main

import (
  "log"
  "./raven"
)

var ConfigFile = "../etc/raven.ini"
var done = make( chan bool)

func main() {
  log.Printf( "Starting...")
  log.Printf( "Config File: %s", ConfigFile)
  raven.ReadConfig( ConfigFile)
  raven.BuildSchedule()
  raven.DumpStorage()
  raven.DumpSchedule()
  raven.StartSchedule( 3)
  raven.StartWebserver(":8000")
  <-done
  
}
