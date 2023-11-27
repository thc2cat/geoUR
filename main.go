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

// History :
// v1.0 combine geocli and geoASN, doesn't resolve ip anymore
//
//	use "| adnsresfilter -ua" if resolution is needed
//

var (
	maxrequests = 128
	// Version given by git tag via Makefile
	Version string
	// Puts your privates networks at the end of this regexg
	// this allow a tag "local" to be printed for ip exclusion
	// Privates  = regexp.MustCompile(`^(10|172\.(1[6789]|2[0-9]|3[01])|192\.168|193.51\.(2[456789]|3[0-9]|4[12]))\.`)
	FullPrint bool
)

// printVersionUsage only print Version and Usage anq exit
func printVersionUsage() {
	fmt.Printf("=== %s %s embedded geoip2 databases\n\t- %s\n\t- %s\n\n",
		os.Args[0], Version, AssetCountryName, AssetASNname)
	fmt.Printf("Usage: %s ip or stdin input\n", os.Args[0])
	os.Exit(0)
}
func main() {

	// embedded values are build from dynembedded.go
	ASN := initGeo(embeddedASN)
	Country := initGeo(embeddedCountry)

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

// initGeo Convert Asset to Geoip2.Reader
func initGeo(Asset []byte) *geoip2.Reader {

	db, err := geoip2.FromBytes(Asset)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// parseandprint return string with ASN and Country
func parseandprint(ips string, dbCountry, dbASN *geoip2.Reader) string {
	var record *geoip2.Country

	ip := net.ParseIP(ips)
	if ip == nil {
		log.Printf("Unable to parse ip : \"%s\"", ips)
		return ""
	}

	var err error

	record, err = dbCountry.Country(ip)
	if err != nil || record.Country.Names["en"] == "" {
		log.Printf("Unable to geoloc Country of \"%s\"", ips)
		return ""
	}

	ASN, err := dbASN.ASN(ip)
	if err != nil {
		log.Printf("ERR db.ASN \"%s\"", ips)
	}
	output := fmt.Sprintf("%s (%s / AS%d %s)", ips,
		record.Country.Names["en"],
		ASN.AutonomousSystemNumber, ASN.AutonomousSystemOrganization)

	return output
}

// readandprintbulk read stdin and printout results
// channel is used when parseandprint use ip to name resolution
// in that way, we can use
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

		go func(line string, mywg *sync.WaitGroup, mychan chan bool) {
			out := parseandprint(line, dbCountry, dbASN)
			if len(out) > 0 {
				fmt.Printf("%s\n", out)
			}
			<-mychan
			mywg.Done()
		}(line, &wg, limitChan)

	}
	wg.Wait()
}
