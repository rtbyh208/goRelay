package pkg

import "strings"

func IsWhitelisted(ipaddr string, whitelist []string) bool {

	if len(whitelist) == 0 {
		return true
	}

	ip := strings.Split(ipaddr, ":")[0]

	for _, v := range whitelist {
		if v == ip {
			return true
		}
	}

	return false
}
