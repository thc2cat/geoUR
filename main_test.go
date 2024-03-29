package main

import (
	"testing"
)

func Test_parseandprint(t *testing.T) {

	tests := []struct {
		name string
		ips  string
		want string
	}{
		{"test1", "193.51.24.1", "193.51.24.1 (France / AS2200 Renater)"},
		{"localhost", "127.0.0.1", "127.0.0.1 (ERROR / AS0 Error db.Country)"},
		{"broadcast", "255.255.255.0", "255.255.255.0 (ERROR / AS0 Error db.Country)"},
		{"network", "193.51.24.0/24", "193.51.24.0/24 (ERROR / AS0 Error net.parseIP)"},
		{"not a ipv4", "289.0.0.1", "289.0.0.1 (ERROR / AS0 Error net.parseIP)"},
		{"not a ipv4", "2001::0:0:1", "2001::0:0:1 (ERROR / AS0 Error db.Country)"},
		{"DEL '", "193.56.242.17", "193.56.242.17 (France / AS25117 Euro-Information-Europeenne de Traitement de l Information SAS)"},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseandprint(tt.ips); got != tt.want {
				t.Errorf("parseandprint() \nreturn \t%v\nwant \t%v", got, tt.want)
			}
		})
	}
}
