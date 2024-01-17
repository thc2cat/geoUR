package main

// geoUR : read ips from stdin, display for Country and ASN

import (
	"fmt"
	"os"
)

// History :
// v1.0 combine geoip2 Country and ASN,
// v1.2 testing functions as parameter
// v1.3 Better json output
// v1.4 correcting for & ' and \ caracters
//
// it doesn't resolve ip anymore
//	use "| adnsresfilter -ua" if resolution is needed
//

var (
	maxrequests = 128
	// Version given by git tag via Makefile
	Version = "v1.4"
)

// printVersionUsage only print Version and Usage anq exit
func printVersionUsage() {
	fmt.Printf("geoUR (%s) build with MaxMind free edition db :\n\t- %s\n\t- %s\n\n",
		Version, AssetCountryName, AssetASNname)
	fmt.Printf("Usage: geoUR ip or stdin input\nuse BULKFORMAT env for JSON, SED format")
	os.Exit(0)
}

func main() {

	bulkformat := getenv("BULKFORMAT", "")
	var printfunc FuncStringtoString

	switch bulkformat {
	case "JSON":
		printfunc = parseandprintJSON
	case "SED":
		printfunc = parseandprintSED
	default:
		printfunc = parseandprint
	}

	switch {
	case len(os.Args) == 1:
		// Read stdin
		readandprintbulk(printfunc)

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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
