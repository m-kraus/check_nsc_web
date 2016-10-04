TODO

## Usage examples
* Alive check
```
go run nscrestc.go -k -p "password from nsclient.ini" -u "https://localhost:8443"
OK: NSClient API reachable on https://localhost:8443
```

* CPU usage
```
go run nscrestc.go -k -p "password from nsclient.ini" -u "https://localhost:8443" check_cpu
OK: CPU load is ok.|'total 5m'=16%;80;90 'total 1m'=8%;80;90 'total 5s'=8%;80;90
```