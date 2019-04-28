package ravenTypes

import (
  "time"
)

type Kwargs map[string]string
type CheckInitType func( Kwargs) interface{}
type CheckFuncType func( *HostEntry, interface{}) *ExitReturn

type HostEntry struct {
  IPv4        string
  Hostname    string
  DisplayName string
  Group       string
}

type ExitReturn struct {
  Exit int
  Text string
  Perf string
  Long string
}

type CheckEntry struct {
  DisplayName string
  CheckF      CheckFuncType
  Options     interface{}
  CheckN      string
  Interval    [4]time.Duration
  Hosts       []string
}

// the basic entry for scheduling a check against a host
type StatusEntry struct {
  Check     *CheckEntry
  Host      *HostEntry
  Queued    bool
  Next      time.Time
  Last      time.Time
  Return    *ExitReturn
}
