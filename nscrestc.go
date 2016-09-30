package main

// TODO
// - specify cert
// - specify ciphers
// - GNU preamble and copyright information
// - usage header
// - make Perf optional -> no perfdata pipe gets printed on Unknown for example
// FIXME

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
	"bytes"
	"strings"
)

type Query struct {
	HeaGder struct {
		SourceID string `json:"source_id"`
	} `json:"header"`
	Payload []struct {
		Command string `json:"command"`
		Lines   []struct {
			Message string `json:"message"`
			Perf    []struct {
				Alias    string `json:"alias"`
				IntValue struct {
					Value    *float64 `json:"value,omitempty"`
					Unit     *string  `json:"unit,omitempty"`
					Warning  *float64 `json:"warning,omitempty"`
					Critical *float64 `json:"critical,omitempty"`
					Minimum  *float64 `json:"mininum,omitempty"`
					Maximum  *float64 `json:"maximum,omitempty"`
				} `json:"int_value"`
			} `json:"perf"`
		} `json:"lines"`
		Result string `json:"result"`
	} `json:"payload"`
}

func main() {
	var flagURL string
	var flagPassword string
	var flagTimeout int
	var flagVerbose bool
	var flagInsecure bool

	flag.StringVar(&flagURL, "u", "", "NSCLient++ URL, for example https://10.1.2.3:8443.")
	flag.StringVar(&flagPassword, "p", "", "NSClient++ webserver password.")
	flag.IntVar(&flagTimeout, "t", 10, "Connection timeout in seconds, defaults to 10.")
	flag.BoolVar(&flagVerbose, "v", false, "Enable verbose output.")
	flag.BoolVar(&flagInsecure, "k", false, "Insecure mode - skip TLS verification.")

	ReturncodeMap := map[string]int{
		"OK":       0,
		"WARNING":  1,
		"CRITICAL": 2,
		"UNKNOWN":  3,
	}

	required := []string{"u", "p"}
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "UNKNOWN: Missing required -%s argument\n", req)
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(3)
		}
	}

	var Url *url.URL
	Url, err := url.Parse(flagURL)
	if err != nil {
		fmt.Println("UNKNOWN: " + err.Error())
		os.Exit(3)
	}

	if len(flag.Args()) == 0 {
		Url.Path += "/"
	} else if len(flag.Args()) == 1 {
		Url.Path += "/query/" + flag.Arg(0)
	} else {
		Url.Path += "/query/" + flag.Arg(0)
		parameters := url.Values{}
		for i, a := range flag.Args() {
			if i == 0 {
				continue
			}
			p := strings.SplitN(a, "=", 2)
			if len(p) == 1 {
				// FIXME it is unclear if a trailing "=" e.g. on show-all can lead to errors
				parameters.Add(p[0], "")
			} else {
				parameters.Add(p[0], p[1])
			}
			if err != nil {
				fmt.Println("UNKNOWN: " + err.Error())
				os.Exit(3)
			}
		}
		Url.RawQuery = parameters.Encode()
	}

	var hTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: flagInsecure,
		},
		TLSHandshakeTimeout: time.Second * time.Duration(flagTimeout),
	}
	var hClient = &http.Client{
		Timeout:   time.Second * time.Duration(flagTimeout),
		Transport: hTransport,
	}

	req, err := http.NewRequest("GET", Url.String(), nil)
	if err != nil {
		fmt.Println("UNKNOWN: " + err.Error())
		os.Exit(3)
	}
	req.Header.Add("password", flagPassword)

	if flagVerbose {
		dumpreq, _ := httputil.DumpRequestOut(req, true)
		fmt.Printf("REQUEST:\n%q\n", dumpreq)
	}
	res, err := hClient.Do(req)
	if err != nil {
		fmt.Println("UNKNOWN: " + err.Error())
		os.Exit(3)
	}
	defer res.Body.Close()

	if flagVerbose {
		dumpres, _ := httputil.DumpResponse(res, true)
		fmt.Printf("RESPONSE:\n%q\n", dumpres)
	}

	if len(flag.Args()) == 0 {
		fmt.Println("OK: NSClient API reachable on " + flagURL)
		os.Exit(0)
	} else {
		QueryResult := new(Query)
		json.NewDecoder(res.Body).Decode(QueryResult)

		// FIXME as payload is a slice, does it have to be iterable ?
		Result := QueryResult.Payload[0].Result

		var nagiosMessage string
		var nagiosPerfdata bytes.Buffer

		// FIXME as payload is a slice, does it have to be iterable ?
		// FIXME how to iterate the slice of lines safely ?
		for _, l := range QueryResult.Payload[0].Lines {

			nagiosMessage = strings.TrimSpace(l.Message)

			val := ""
			uni := ""
			cri := ""
			war := ""
			min := ""
			max := ""
			for _, p := range l.Perf {
				// REFERENCE 'label'=value[UOM];[warn];[crit];[min];[max]
				if p.IntValue.Value != nil {
					val = strconv.FormatFloat(*(p.IntValue.Value), 'f', -1, 64)
					nagiosPerfdata.WriteString(" '" + p.Alias + "'=" + val)
				} else {
					continue
				}
				if p.IntValue.Unit != nil {
					uni = *(p.IntValue.Unit)
					nagiosPerfdata.WriteString(uni)
				}
				if p.IntValue.Warning != nil {
					war = strconv.FormatFloat(*(p.IntValue.Warning), 'f', -1, 64)
					nagiosPerfdata.WriteString(";" + war)
				}
				if p.IntValue.Critical != nil {
					cri = strconv.FormatFloat(*(p.IntValue.Critical), 'f', -1, 64)
					nagiosPerfdata.WriteString(";" + cri)
				}
				if p.IntValue.Minimum != nil {
					min = strconv.FormatFloat(*(p.IntValue.Minimum), 'f', -1, 64)
					nagiosPerfdata.WriteString(";" + min)
				}
				if p.IntValue.Maximum != nil {
					max = strconv.FormatFloat(*(p.IntValue.Maximum), 'f', -1, 64)
					nagiosPerfdata.WriteString(";" + max)
				}
			}
		}

		if nagiosPerfdata.Len() == 0 {
			fmt.Println(nagiosMessage)
		} else {
			fmt.Println(nagiosMessage + "|" + strings.TrimSpace(nagiosPerfdata.String()))
		}
		os.Exit(ReturncodeMap[Result])
	}

}
