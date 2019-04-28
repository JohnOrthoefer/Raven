package ravenChecks

import (
  "../ravenTypes"
)

type CheckIMap map[string]ravenTypes.CheckInitType
type CheckFMap map[string]ravenTypes.CheckFuncType

var CheckInit CheckIMap
var CheckFunc CheckFMap

func init() {
}
