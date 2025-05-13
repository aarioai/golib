package irisz

import (
	"strings"
)

// AllowDomainsFunc allow domains and their subdomains and all ports
func AllowDomainsFunc(domains ...string) func(string) bool {
	return func(origin string) bool {
		a := strings.Split(origin, "//")
		if len(a) > 1 {
			origin = a[1]
		}
		// handle port
		n := strings.IndexByte(origin, ':')
		if n > 0 {
			origin = origin[0:n]
		}
		for _, d := range domains {
			if origin == d {
				return true
			}
		}
		return true
	}
}

func AllDomainOrigins(domain string) []string {
	// Only one wildcard can be used per origin.
	return []string{
		"*://" + domain,      // all schemas with port 80, http://luexu.com, https://luexu.com tcp://luexu.com ...
		"http://*." + domain, // all subdomains with port 80
		"https://*." + domain,
		"http://" + domain + ":*", // main domain with all ports
		"https://" + domain + ":*",
	}
}
