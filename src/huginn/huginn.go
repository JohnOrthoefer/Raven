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

package main

import (
	"./license"
	"./raven"
	"./raven/ravenLog"
  "./raven/ravenChecks"
	"flag"
	"fmt"
	"log"
)

// this will never have true in it.. it's jsut a simple way to hold the main thread open
var done = make(chan bool)

func main() {
	license.LogLicense(VERSION, COMMIT)
	log.Printf("Commit: %s", COMMIT)
	ravenLog.SendError(10, "main", "Starting...")

	configFile := flag.String("config", "../etc/raven.ini", "Configuration File")
	webPort := flag.String("port", ":8000", "Webserver Port")
	workers := flag.Int("workers", 3, "Worker Process")
	version := flag.Bool("version", false, "Display Full Version")
  plugins := flag.String("plugdir", "./plugins", "Plugins Directory")
	flag.Parse()

	if *version {
		log.Fatal(FULL)
	}

	ravenLog.SendError(10, "main", fmt.Sprintf("Config File: %s", *configFile))
	ravenLog.SendError(10, "main", fmt.Sprintf("PluginDir: %s", *plugins))
	ravenLog.SendError(10, "main", fmt.Sprintf("Listen Port: %s", *webPort))
	ravenLog.SendError(10, "main", fmt.Sprintf("Workers: %d", *workers))

  ravenChecks.LoadPlugins(*plugins)
	raven.ReadConfig(*configFile)
	raven.BuildSchedule()
	raven.StartSchedule(*workers)
	raven.StartWebserver(*webPort)
	<-done
}
