package raven

// ping check command

import (
  "log"
)

func Ping( he HostEntry, opts map[string]string) (int, []string) {
  var rtnOut  []string // 0 = text; 1 = perf; 2 = extended text

  rtnExit, output := runExternal( "/usr/bin/ping", "-c", "5", he.Hostname)
  log.Printf( "%s(Ping) exit:%d out=%s", he.Hostname, rtnExit, output)

  return rtnExit, rtnOut
}

