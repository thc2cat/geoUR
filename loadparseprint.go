package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var (
	// embedded values are generated by init.sh in init.go
	dbASN     = newgeoip2Reader(embeddedASN)
	dbCountry = newgeoip2Reader(embeddedCountry)
)

// newgeoip2Reader wrapper from []byte to *geoip2.Reader
func newgeoip2Reader(Asset []byte) *geoip2.Reader {
	db, err := geoip2.FromBytes(Asset)
	if err != nil {
		log.Fatal("ERR geoip2.FromBytes() ", err)
	}
	return db
}

// parseandprint return string with ASN and Country from ip string
func parseandprint(ips string) string {
	ip := net.ParseIP(ips)
	if ip == nil {
		return fmt.Sprintf("%s (ERR net.ParseIP / AS0ERR)", ips)
	}

	record, err := dbCountry.Country(ip)
	if err != nil || record.Country.Names["en"] == "" {
		return fmt.Sprintf("%s (ERR geoip2.Country / AS0ERR)", ips)
	}

	ASN, err := dbASN.ASN(ip)
	if err != nil {
		return fmt.Sprintf("%s (ERR geoip2.ASN / AS0ERR)", ips)
	}

	return fmt.Sprintf("%s (%s / AS%d %s)", ips,
		record.Country.Names["en"],
		ASN.AutonomousSystemNumber, ASN.AutonomousSystemOrganization)
}

// Trying ot pass a generic ParsePrint function to readandprintbulk
type FuncStringtoString func(string) string

// readandprintbulk read stdin and printout results
// channel was used when parseandprint resolved ip to name
// wg was used to limit concurrent requests
func readandprintbulk(doit FuncStringtoString) {
	var line string
	var wg sync.WaitGroup
	var limitChan = make(chan bool, maxrequests)
	cached := make(map[string]bool)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if cached[line] {
			continue
		}
		cached[line] = true
		wg.Add(1)
		limitChan <- true // will block after maxrequests

		go func(myline string, mywg *sync.WaitGroup, mychan chan bool) {
			out := doit(myline)
			if len(out) > 0 {
				fmt.Printf("%s\n", out)
			}
			<-mychan
			mywg.Done()
		}(line, &wg, limitChan)

	}
	wg.Wait()
}
