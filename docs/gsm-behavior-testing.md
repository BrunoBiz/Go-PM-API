- [Intro](#intro)
  - [Start](#start)
  - [Stop](#stop)
  - [Restart](#restart)
  - [Details](#details)


# Intro
This document presents the outcomes of LinuxGSM and GSM, in response to each of the main options provided by the API, in a way to compare both and find common semantic meaning between them.

The purpose of these tests if to streamline how the response in the API will be handled, while still being agnostic to which one of the managers is currently hosting the server.

The data collected here will represent if a manager is idempotent, if a request will result in an error or not, and the server state before and after it's sent.

Below is a list of all the tests that will be performed to collect the data for this document:

| **Initial state** | **Operation** | **LinuxGSM**  | **GSM**       |
|-------------------|---------------|---------------|---------------|
| on                | start         | actual output | actual output |
| off               | start         | actual output | actual output |
| on                | stop          | actual output | actual output |
| off               | stop          | actual output | actual output |
| on                | restart       | actual output | actual output |
| off               | restart       | actual output | actual output |
| on                | details       | actual output | actual output |
| off               | details       | actual output | actual output |

Whenever possible, this is the data that will be collected from each test:
- Exit code
- Stdout / Stderr
- Server state after test
- Wheter the command succeded

## Start

### LinuxGSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: On
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
[ .... ] Starting untserver: grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Starting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[  OK  ] Starting untserver: LinuxGSM
```

Initial Status: On

```
- Successful: False
- Server state after: On
- Exit code: 2
- Stdout / Stderr:

Using LinuxGSM - Wrapper
[ .... ] Starting untserver: grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Starting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[ SKIP ] Starting untserver: LinuxGSM is already running
```

### GSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: On 
- Exit code: 0
- Stdout / Stderr:

""="GSM - START - SUCCESS: true"
""="GSM - START - COMMAND: /bin/tmux new -d -s mockServerTest '/opt/mockServerTest'"
""="GSM - START - MESSAGE: Server started"
```

Initial Status: On

```
- Successful: False
- Server state after: On 
- Exit code: 0
- Stdout / Stderr:

""="GSM - START - SUCCESS: false"
""="GSM - START - ERROR: exit status 1"
""="GSM - START - COMMAND: /bin/tmux new -d -s mockServerTest '/opt/mockServerTest'"
""="GSM - START - COMMAND RESULT: duplicate session: mockServerTest\n"
""="GSM - START - MESSAGE: Server is already running"
```


## Stop

### LinuxGSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: Off 
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
[ WARN ] Stopping untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[ SKIP ] Stopping untserver: LinuxGSM is already stopped
```

Initial Status: On

```
- Successful: True
- Server state after: Off
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
[ WARN ] Stopping untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[  OK  ] Stopping untserver: Graceful: CTRL+c: 4 ... OK
```

### GSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: False
- Server state after: Off
- Exit code: 1
- Stdout / Stderr:

""="Stopping server..."
""="GSM - STOP - SUCCESS: false"
""="GSM - STOP - ERROR: exit status 1"
""="GSM - STOP - COMMAND: /bin/tmux send-keys -t mockServerTest shutdown ENTER"
""="GSM - STOP - COMMAND RESULT: no server running on /tmp/tmux-1000/default\n"
""="GSM - STOP - MESSAGE: shutdown - Script failed to run"
```

Initial Status: On

```
- Successful: True
- Server state after: Off 
- Exit code: 0
- Stdout / Stderr:

""="Stopping server..."
Time elapsed: 15s / 300s

""="GSM - STOP - SUCCESS: true"
""="GSM - STOP - COMMAND: /bin/tmux send-keys -t mockServerTest shutdown ENTER"
""="GSM - STOP - MESSAGE: Server stopped"
```

## Restart

### LinuxGSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: On
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Restarting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[ WARN ] Stopping untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[ SKIP ] Stopping untserver: LinuxGSM is already stopped
[ .... ] Starting untserver: grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Starting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[  OK  ] Starting untserver: LinuxGSM
```

Initial Status: On

```
- Successful: True
- Server state after: On
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Restarting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[ WARN ] Stopping untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[  OK  ] Stopping untserver: Graceful: CTRL+c: 3 ... OK
[ .... ] Starting untserver: grep: /etc/apt/sources.list: No such file or directory
[ WARN ] Starting untserver: Configuration file missing!
/home/untserver/serverfiles/Servers/untserver/Config.json
[  OK  ] Starting untserver: LinuxGSM
```

### GSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: On 
- Exit code: 0
- Stdout / Stderr:

""="GSM - DETAILS - SUCCESS: true"
""="GSM - DETAILS - SERVER STATUS: false"
""="GSM - DETAILS - COMMAND: /bin/tmux ls"
""="GSM - DETAILS - COMMAND RESULT: no server running on /tmp/tmux-1000/default\n"
""="GSM - DETAILS - MESSAGE: No server running"
""="GSM - START - SUCCESS: true"
""="GSM - START - COMMAND: /bin/tmux new -d -s mockServerTest '/opt/mockServerTest'"
""="GSM - START - MESSAGE: Server started"
```

Initial Status: On

```
- Successful: True
- Server state after: On 
- Exit code: 0
- Stdout / Stderr:

""="GSM - DETAILS - SUCCESS: true"
""="GSM - DETAILS - SERVER STATUS: true"
""="GSM - DETAILS - COMMAND: /bin/tmux ls"
""="GSM - DETAILS - COMMAND RESULT: mockServerTest: 1 windows (created Sat Aug 29 02:30:06 2026)\n"
""="GSM - DETAILS - MESSAGE: Server running"
""="Stopping server..."
Time elapsed: 15s / 300s

""="GSM - STOP - SUCCESS: true"
""="GSM - STOP - COMMAND: /bin/tmux send-keys -t mockServerTest shutdown ENTER"
""="GSM - STOP - MESSAGE: Server stopped"
""="GSM - START - SUCCESS: true"
""="GSM - START - COMMAND: /bin/tmux new -d -s mockServerTest '/opt/mockServerTest'"
""="GSM - START - MESSAGE: Server started"
```

## Details

### LinuxGSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: Off 
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
grep: /etc/apt/sources.list: No such file or directory

Distro Details
==============================================================================================================================================================================================================================================================================================
Date:         Sat Aug 22 22:40:49 UTC 2026
Distro:       Debian GNU/Linux 13 (trixie)
Arch:         x86_64
Kernel:       6.17.13-1-pve
Hostname:     UnturnedServer
Environment:  lxc
Uptime:       0d, 21h, 28m
tmux:         3.5a
glibc:        2.41

Server Resource
==============================================================================================================================================================================================================================================================================================
CPU
Model:      AMD Ryzen 5 7430U with Radeon Graphics
Cores:      6
Frequency:  4340.882MHz
Avg Load:   0.53, 0.66, 0.60

Memory
Mem:       total  used   free   cached  available
Physical:  4.0GB  55MB   4.0GB  496MB   4.0GB
Swap:      4.0GB  4.8MB  4.0GB

Storage
Filesystem:  /dev/mapper/pve-vm--101--disk--0
Total:       20G
Used:        4.7G
Available:   14G

Network
IP:           0.0.0.0
Internet IP:  177.92.54.179

Unturned Resource Usage
==============================================================================================================================================================================================================================================================================================
CPU Used:  0%
Mem Used:  0%  0MB

Storage
Total:        3.0G
Serverfiles:  2.8G

Unturned Server Details
==============================================================================================================================================================================================================================================================================================
Server name:  LinuxGSM
App ID:       1110390
Server IP:    0.0.0.0:27015
Internet IP:  177.92.54.179:27015
Maxplayers:   20
Default map:  PEI
Game mode:    normal
Status:       STOPPED
Query Check:  https://ismygameserver.online/valve/177.92.54.179:27015

untserver Script Details
==============================================================================================================================================================================================================================================================================================
Script name:       untserver
LinuxGSM version:  v25.2.0
glibc required:    2.15
IFTTT alert:       off
Update on start:   off
User:              untserver
Location:          /home/untserver
Config file:       /home/untserver/serverfiles/Servers/untserver/Config.json (FILE MISSING)

Backups
==============================================================================================================================================================================================================================================================================================
No Backups created

Command-line Parameters
==============================================================================================================================================================================================================================================================================================
 ./Unturned_Headless.x86_64 -nographics -batchmode -bind 0.0.0.0 -port 27015 -maxplayers 20 -perspective first -mode normal -name LinuxGSM -map PEI -gslt  +InternetServer/untserver

Ports
==============================================================================================================================================================================================================================================================================================
Change ports by editing the parameters in:
/home/untserver/lgsm/config-lgsm/untserver

Useful port diagnostic command:
ss -tuplwn | grep Unturned_Headle

DESCRIPTION  PORT   PROTOCOL  LISTEN
Game         27015  udp       0
Query        27015  udp       0
Steam        27016  udp       0

Status: STOPPED
```

Initial Status: On

```
- Successful: True
- Server state after: On 
- Exit code: 0
- Stdout / Stderr:

Using LinuxGSM - Wrapper
grep: /etc/apt/sources.list: No such file or directory

Distro Details
==============================================================================================================================================================================================================================================================================================
Date:         Sat Aug 22 22:17:13 UTC 2026
Distro:       Debian GNU/Linux 13 (trixie)
Arch:         x86_64
Kernel:       6.17.13-1-pve
Hostname:     UnturnedServer
Environment:  lxc
Uptime:       0d, 21h, 5m
tmux:         3.5a
glibc:        2.41

Server Resource
==============================================================================================================================================================================================================================================================================================
CPU
Model:      AMD Ryzen 5 7430U with Radeon Graphics
Cores:      6
Frequency:  3836.092MHz
Avg Load:   0.85, 0.58, 0.38

Memory
Mem:       total  used   free   cached  available
Physical:  4.0GB  3.4GB  651MB  496MB   651MB
Swap:      4.0GB  4.8MB  4.0GB

Storage
Filesystem:  /dev/mapper/pve-vm--101--disk--0
Total:       20G
Used:        4.7G
Available:   14G

Network
IP:           0.0.0.0
Internet IP:  177.92.54.179

Unturned Resource Usage
==============================================================================================================================================================================================================================================================================================
CPU Used:  77.4%
Mem Used:  83.8%  3434MB

Storage
Total:        3.0G
Serverfiles:  2.8G

Unturned Server Details
==============================================================================================================================================================================================================================================================================================
Server name:  LinuxGSM
App ID:       1110390
Server IP:    0.0.0.0:27015
Internet IP:  177.92.54.179:27015
Maxplayers:   20
Default map:  PEI
Game mode:    normal
Status:       STARTED
Query Check:  https://ismygameserver.online/valve/177.92.54.179:27015

untserver Script Details
==============================================================================================================================================================================================================================================================================================
Script name:       untserver
LinuxGSM version:  v25.2.0
glibc required:    2.15
IFTTT alert:       off
Update on start:   off
User:              untserver
Location:          /home/untserver
Config file:       /home/untserver/serverfiles/Servers/untserver/Config.json (FILE MISSING)

Backups
==============================================================================================================================================================================================================================================================================================
No Backups created

Command-line Parameters
==============================================================================================================================================================================================================================================================================================
 ./Unturned_Headless.x86_64 -nographics -batchmode -bind 0.0.0.0 -port 27015 -maxplayers 20 -perspective first -mode normal -name LinuxGSM -map PEI -gslt  +InternetServer/untserver

Ports
==============================================================================================================================================================================================================================================================================================
Change ports by editing the parameters in:
/home/untserver/lgsm/config-lgsm/untserver

Useful port diagnostic command:
ss -tuplwn | grep Unturned_Headle

DESCRIPTION  PORT   PROTOCOL  LISTEN
Game         27015  udp       1
Query        27015  udp       1
Steam        27016  udp       1

Status: STARTED
```

### GSM <!-- omit from toc -->

Initial Status: Off

```
- Successful: True
- Server state after: Off
- Exit code: 0
- Stdout / Stderr:

""="GSM - DETAILS - SUCCESS: true"
""="GSM - DETAILS - SERVER STATUS: false"
""="GSM - DETAILS - COMMAND: /bin/tmux ls"
""="GSM - DETAILS - COMMAND RESULT: no server running on /tmp/tmux-1000/default\n"
""="GSM - DETAILS - MESSAGE: No server running"
```

Initial Status: On

```
- Successful: True
- Server state after: On 
- Exit code: 0 
- Stdout / Stderr:

""="GSM - DETAILS - SUCCESS: true"
""="GSM - DETAILS - SERVER STATUS: true"
""="GSM - DETAILS - COMMAND: /bin/tmux ls"
""="GSM - DETAILS - COMMAND RESULT: mockServerTest: 1 windows (created Sat Aug 29 02:30:38 2026)\n"
""="GSM - DETAILS - MESSAGE: Server running"
```