# Configuration
## Auto configuration
Included is `nmapparse`.  This program attempts to scan your local subnet with nmap and then parses the XML output.  You can control which network it probes, and even provide your own nmap XML for it parse and output a ini file. 

## Deamon 
* `-config` - points the daemon at the configration file to use. The default is `../etc/raven.ini` However, you can name and place the file where ever you want.
* `-port` - what port to use for the webserver.  **Default** `:8000`
* `-workers` - how many *worker threads* to spawn.  Understand these are go routines so it will not map 1 to 1 to what your OS shows.

## Monitoring file
All the monitoring goes into a single ini file.   

There are two section types:
* **Checks** Monitors to run, to be considered a check, the section **must have** a `checkwith =` key in it.  Currently there are 4 built in monitor/checks:
  * **ping** - Simply use the ping command to check for a host reachablity
  * **fping** - This is a secondary function of the ping it uses the ``fping`` program instead of the regular OS `ping` command.
  * **nagios** - Runs a nagios check.  Nagios is another NMS system, which has a rich set of checks written as a stand alone programs.  Instead of needing to *reinvent* all of them I choose to leveage them.
  * **viassh** - This runs `ssh` to reach out and check another host, mostly tested with nagios-check commands on the far end. 

* **Hosts** - this section describes a host.  A host either needs a `hostname =` or an `ipv4 =` key in it.  It is not excusive a host can have both.
  * All hosts belong to a *group* which is currently only used for display grouping. 

Each section is identified by a unique name inside square brackets, `[]`.  This name is uses as the friendly/display name of the check or host.

## Checks
### All checks
there are a few Keys that are common to every check.  

### `checkwith = ping`
* `program = /usr/bin/ping` The fully qualified path to the ping program on your system. 
* `count = 5` This will be passed to the ping program as `-c N` for the number of times to send a ping each check.
* `rtt_warn = 20.0` and `rtt_crit = 30.0` The round trip time(rtt) in ms if exceeded on average to consider warning or critical.
* `loss_warn = 20` or `loss_crit = 40` The percentage of loss packets to consider warning or critical.

### `checkwith = fping`
This is the same as ping, with the exception of the program to run which is by default `program = /usr/bin/fping`.  It does change 

### `checkwith = nagios`
This 
