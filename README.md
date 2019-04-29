# Raven
Home Network Monitoring 

## Introduction
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
* **Configuration**
  * Add IPv6
* **Check Commands**
  * Add options to the Ping Command built-in.  
  * Add options to the Fping Command built-in.  
  * Implement calling Nagios Checks
  * Implement calling Nagios Checks via SSH (for remote machines)
  * Implement SNMP checks 
* **Logging**
  * reduce the chattiness of the server (add log levels)
* **WebServer**
  * Make column sortable (DataTables.js?) 
  * add tabs for the "Groups"?
  * add logs visible on webserver
  * add thread status to webserver
