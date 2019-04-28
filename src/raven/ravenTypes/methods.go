package ravenTypes

import (
//  "log"
  "strconv"
)

func (kw Kwargs) GetKwargStr( s string, d string) string {
  if rtn,ok := kw[s]; ok {
    return rtn
  }
  return d
}

func (kw Kwargs) GetKwargFloat( s string, d float64) float64 {
  if rtn, ok := kw[s]; ok {
    if fl, err := strconv.ParseFloat( rtn, 64); err == nil {
      return fl
    }
  }
  return d
}

func (kw Kwargs) GetKwargInt( s string, d int64) int64 {
  if rtn, ok := kw[s]; ok {
    if num, err := strconv.ParseInt( rtn, 10, 0); err == nil {
      return num
    }
  }
  return d
}

