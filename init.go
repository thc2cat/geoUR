package main

import _ "embed"

//go:embed data/GeoLite2-ASN.mmdb
var embeddedASN []byte

//go:embed data/GeoLite2-Country.mmdb
var embeddedCountry []byte

var AssetCountryName = "GeoLite2-Country_20231215"
var AssetASNname = "GeoLite2-ASN_20231215"
