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
		{"localhost", "127.0.0.1", "127.0.0.1 (ERR geoip2.Country / AS0ERR)"},
		{"broadcast", "255.255.255.0", "255.255.255.0 (ERR geoip2.Country / AS0ERR)"},
		{"network", "193.51.24.0/24", "193.51.24.0/24 (ERR net.ParseIP / AS0ERR)"},
		{"not a ipv4", "289.0.0.1", "289.0.0.1 (ERR net.ParseIP / AS0ERR)"},
		{"not a ipv4", "2001::0:0:1", "2001::0:0:1 (ERR geoip2.Country / AS0ERR)"},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseandprint(tt.ips); got != tt.want {
				t.Errorf("parseandprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
