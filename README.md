# About *nscrestc*

*nscrestc* collects check results from NSClient++ agents using its brand-new REST API. It is an alternative to check_nrpe et al.
*nscrestc* can be used with any monitoring tool, that can use Nagios compatible plugins.

To be easily portable, *nscrestc* is written in Go. Binary builds for Linux, Windows and MacOS are available in the ```build``` subdirectory.
 
*nscrestc* is released under the GNU GPL v3.

## Usage examples
* Alive check
```
go run nscrestc.go -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443"
OK: NSClient API reachable on https://localhost:8443
```

* CPU usage
```
nscrestc -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_cpu
OK: CPU load is ok.|'total 5m'=16%;80;90 'total 1m'=8%;80;90 'total 5s'=8%;80;90
```
* CPU usage with threshodlds
```
nscrestc -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_cpu show-all "warning=load > 75" "critical=load > 90"
OK: 5m: 1%, 1m: 0%, 5s: 0%|'total 5m'=1%;75;90 'total 1m'=0%;75;90 'total 5s'=0%;75;90
```

* Service status
```
nscrestc -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_service "service=BvSshServer"
OK: All 1 service(s) are ok.|'BvSshServer'=4;0;0
```

* Complex eventlog check
```
nscrestc -k -p "password from nsclient.ini" -u "https://<SERVER_RUNNING_NSCLIENT>:8443" check_eventlog "file=system" "filter=id=8000" "crit=count>0" "detail-syntax=\${message}" show-all "scan-range=-900m"
OK: No entries found|'count'=0;0;0 'problem_count'=0;0;0
```

## Program help
```
Usage of ./nscrestc:
  -f  int
      Round performance data float values to this number of digits
  -k	Insecure mode - skip TLS verification.
  -p string
    	NSClient++ webserver password.
  -t int
    	Connection timeout in seconds, defaults to 10. (default 10)
  -u string
    	NSCLient++ URL, for example https://10.1.2.3:8443.
  -v	Enable verbose output.
```