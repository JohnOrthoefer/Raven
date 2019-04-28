package ravenTypes


import (
  "time"
)

// the basic entry for scheduling a check against a host
type StatusEntry struct {
  Check     *CheckEntry
  Host      *HostEntry
  ExitCode  int
  Queued    bool
  Next      time.Time
  Last      time.Time
}
