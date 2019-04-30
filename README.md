# Raven
Home Network Monitoring 

## Introduction
Home networks are getting a lot of *devices* on them, what is better known as IoT. There is monitoring software out there to give you a view of what is and isn't working on a network.  The issue is all this monitoring is geared around an **Enterprise**.  

I've worked with deploying Network Monitoring Systems (NMS) for over 20 years.   A lot of the enterpise feature are honestly not needed for a home network.  They systems are boated with:

* Notififcation systems
* Multi-level access controls
* Network isolation detection
* Service Level Agreement(SLA), enforcement and documenation
* Performace Archiving/History Database
* Intergration to Ticketing systems

What I want, is just a webpage I can go to and see if the reason my Xbox can't get on the net is because my router crashed or the switch it's plugged into is offline.  

## Goals
* Simple config, .ini files
* Lightweight process, Raspberry Pi Model 2B+ is the target platform
* Little Disk I/O, Raspberry Pi uses a flash card

## History
This is a rewrite of my prototype Kassandra, which was written in Python.  When I showed it to [MarkLlama](https://github.com/markllama), he suggest GoLang, and a name change, since Cassandra is Distrubuted Database from Apache.  So thinking about it and liking to name things after mythological stuff, I decided on Raven.  Odin has a pair of [Ravens](https://en.wikipedia.org/wiki/Huginn_and_Muninn) which worked as agents for him.  

So the last week or so I've been working on getitng the parts working in Go.  The structure is still basically the same.

## Changes
* The ini file now does not use type specifiers for the sections.  They are implied by the type of data in it.
  * Checks **MUST** have a `checkwith` key
  * Hosts **MUST** have a `hostname` key
* Hosts are now associated with Checks
* Hosts are not able to *modify* a check
  * if you want two checks that do almost the same thing, it's two checks.  
  * this was done to emulate how Nagios does things.  But after turning up a bunch of nagios checks in Kassandra, I decided for a home deplyoment it wasn't need.

## To Do
* **Documentation**
  * Shocker needs documenations
* **Configuration**
  * Add IPv6
  * Better error checking in the configuation file
  * Allow "Groups" to be monitored
  * Allow a host to be part of multiple groups?
* **Check Commands**
  * Add options to the Ping Command built-in.  
  * Add options to the Fping Command built-in.  
  * Implement calling Nagios Checks
  * Implement calling Nagios Checks via SSH (for remote machines)
  * Implement SNMP checks 
* **Logging**
  * reduce the chattiness of the server (add log levels)
  * add a central "logging" command at the end of a channel, because you know channels are cool
* **WebServer**
  * Make column sortable (DataTables.js?) 
  * add tabs for the "Groups"?
  * add logs visible on webserver
  * add thread status to webserver
