package main

// geoUR : read ips from stdin, display for Country and ASN

import (
	"fmt"
	"os"
)

// History :
// v1.0 combine geoip2 Country and ASN,
// v1.2 testing functions as parameter
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

	switch {
	case len(os.Args) == 1:
		// Readstdin
		readandprintbulk(parseandprint)

	case (len(os.Args) == 2) && (os.Args[1] == "-h" || os.Args[1] == "-v"):
		printVersionUsage()

	case len(os.Args) == 2:
		// Single input as parameter
		fmt.Printf("%s\n", parseandprint(os.Args[1]))

	default:
		printVersionUsage()
	}

	os.Exit(0)
}
