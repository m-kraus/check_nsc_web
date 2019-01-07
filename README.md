# About *check_nsc_web*

*check_nsc_web* collects check results from NSClient++ agents using its brand-new REST API. It is an alternative to check_nrpe et al.
*check_nsc_web* can be used with any monitoring tool, that can use Nagios compatible plugins.

To be easily portable, *check_nsc_web* is written in Go. Binary builds for Linux, Windows and MacOS are available in the ```build``` subdirectory.

*check_nsc_web* is released under the GNU GPL v3.

## Usage examples
* Alive check
```
go run check_nsc_web.go -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443"
OK: NSClient API reachable on https://localhost:8443
```

* CPU usage
```
check_nsc_web -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_cpu
OK: CPU load is ok.|'total 5m'=16%;80;90 'total 1m'=8%;80;90 'total 5s'=8%;80;90
```
* CPU usage with threshodlds
```
check_nsc_web -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_cpu show-all "warning=load > 75" "critical=load > 90"
OK: 5m: 1%, 1m: 0%, 5s: 0%|'total 5m'=1%;75;90 'total 1m'=0%;75;90 'total 5s'=0%;75;90
```

* Service status
```
check_nsc_web -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_service "service=BvSshServer"
OK: All 1 service(s) are ok.|'BvSshServer'=4;0;0
```

* Complex eventlog check
```
check_nsc_web -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_eventlog "file=system" "filter=id=8000" "crit=count>0" "detail-syntax=\${message}" show-all "scan-range=-900m"
OK: No entries found|'count'=0;0;0 'problem_count'=0;0;0
```

## Program help
```
Usage of ./check_nsc_web:

  check_nsc_web is a REST client for the NSClient++ webserver for querying
  and receiving check information over HTTPS.

  Copyright 2016 Michael Kraus <Michael.Kraus@consol.de>

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.

  Example:
  check_nsc_web -p "password" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_cpu

  Usage:
  check_nsc_web [options] [NSClient query parameters]

  check_nsc_web can and should be built with CGO_ENABLED=0

  Options:
  -V	Print program version.
  -a string
    	API version of NSClient++ (legacy or 1). (default "legacy")
  -f int
    	Round performance data float values to this number of digits. (default -1)
  -j	Print out JOSN response body.
  -k	Insecure mode - skip TLS verification.
  -p string
    	NSClient++ webserver password.
  -t int
    	Connection timeout in seconds, defaults to 10. (default 10)
  -u string
    	NSCLient++ URL, for example https://10.1.2.3:8443.
  -v	Enable verbose output.
  -x string
    	Extra text to appear in output.
  -l string
        NSClient++ webserver login. (default "admin")
```
