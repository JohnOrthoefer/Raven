package raven

import (
	"./ravenLog"
	"./ravenTypes"
	"fmt"
	ini "github.com/ochinchina/go-ini"
)

func makeMap(s *ini.Section) ravenTypes.Kwargs {
	r := make(ravenTypes.Kwargs)
	for _, key := range s.Keys() {
		n := key.Name()
		v, _ := s.GetValue(n)
		r[n] = v
	}
	return r
}

func ReadConfig(iniFile string) {
	ini := ini.Load(iniFile)
	for _, section := range ini.Sections() {
		secName := section.Name
		ravenLog.SendError(10, "Configuration", fmt.Sprintf("Parsing Section %s", secName))
		keyVal := makeMap(section)
		AddEntry(secName, keyVal)
	}
}
