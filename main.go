package main

// geoUR : read ips from stdin, display for Country and ASN

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

// History :
// v1.0 combine geoip2 Country and ASN,
//
// it doesn't resolve ip anymore
//	use "| adnsresfilter -ua" if resolution is needed
//

var (
	maxrequests = 128
	// Version given by git tag via Makefile
	Version string
)

// printVersionUsage only print Version and Usage anq exit
func printVersionUsage() {
	fmt.Printf("%s %s embedded geoip2 databases versions\n\t- %s\n\t- %s\n\n",
		os.Args[0], Version, AssetCountryName, AssetASNname)
	fmt.Printf("Usage: %s ip or stdin input\n", os.Args[0])
	os.Exit(0)
}
func main() {

	// embedded values are build from dynembedded.go
	ASN := newgeoip2Reader(embeddedASN)
	Country := newgeoip2Reader(embeddedCountry)

	switch {
	case len(os.Args) == 1:
		// Readstdin
		readandprintbulk(Country, ASN)

	case (len(os.Args) == 2) && os.Args[1] == "-h":
		printVersionUsage()

	case len(os.Args) == 2:
		// Single input as parameter
		fmt.Printf("%s\n", parseandprint(os.Args[1], Country, ASN))

	default:
		printVersionUsage()
	}

	os.Exit(0)
}

// newgeoip2Reader wrapper from []byte to *geoip2.Reader
func newgeoip2Reader(Asset []byte) *geoip2.Reader {
	db, err := geoip2.FromBytes(Asset)
	if err != nil {
		log.Printf("ERR geoip2.FromBytes()")
		log.Fatal(err)
	}
	return db
}

// parseandprint return string with ASN and Country from ip string
func parseandprint(ips string, dbCountry, dbASN *geoip2.Reader) string {
	ip := net.ParseIP(ips)
	if ip == nil {
		log.Printf("ERR net.ParseIP(%s)", ips)
		return ""
	}

	record, err := dbCountry.Country(ip)
	if err != nil || record.Country.Names["en"] == "" {
		log.Printf("ERR dbCountry.Country(\"%s\")", ips)
		return ""
	}

	ASN, err := dbASN.ASN(ip)
	if err != nil {
		log.Printf("ERR dbASN.ASN(%s)", ips)
		return ""
	}

	output := fmt.Sprintf("%s (%s / AS%d %s)", ips,
		record.Country.Names["en"],
		ASN.AutonomousSystemNumber, ASN.AutonomousSystemOrganization)

	return output
}

// readandprintbulk read stdin and printout results
// channel was used when parseandprint resolved ip to name
// wg was used to limit concurrent requests
func readandprintbulk(dbCountry, dbASN *geoip2.Reader) {
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
			out := parseandprint(myline, dbCountry, dbASN)
			if len(out) > 0 {
				fmt.Printf("%s\n", out)
			}
			<-mychan
			mywg.Done()
		}(line, &wg, limitChan)

	}
	wg.Wait()
}
