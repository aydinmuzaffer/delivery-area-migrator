package envutils

import (
	"os"
	"testing"
)

func TestGetENVWhenPassed(t *testing.T) {
	var tests = []struct {
		envName         string
		envValue        string
		envDefaultValue string
		want            string
	}{
		{"PricingDbAddres", "This is PricingDbAddres", "This is PricingDbAddresDefault", "This is PricingDbAddres"},
		{"VendorDbAddres", "This is VendorDbAddres", "This is VendorDbAddresDefault", "This is VendorDbAddres"},
	}

	for _, tt := range tests {
		// t.Run enables running "subtests", one for each
		// table entry. These are shown separately
		// when executing `go test -v`.
		os.Setenv(tt.envName, tt.envValue)
		t.Run(tt.envName, func(t *testing.T) {
			got := GetENV(tt.envName, tt.envDefaultValue)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestGetENVWhenNotPassed(t *testing.T) {
	var tests = []struct {
		envName         string
		envDefaultValue string
		want            string
	}{
		{"PricingDbAddres", "This is PricingDbAddres", "This is PricingDbAddres"},
		{"VendorDbAddres", "This is VendorDbAddres", "This is VendorDbAddres"},
	}

	for _, tt := range tests {
		t.Run(tt.envName, func(t *testing.T) {
			got := GetENV(tt.envName, tt.envDefaultValue)
			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
