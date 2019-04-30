package raven

import (
  ini "github.com/ochinchina/go-ini"
  "fmt"
  "./ravenLog"
)

type stringMap map[string]string

func makeMap( s *ini.Section) stringMap {
  r:=make(stringMap)
  for _,key := range s.Keys() {
    n := key.Name()
    v,_ := s.GetValue(n)
    r[n] = v
  }
  return r
}

func ReadConfig(iniFile string) {
  ini := ini.Load( iniFile)
  for _,section:= range ini.Sections() {
    secName := section.Name
    ravenLog.SendError( 10, "Configuration", fmt.Sprintf( "Parsing Section %s", secName))
    keyVal := makeMap(section)
    AddEntry( secName, keyVal)
  }
}
