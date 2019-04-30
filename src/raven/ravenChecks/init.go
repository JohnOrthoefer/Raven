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

func registerHandler( n string,
                      iptr ravenTypes.CheckInitType,
                      fptr ravenTypes.CheckFuncType) {
  if CheckFunc == nil {
    CheckFunc = make( CheckFMap)
    CheckInit = make( CheckIMap)
  }
  CheckInit[n] = iptr
  CheckFunc[n] = fptr
}
