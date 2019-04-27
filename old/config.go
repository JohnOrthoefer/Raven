package main

import (
  "log"
  toml "github.com/BurntSushi/toml"
  "../lib"
)

type tomlMap map[string]map[string]string

var config tomlMap

func InitConfig() {
	if _, err := toml.DecodeFile( "../etc/config.toml", &config); err != nil {
    log.Fatal(err)
  }

  for name,_ := range config {
    if h,ok := config[name]["host"]; ok {
      if e,ok := config[name]["enabled"]; ok && e != "false" {
	log.Printf( "%v", config[name])
        g := config[name]["group"]
        storage.NewHost( name, h, g)
      }
    }
  }

}

func main() {
  InitConfig()
  storage.PrintHosts()
}
