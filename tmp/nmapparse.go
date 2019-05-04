package main

import (
  "log"
  "fmt"
  "strings"
  "strconv"
  "time"
  "encoding/xml"
  "io/ioutil"
  "bytes"
  ini "github.com/ochinchina/go-ini"
)

type NmapInfo struct {
//  XMLName xml.Name  `xml:"nmaprun"`
  Scanner string `xml:"scanner,attr"`
  Args    string `xml:"args,attr"`
  Start   int64 `xml:"start,attr"`
  StartS  string  `xml:"startstr,attr"`
  version string  `xml:"version,attr"`
  xmloutputvers string  `xml:"xmloutputversion,attr"`
  Info  InfoStruct `xml:"scaninfo"`
  Host  []HostStruct `xml:"host"`
}

type InfoStruct struct {
  Type  string  `xml:"type,attr"`
  Proto string  `xml:"protocol,attr"`
  count string  `xml:"numservices,attr"`
  svrs  string  `xml:"services,attr"`
}

type StatusStruct struct {
  State     string    `xml:"state,attr"`
  Reason    string    `xml:"reason,attr"`
  Reasonttl int       `xml:"reason_ttl,attr"`
}
type AddrStruct struct {
  Addr      string    `xml:"addr,attr"`
  Type      string    `xml:"addrtype,attr"`
}
type HostnameStruct struct {
  Name      string    `xml:"name,attr"`
  Type      string    `xml:"type,attr"`
}
type StateStruct struct {
  State     string    `xml:"state,attr"`
  Reason    string    `xml:"reason,attr"`
  Reasonttl int       `xml:"reason_ttl,attr"`
}
type ServiceStruct struct {
  Name      string    `xml:"name,attr"`
  Method    string    `xml:"method,attr"`
  Conf      int       `xml:"conf,attr"`
}
type PortStruct struct {
  Protocol  string    `xml:"protocol,attr"`
  PortID    int    `xml:"portid,attr"`
  State     StateStruct `xml:"state"`
  Service   ServiceStruct `xml:"service"`
}

type HostStruct struct {
  StartTime int64              `xml:"starttime,attr"`
  EndTime   int64              `xml:"endtime,attr"`
  Status    StatusStruct        `xml:"status"`
  Addr      []AddrStruct        `xml:"address"`
  Hostname  []HostnameStruct    `xml:"hostnames>hostname"`
  Ports     []PortStruct        `xml:"ports>port"`
}

type HostJSON struct {
  Name    string    `json:"name"`
  Hostname    string    `json:"hostname"`
  IPv4    string    `json:"ipv4"`
  When    time.Time `json:"lastscan"`
  Ports   []int     `json:"ports"`
}


var scanfile string = "scan.xml"
var nmap NmapInfo
var groupName string = "Internal-LAN"
var hosts []HostJSON

func getHostInfo( h HostStruct) (name, hn, hi string) {
  if len(h.Hostname)> 0 {
     hn = h.Hostname[0].Name
  }

  if len(h.Addr) > 0 {
    hi = h.Addr[0].Addr
  }

  if hn == "" {
    s := strings.Split( hi, ".")
    v,_ := strconv.Atoi( s[3])
    hn = fmt.Sprintf( "NO-Name-%d", v)
  }
  s := strings.Split(hn, ".")
  return s[0], hn, hi
}

func main() {
  ini := ini.NewIni()

  log.Printf( "Starting reading %s", scanfile)
  xmlblob, err := ioutil.ReadFile( scanfile)
  if  err != nil {
    log.Fatal( err)
  }

  if err := xml.Unmarshal( xmlblob, &nmap); err != nil {
    log.Fatal( err)
  }


  for _,v:= range nmap.Host {
    hr := new( HostJSON)
    hr.Name,hr.Hostname,hr.IPv4 = getHostInfo( v)

    for _,p := range v.Ports {
      if p.State.State == "open" {
        hr.Ports = append(hr.Ports, p.PortID)
      }
    }
    hr.When = time.Unix( v.EndTime, 0)
    section := ini.NewSection( hr.Name)
    section.Add( "hostname", hr.Hostname)
    section.Add( "group", groupName)
    section.Add( "enabled", "true")
    section.Add( "ipv4", hr.IPv4)
    log.Printf( "%s(%s) %v", hr.Name, hr.IPv4, hr.Ports)
    hosts = append(hosts, *hr)
  }

  log.Printf( "%v", hosts)

  buf := bytes.NewBufferString("")
  ini.Write(buf)
  log.Printf( "%v", buf)
}
