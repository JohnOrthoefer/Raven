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

// Nagios externals check command
import (
	"../ravenLog"
	"../ravenTypes"
	. "../ravenRunCommand"
	"fmt"
	"strings"
)

type nagiosOpts struct {
	prog     string
	progOpts []string
	addH     bool
	useDNS   bool
	//resplit   *regexp.Regexp
}

const checkMsgName = "Check nagios (plugin)"
var CheckName = "nagios"

func init() {
}

func InitCheck(kw ravenTypes.Kwargs) interface{} {
	var r interface{}
	rtn := new(nagiosOpts)
	rtn.prog = kw.GetKwargStr("program", "/usr/lib/monitoring-plugins/check_ping")
	rtn.progOpts = []string{"-w", "20,20%", "-c", "40,40%"}
	rtn.progOpts = kw.GetKwargStrA("options", rtn.progOpts)
	rtn.addH = kw.GetKwargBool("addhost", true)
	rtn.useDNS = kw.GetKwargBool("usedns", false)
	r = rtn
	return r
}

func RunCheck(he *ravenTypes.HostEntry, options interface{}) *ravenTypes.ExitReturn {
	e := new(ravenTypes.ExitReturn)
	opts := options.(*nagiosOpts)

	fullOpts := opts.progOpts
	if opts.addH {
		target := he.Hostname
		if he.IPv4 != "" { 
			target = he.IPv4
		}
    if opts.useDNS && he.Hostname != "" {
      target = hr.Hostname
    }
		fullOpts = append(fullOpts, "-H", target)
	}
	rtnExit, output := RunExternal(opts.prog, fullOpts...)

	switch rtnExit {
	case 0:
		fallthrough
	case 1:
		fallthrough
	case 2:
		e.Exit = rtnExit
		s := strings.Split(output, "|")
		switch len(s) {
		case 3:
			e.Long = s[2]
			fallthrough
		case 2:
			e.Perf = s[1]
			fallthrough
		case 1:
			e.Text = s[0]
		default:
			ravenLog.SendError(10, checkMsgName, "Failed to split output")
		}
	default:
		e.Exit = 3
		e.Text = "Nagios Unknown"
		e.Perf = ""
		e.Long = output
	}

	ravenLog.SendMessage(10, checkMsgName, fmt.Sprintf("%s(Nagios) exit:%d out=%s, perf=%s, long=%s", he.Hostname,
		e.Exit, e.Text, e.Perf, e.Long))
	return e
}
