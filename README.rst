.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0

This project consists of three main components: API, orchestration and
adapter. Those three components communicate with each through RPC
mechanism (using jsonRPC, only support point-point connection now).

API module manages the request about storage resources, such as volumes,
databases, file systems, policys and so forth.

Orchestration module has three roles:
1. Handles the request from API module.
2. Collects the statistics (connection information, feature and so on) of
   storage resources through adapter module.
3. Orchestrates storage resources and shows appropriate resources to users
   according to scenarios.

Adapter module contains plugins to integrate open source projects (such
as Cinder, Manila, Swift and so on) and enterprise projects (such as
OceanStor DJ).

Besides, log module provides Log function for system and it can debug
the error when the system breaks down. And we can directly get access to
the system by using CLI with cmd module. Lastly, there are some test cases
in test module to test the functionality of the system.  
