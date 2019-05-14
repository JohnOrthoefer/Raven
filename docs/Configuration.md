# Configuration
## Auto configuration
`nmapparse` is used for inital configuration/network descovery.  This program attempts to scan your local subnet with `/usr/bin/nmap` and then parses the XML output.  

* **Discovery**
These options control how network discovery is done. 
  * `-net NETWORK/CIDR` the network to scan.  If you don't provide one, the program *makes a guess*
* **Input**
  * `-xml FILENAME` provide a file instead of allowing it to run nmap.
* **Output**
If you do not provide at either `-json` or `-ini` the program only prints to the terminal a summary of what it found. 
  * `-json` output GROUPNAME.json file with all the found hosts
  * `-ini` output GROUPNAME.ini file for all the found hosts.  This file will have some basic checks enabled for all all hosts, currently `ping` for all hosts, and for hosts with the following ports open-
    * port 22 /usr/lib/monitoring-plugins/check_ssh
    * port 80 /usr/lib/monitoring-plugins/check_http
    * port 443 /usr/lib/monitoring-plugins/check_http *to check certificate*
  * `-skel base.ini` the file with the base checks in it for.ini 
  * `-disabled` mark all hosts as `enabled=false`
  * `-group Internal-LAN` this is the tag all the hosts found will be added to
  * `-dhcp 100-200` in the ini file do not add an `ipv4=` entry if the least sigificate octet falls in this range.   This option won't work correctly if your CIDR block is larger is /24. To disable, set the range to `256-256`.


## Deamon 
* `-config` points the daemon at the configration file to use. The default is `../etc/raven.ini` However, you can name and place the file where ever you want.
* `-port` what port to use for the webserver.  **Default** `:8000`
* `-workers` how many *worker threads* to spawn.  Understand these are go routines so it will not map 1 to 1 to what your OS shows.

## Monitoring file
All the monitoring goes into a single ini file.   

There are two section types:
* **Checks** Monitors to run, to be considered a check, the section **must have** a `checkwith =` key in it.  Currently there are 3 built in monitor/checks:
  * **ping** - Simply use the ping command to check for a host reachablity
  * **nagios** - Runs a nagios check.  Nagios is another NMS system, which has a rich set of checks written as a stand alone programs.  Instead of needing to *reinvent* all of them I choose to leveage them.
  * **viassh** - This runs `ssh` to reach out and check another host, mostly tested with nagios-check commands on the far end. 

* **Hosts** - this section describes a host.  A host either needs a `hostname =` or an `ipv4 =` key in it.  It is not excusive a host can have both.
  * All hosts belong to a *group* which is currently only used for display grouping. 

Each section is identified by a unique name inside square brackets, `[]`.  This name is uses as the friendly/display name of the check or host.

## Checks
### All checks
There are a few Keys that are common to every check.  
* `checkwith = ping` what check to run
* `interval = 1m30s 1m 30s 30s` how often to check, based on the return code of the last run.  space seperate, it understand h=hours, m=minutes, s=seconds.  The 4 values are okay, warning, critical, unknown.
* `threshold = 5` number of times a check is required to *change* status.  1 means if it fails the check changes, this is used for checks like certificate experation, if the cert has expire, it is expired.
* `hosts = ` the hosts that this check needs to be run for.
* `enabled = true` used to disable this check completely, without deleting it from the configuration.

### `checkwith = ping`
* `program = /usr/bin/ping` The fully qualified path to the ping program on your system. 
* `count = 5` This will be passed to the ping program as `-c N` for the number of times to send a ping each check.
* `rtt_warn = 20.0` and `rtt_crit = 30.0` The round trip time(rtt) in ms if exceeded on average to consider warning or critical.
* `loss_warn = 20` or `loss_crit = 40` The percentage of loss packets to consider warning or critical.
* `usefping = false` when set to true it will call `/usr/bin/fping` instead

### `checkwith = nagios`
This runs a nagios check.  
* `program = /usr/lib/monitoring-plugins/check_ping` The program to run
* `options = -w 20,20% -c 40,40%` options to pass to the program
* `addhost = true` If you want `-H [HOST/IP]` it will pass the IP address if there is one then it falls back to passing the `hostname`
* `usedns = false` the default is to use the IP if available, if you want to use the hostname always set this to true

### `checkwith = viassh`
This runs a check via ssh.  Typically a remote nagios check.
* `ssh = /usr/bin/ssh` 
* `sshoptions = ""` options passed before the hostname and executable.  This should not be nessary since you can do set most options in the `~/.ssh/config`
* `program = /usr/lib/monitoring-plugins/check_ping` The program to run on the remote host
* `options = -w 20,20% -c 40,40%` options to pass to the remote program 
* `addhost = true` If you want `-H [HOST/IP]` it will pass the IP address if there is one then it falls back to passing the `hostname`
* `usedns = false` the default is to use the IP if available, if you want to use the hostname always set this to true

## Host
* Host address
  * `host = ` DNS Hostname of the machine
  * `ipv4 = ` IPv4 address of the machine
  * if both exist the v4 address is used.  Some check commands can override this behavior
* `group = Internal-LAN` the default display group name
* `enabled = true` used to disable a host without removing it from the configuration

