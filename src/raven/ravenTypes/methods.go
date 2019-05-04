package ravenTypes

import (
//  "log"
  "regexp"
  "strings"
  "strconv"
)

var reSpaces *regexp.Regexp

func init() {
  reSpaces = regexp.MustCompile( `\s+`)
}

func (kw Kwargs) GetKwargStr( s string, d string) string {
  s = strings.ToLower(s)
  if rtn,ok := kw[s]; ok {
    return rtn
  }
  return d
}

func (kw Kwargs) GetKwargStrTrim( s string, d string) string {
  return strings.TrimSpace(kw.GetKwargStr(s,d))
}


func (kw Kwargs) GetKwargStrA( s string, d []string) []string {
  s = strings.ToLower(s)
  if t, ok := kw[s]; ok {
    return reSpaces.Split( strings.TrimSpace(t), -1)
  }
  return d
}

func (kw Kwargs) GetKwargBool( s string, d bool) bool {
  s = strings.ToLower(s)
  if t, ok := kw[s]; ok {
    t = strings.ToLower( t)
    for _,i:=range []string{"true", "t", "yes", "y"} {
      if i == t {
        return true
      }
    }
    return false
  }
  return d
}

func (kw Kwargs) GetKwargFloat( s string, d float64) float64 {
  s = strings.ToLower(s)
  if rtn, ok := kw[s]; ok {
    if fl, err := strconv.ParseFloat( rtn, 64); err == nil {
      return fl
    }
  }
  return d
}

func (kw Kwargs) GetKwargInt( s string, d int64) int64 {
  s = strings.ToLower(s)
  if rtn, ok := kw[s]; ok {
    if num, err := strconv.ParseInt( rtn, 10, 0); err == nil {
      return num
    }
  }
  return d
}

