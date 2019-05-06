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

package ravenChecks

// ping check command

import (
	"../ravenLog"
	"../ravenTypes"
	"fmt"
	"regexp"
	"strconv"
)

type regmap map[string]*regexp.Regexp

var rePing regmap

type pingOpts struct {
	name     string
	pingProg string
	rttWarn  float64
	rttCrit  float64
	lossWarn int64
	lossCrit int64
	rttReg   *regexp.Regexp
	lossReg  *regexp.Regexp
	count    string
}

func init() {
	registerHandler("ping", PingInit, Ping)

	rePing = make(regmap)
	rePing["ping-rtt"], _ = regexp.Compile(`\d+\.?\d+/(\d+\.?\d+)/\d+\.?\d+/\d+\.?\d+`)
	rePing["ping-loss"], _ = regexp.Compile(`(\d+)\% packet loss`)

	registerHandler("fping", FpingInit, Ping)
	rePing["fping-rtt"], _ = regexp.Compile(`min/avg/max = \d+\.?\d+/(\d+\.?\d+)/\d+\.?\d+`)
	rePing["fping-loss"], _ = regexp.Compile(`xmt/rcv/\%loss = \d+/\d+/(\d+)\%,`)
}

func pingComm(kw ravenTypes.Kwargs) *pingOpts {
	rtn := new(pingOpts)
	rtn.rttWarn = kw.GetKwargFloat("rtt_warn", 20.0)
	rtn.lossWarn = kw.GetKwargInt("loss_warn", 20)
	rtn.rttCrit = kw.GetKwargFloat("rtt_crit", 30.0)
	rtn.lossCrit = kw.GetKwargInt("loss_crit", 40)
	rtn.count = kw.GetKwargStr("count", "5")
	return rtn
}

func PingInit(kw ravenTypes.Kwargs) interface{} {
	var r interface{}
	rtn := pingComm(kw)
	rtn.name = "Ping"
	rtn.pingProg = kw.GetKwargStr("program", "/usr/bin/ping")
	rtn.rttReg = rePing["ping-rtt"]
	rtn.lossReg = rePing["ping-loss"]
	r = rtn
	return r
}

func FpingInit(kw ravenTypes.Kwargs) interface{} {
	var r interface{}
	rtn := pingComm(kw)
	rtn.name = "Fping"
	rtn.pingProg = kw.GetKwargStr("program", "/usr/bin/fping")
	rtn.rttReg = rePing["fping-rtt"]
	rtn.lossReg = rePing["fping-loss"]
	r = rtn
	return r
}

func Ping(he *ravenTypes.HostEntry, options interface{}) *ravenTypes.ExitReturn {
	e := new(ravenTypes.ExitReturn)
	opts := options.(*pingOpts)

	target := he.Hostname
	if he.IPv4 != "" {
		target = he.IPv4
	}
	rtnExit, output := runExternal(opts.pingProg, "-c", opts.count, target)

	switch rtnExit {
	case 0:
		rtt := opts.rttReg.FindAllStringSubmatch(output, -1)
		pls := opts.lossReg.FindAllStringSubmatch(output, -1)
		rttAvg, _ := strconv.ParseFloat(rtt[0][1], 32)
		loss, _ := strconv.ParseInt(pls[0][1], 10, 32)
		e.Exit = 0
		e.Text = fmt.Sprintf("%s Okay", opts.name)
		e.Perf = fmt.Sprintf("RTT Average: %4.2f, Loss: %d", rttAvg, loss)
		e.Long = fmt.Sprintf("Count:%s Warn:%d,%4.2f Crit:%d,%4.2f",
			opts.count, opts.lossWarn, opts.rttWarn,
			opts.lossCrit, opts.rttCrit)
		if opts.rttWarn < rttAvg || opts.lossWarn < loss {
			e.Exit = 1
			e.Text = fmt.Sprintf("%s Warning", opts.name)
			e.Long = fmt.Sprintf("WARNING %4.2f < %4.2f or %d < %d",
				opts.rttWarn, rttAvg, opts.lossWarn, loss)
		}
		if opts.rttCrit < rttAvg || opts.lossCrit < loss {
			e.Exit = 2
			e.Text = fmt.Sprintf("%s Critical", opts.name)
			e.Long = fmt.Sprintf("CRITICAL %4.2f < %4.2f or %d < %d",
				opts.rttCrit, rttAvg, opts.lossCrit, loss)
		}
	default:
		e.Exit = 3
		e.Text = fmt.Sprintf("%s unknown", opts.name)
		e.Perf = ""
		e.Long = output
	}

	ravenLog.SendMessage(10, fmt.Sprintf("Check %s", opts.name), fmt.Sprintf("%s(Ping) exit:%d out=%s, perf=%s", he.Hostname,
		e.Exit, e.Text, e.Perf))
	return e
}
