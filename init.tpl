package main

import _ "embed"

//go:embed data/GeoLite2-ASN.mmdb
var embeddedASN []byte

//go:embed data/GeoLite2-Country.mmdb
var embeddedCountry []byte

var AssetCountryName = "%COUNTRYFILE%"
var AssetASNname = "%ASNFILE%"
