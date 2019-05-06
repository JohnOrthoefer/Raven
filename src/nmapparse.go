/*
   Raven Network Discovery and Monitoring
   Copyright (C) 2019 John{at}Orthoefer{dot}org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/
package main

import (
	"./license"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	ini "github.com/ochinchina/go-ini"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var nmap NmapInfo

type NmapInfo struct {
	//  XMLName xml.Name  `xml:"nmaprun"`
	Scanner       string       `xml:"scanner,attr"`
	Args          string       `xml:"args,attr"`
	Start         int64        `xml:"start,attr"`
	StartS        string       `xml:"startstr,attr"`
	version       string       `xml:"version,attr"`
	xmloutputvers string       `xml:"xmloutputversion,attr"`
	Info          InfoStruct   `xml:"scaninfo"`
	Host          []HostStruct `xml:"host"`
}

type InfoStruct struct {
	Type  string `xml:"type,attr"`
	Proto string `xml:"protocol,attr"`
	count string `xml:"numservices,attr"`
	svrs  string `xml:"services,attr"`
}

type StatusStruct struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	Reasonttl int    `xml:"reason_ttl,attr"`
}
type AddrStruct struct {
	Addr string `xml:"addr,attr"`
	Type string `xml:"addrtype,attr"`
}
type HostnameStruct struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
}
type StateStruct struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	Reasonttl int    `xml:"reason_ttl,attr"`
}
type ServiceStruct struct {
	Name   string `xml:"name,attr"`
	Method string `xml:"method,attr"`
	Conf   int    `xml:"conf,attr"`
}
type PortStruct struct {
	Protocol string        `xml:"protocol,attr"`
	PortID   int           `xml:"portid,attr"`
	State    StateStruct   `xml:"state"`
	Service  ServiceStruct `xml:"service"`
}

type HostStruct struct {
	StartTime int64            `xml:"starttime,attr"`
	EndTime   int64            `xml:"endtime,attr"`
	Status    StatusStruct     `xml:"status"`
	Addr      []AddrStruct     `xml:"address"`
	Hostname  []HostnameStruct `xml:"hostnames>hostname"`
	Ports     []PortStruct     `xml:"ports>port"`
}

type HostJSON struct {
	Name     string    `json:"name"`
	Hostname string    `json:"hostname"`
	Enabled  bool      `json:"enabled"`
	DHCP     bool      `json:"dhcp"`
	IPv4     string    `json:"ipv4"`
	When     time.Time `json:"lastseen"`
	Ports    []int     `json:"ports"`
}

func runExternal(prog string, args ...string) (int, string) {

	var out bytes.Buffer

	cmd := exec.Command(prog, args...)

	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v")
	}

	rtnExit := 0
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

func getHostInfo(h HostStruct, dl, dh int) (name, hn, hi string, dhcp bool) {
	if len(h.Hostname) > 0 {
		hn = h.Hostname[0].Name
	}

	if len(h.Addr) > 0 {
		hi = h.Addr[0].Addr
	}

	octets := strings.Split(hi, ".")
	lsv, _ := strconv.Atoi(octets[3])
	dhcp = (lsv >= dl) && (lsv <= dh)
	if hn == "" {
		hn = fmt.Sprintf("NO-Name-%d", lsv)
	}
	s := strings.Split(hn, ".")
	return s[0], hn, hi, dhcp
}

func findLocal() string {
	addr, _ := net.InterfaceAddrs()
	for _, j := range addr {
		//log.Printf( " -Addesses %s", j.String())
		ip, _, err := net.ParseCIDR(j.String())
		if err != nil {
			log.Printf("ParseCIDR failed")
			continue
		}
		if ip.IsGlobalUnicast() {
			log.Printf("Found %v", j.String())
			return j.String()
		}
	}

	return ""
}

func runNmap(lnet string) []byte {
	nmapExec := "/usr/bin/nmap"
	nmapOpts := []string{"-oX", "-",
		"-p", "22,23,25,80,123,161,162,443",
		lnet}
	log.Printf("Running %s %s", nmapExec, strings.Join(nmapOpts, " "))
	e, out := runExternal(nmapExec, nmapOpts...)
	log.Printf("Done Exit: %d", e)
	log.Printf("Output length %d characters", len(out))
	if e != 0 {
		log.Fatal("Exit %d when running %s", e, nmapExec)
	}
	return []byte(out)
}

func main() {
	var err error
	xmlblob := []byte{}

	// print license
	license.LogLicense()

	// Do the CLI flags
	scanfile := flag.String("xml", "", "XML Output from nmap")
	gname := flag.String("group", "Internal-LAN",
		"Group for hosts on this network")
	lnet := flag.String("net", "", "CIDR network to scan with nmap")
	dhcpRange := flag.String("dhcp", "100-200", "DHCP address range")
	iniFile := flag.Bool("ini", false, "Output to file (ini format)")
	jsonFile := flag.Bool("json", false, "Output to file (json format)")
	baseini := flag.String("skel", "base.ini", "Skeleton file for ini file generaton")
	disabled := flag.Bool("disabled", false, "Mark all hosts as enabled=false")
	flag.Parse()

	groupName := *gname

	// if scanfile is empty try running nmap
	if *scanfile == "" {
		localNet := *lnet
		if localNet == "" {
			log.Printf("Detecting interface")
			localNet = findLocal()
		}
		// give up, no local net found and none provided
		if localNet == "" {
			log.Fatal("no interfaces found")
		}
		xmlblob = runNmap(localNet)
	} else {
		// read a provided XML file
		log.Printf("Starting reading %s", scanfile)
		xmlblob, err = ioutil.ReadFile(*scanfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	// We should have an xmlblob by now
	if err := xml.Unmarshal(xmlblob, &nmap); err != nil {
		log.Fatal(err)
	}

	//ini := ini.NewIni()
	ini := ini.Load(baseini)
	hosts := make(map[string]*HostJSON)
	portsEnabled := make(map[int][]string)

	// This may need to get better
	dhcp := strings.Split(*dhcpRange, "-")
	dhcplow, _ := strconv.Atoi(dhcp[0])
	dhcphi, _ := strconv.Atoi(dhcp[1])
	log.Printf("DHCP range is %d to %d", dhcplow, dhcphi)

	for _, v := range nmap.Host {
		hr := new(HostJSON)
		hr.Name, hr.Hostname, hr.IPv4, hr.DHCP = getHostInfo(v, dhcplow, dhcphi)
		hr.When = time.Unix(v.EndTime, 0)
		hr.Enabled = !*disabled
		section := ini.NewSection(hr.Name)
		section.Add("hostname", hr.Hostname)
		section.Add("group", groupName)
		section.Add("enabled", fmt.Sprintf("%t", hr.Enabled))
		if !hr.DHCP {
			section.Add("ipv4", hr.IPv4)
		}
		portsEnabled[0] = append(portsEnabled[0], hr.Name)
		for _, p := range v.Ports {
			if p.State.State == "open" {
				hr.Ports = append(hr.Ports, p.PortID)
				portsEnabled[p.PortID] = append(portsEnabled[p.PortID], hr.Name)
			}
		}
		log.Printf("%s(%s) %v", hr.Name, hr.IPv4, hr.Ports)
		hosts[hr.Name] = hr
	}
	for k, v := range portsEnabled {
		log.Printf("Port:%d %s", k, strings.Join(v, " "))
		c := []string{"Check", strconv.Itoa(k)}
		section, err := ini.GetSection(strings.Join(c, "-"))
		if err != nil {
			continue
		}
		section.Add("hosts", strings.Join(v, " "))
	}

	if *iniFile {
		var b strings.Builder
		b.WriteString(groupName)
		b.WriteString(".ini")
		outfile := b.String()

		buf := bytes.NewBufferString("")
		ini.Write(buf)
		log.Printf("Writing %s", outfile)
		ini.WriteToFile(outfile)
	}

	if *jsonFile {
		var b strings.Builder
		b.WriteString(groupName)
		b.WriteString(".json")
		outfile := b.String()

		l, _ := json.MarshalIndent(&hosts, "", "  ")
		log.Printf("Writing %s", outfile)
		ioutil.WriteFile(outfile, l, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Done...")
}
