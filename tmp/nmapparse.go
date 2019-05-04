package main

import (
  "log"
  "bytes"
  "fmt"
  "strings"
  "strconv"
  "time"
  "net"
  "syscall"
  "flag"
  "encoding/xml"
  "encoding/json"
  "io/ioutil"
  "os/exec"
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

func runExternal( prog string, args ...string) (int, string) {

  var out bytes.Buffer

  cmd := exec.Command(prog, args...)

  cmd.Stdout = &out
  cmd.Stderr = &out
  if err := cmd.Start(); err != nil {
    log.Fatalf("cmd.Start: %v")
  }

  rtnExit:=0
  if err := cmd.Wait(); err != nil {
    if exiterr, ok := err.(*exec.ExitError); ok {
    // The program has exited with an exit code != 0
    // This works on both Unix and Windows. Although package
    // syscall is generally platform dependent, WaitStatus is
    // defined for both Unix and Windows and in both cases has
    // an ExitStatus() method with the same signature.
      if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
        rtnExit = status.ExitStatus()
        log.Printf("Exit Status: %d", rtnExit)
      }
    } else {
      log.Fatalf("cmd.Wait: %v", err)
    }
  }
  return rtnExit, out.String()
}

var nmap NmapInfo

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

func findLocal() (string) {
  addr,_ := net.InterfaceAddrs()
  for _,j := range addr {
    //log.Printf( " -Addesses %s", j.String())
    ip,_,err := net.ParseCIDR( j.String())
    if err != nil {
      log.Printf( "ParseCIDR failed")
      continue
    }
    if ip.IsGlobalUnicast() {
      log.Printf( "Found %v", j.String())
      return j.String()
    }
  }

  return ""
}


func runNmap( lnet string) []byte {
  nmapExec := "/usr/bin/nmap"
  nmapOpts := []string{ "-oX", "-",
                        "-p", "22,23,25,80,123,161,162,443",
                        lnet}
  log.Printf( "Running %s %s", nmapExec, strings.Join(nmapOpts, " "))
  e, out := runExternal( nmapExec, nmapOpts...)
  log.Printf( "Done Exit: %d", e)
  log.Printf( "Output length %d characters", len(out))
  if e != 0 {
    log.Fatal( "Exit %d when running %s", e, nmapExec)
  }
  return []byte( out)
}

func main() {
  var err error
  xmlblob := []byte{}

  // Do the CLI flags
  scanfile := flag.String( "xml", "", "XML Output from nmap")
  gname := flag.String( "group", "Internal-LAN",
                        "Group for hosts on this network")
  lnet := flag.String( "net", "", "CIDR network to scan with nmap")
  iniFile := flag.Bool( "ini", false, "Output to file (ini format)")
  jsonFile := flag.Bool( "json", false, "Output to file (json format)")
  flag.Parse()

  groupName := *gname

  // if scanfile is empty try running nmap
  if *scanfile == "" {
    localNet := *lnet
    if localNet == "" {
      log.Printf( "Detecting interface")
      localNet = findLocal()
    }
    // give up, no local net found and none provided
    if localNet == "" {
      log.Fatal( "no interfaces found")
    }
    xmlblob = runNmap( localNet)
  } else {
    // read a provided XML file
    log.Printf( "Starting reading %s", scanfile)
    xmlblob, err = ioutil.ReadFile( *scanfile)
    if  err != nil {
      log.Fatal( err)
    }
  }

  // We should have an xmlblob by now
  if err := xml.Unmarshal( xmlblob, &nmap); err != nil {
    log.Fatal( err)
  }

  ini := ini.NewIni()
  hosts := make(map[string]*HostJSON)
  for _,v:= range nmap.Host {
    hr := new( HostJSON)
    hr.Name,hr.Hostname,hr.IPv4 = getHostInfo( v)

    hr.When = time.Unix( v.EndTime, 0)
    section := ini.NewSection( hr.Name)
    section.Add( "hostname", hr.Hostname)
    section.Add( "group", groupName)
    section.Add( "enabled", "true")
    section.Add( "ipv4", hr.IPv4)
    for _,p := range v.Ports {
      if p.State.State == "open" {
        hr.Ports = append(hr.Ports, p.PortID)
      }
    }
    log.Printf( "%s(%s) %v", hr.Name, hr.IPv4, hr.Ports)
    hosts[hr.Name] = hr
  }

  if *iniFile {
    var b strings.Builder
    b.WriteString(groupName)
    b.WriteString(".ini")
    outfile := b.String()

    buf := bytes.NewBufferString("")
    ini.Write(buf)
    log.Printf( "Writing %s", outfile)
    ini.WriteToFile( outfile)
  }

  if *jsonFile {
    var b strings.Builder
    b.WriteString(groupName)
    b.WriteString(".json")
    outfile := b.String()

    l,_ := json.MarshalIndent(&hosts, "", "  ")
    log.Printf("Writing %s", outfile)
    ioutil.WriteFile( outfile, l, 0644)
    if err != nil {
      log.Fatal( err)
    }
  }
  log.Printf( "Done...")
}
