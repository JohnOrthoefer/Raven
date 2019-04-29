# Raven
Home Network Monitoring 

## Introduction
This is a rewrite of my prototype Kassandra, which was written in Python.  When I showed it to MarkLlama he suggest GoLang, and a name change, since Cassandra is Distrubuted Database from Apache.  So thinking about it and liking to name things after mythological stuff, I decided on Raven.  Odin had a pair of Ravens (Huginn/Thought and Muninn/Memory) which worked as agents for him.  

So the last week or so I've been working on getitng the parts working in Go.  The structure is still basically the same.

## Changes
* The ini file now does not use type specifiers for the sections.  They are implied by the type of data in it.
  * Checks MUST HAVE a 'checkwith' key
  * Hosts MUST HAVE a 'hostname' key
* Hosts are now associated with Checks
* Hosts are not now able to "modify" a check
  * if you want two checks that do almost the same thing, it's two checks.  
  * this was done to emulate how Nagios does things.  But after turning up a bunch of nagios checks in Kassandra, I decided for a home deplyoment it wasn't need.

## To Do
* IPv6 - I have IPv6 to monitor so this will happen
* Finish the built-in Ping check (although it uses the external Ping Command, so the check doesn't have to be run as root.)
* Implement calling Nagios Checks
* Implement calling Nagios Checks via SSH (for remote machines)
* Implement SNMP checks 
* Add a logs pages with a circular buffer 
  * reduce the chattiness of the server (add log levels)
