PingScan
========

PROOF OF CONCEPT

Fast concurent ping with JSON output. 

TODO
----

- missing timeouts
- forcing v6


Usage
-----

**pingscan google.com yahoo.com astalavista.com ...**

Rights
------

Linux rights by default don't allow regular users to send ICMP packets. This 
can be fixed either running pingscan by root 

* you can run pingscan by **root**

* or you can allow your user group to manipulate raw sockets with command: **sudo sysctl net.ipv4.ping_group_range="0   1000"** where "0   1000" is range of IDs of groups 

* or you can set setuid bit to pingscan with: **sudo chown root:root pingscan && sudo chmod +s pingscan**
