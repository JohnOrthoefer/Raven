package ravenChecks

import (
  "../ravenTypes"
)

var CheckFunc map[string]func( ravenTypes.HostEntry, map[string]string) (int, [3]string)

func init() {
}
