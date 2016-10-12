/*
  nscrestc

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
*/

package main

// TODO
// - specify cert
// - specify ciphers
// - usage header

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/griesbacher/check_x"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

//Query represents the nsclient response
type Query struct {
	Header struct {
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

	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})
	for _, req := range []string{"u", "p"} {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "UNKNOWN: Missing required -%s argument\n", req)
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(3)
		}
	}

	urlStruct, err := url.Parse(flagURL)
	check_x.ExitOnError(err)

	if len(flag.Args()) == 0 {
		urlStruct.Path += "/"
	} else if len(flag.Args()) == 1 {
		urlStruct.Path += "/query/" + flag.Arg(0)
	} else {
		urlStruct.Path += "/query/" + flag.Arg(0)
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
		}
		urlStruct.RawQuery = parameters.Encode()
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

	req, err := http.NewRequest("GET", urlStruct.String(), nil)
	check_x.ExitOnError(err)

	req.Header.Add("password", flagPassword)

	if flagVerbose {
		dumpreq, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Printf("REQUEST-ERROR:\n%s\n", err.Error())
		}
		fmt.Printf("REQUEST:\n%q\n", dumpreq)
	}
	res, err := hClient.Do(req)
	check_x.ExitOnError(err)
	defer res.Body.Close()

	if flagVerbose {
		dumpres, err := httputil.DumpResponse(res, true)
		if err != nil {
			fmt.Printf("RESPONSE-ERROR:\n%s\n", err.Error())
		}
		fmt.Printf("RESPONSE:\n%q\n", dumpres)
	}

	if len(flag.Args()) == 0 {
		check_x.Exit(check_x.OK, "NSClient API reachable on " + flagURL)
	} else {
		queryResult := new(Query)
		err = json.NewDecoder(res.Body).Decode(queryResult)
		check_x.ExitOnError(err)

		if len(queryResult.Payload) == 0 {
			if flagVerbose {
				fmt.Printf("QUERY RESULT:\n%+v\n", queryResult)
			}
			check_x.Exit(check_x.Unknown, "The resultpayload size is 0")
		}
		result := queryResult.Payload[0].Result

		var nagiosMessage string

		// FIXME how to iterate the slice of lines safely ?
		for _, l := range queryResult.Payload[0].Lines {
			nagiosMessage += strings.TrimSpace(l.Message)

			for _, p := range l.Perf {
				if p.IntValue.Value != nil && p.Alias != "" {
					perf := check_x.NewPerformanceData(p.Alias, *(p.IntValue.Value))

					if p.IntValue.Unit != nil {
						perf.Unit(*(p.IntValue.Unit))
					}
					if p.IntValue.Warning != nil {
						w, err := check_x.NewThreshold(strconv.FormatFloat(*(p.IntValue.Warning), 'f', -1, 64))
						check_x.ExitOnError(err) //Should never happen
						perf.Warn(w)
					}
					if p.IntValue.Critical != nil {
						c, err := check_x.NewThreshold(strconv.FormatFloat(*(p.IntValue.Critical), 'f', -1, 64))
						check_x.ExitOnError(err) //Should never happen
						perf.Crit(c)
					}
					if p.IntValue.Minimum != nil {
						perf.Min(*(p.IntValue.Minimum))
					}
					if p.IntValue.Maximum != nil {
						perf.Max(*(p.IntValue.Maximum))
					}
				}
			}
		}
		check_x.Exit(check_x.StateFromString(result), nagiosMessage)
	}

}
