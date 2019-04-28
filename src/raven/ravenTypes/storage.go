package ravenTypes

import (
  "time"
)

type HostEntry struct {
  IPv4        string
  Hostname    string
  DisplayName string
  Group       string
}

type CheckEntry struct {
  DisplayName string
  CheckF      func( HostEntry, map[string]string) (int, [3]string)
  CheckN      string
  Interval    [4]time.Duration
  Hosts       []string
  Options     map[string]string
}
