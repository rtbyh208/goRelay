package pkg

import "strings"

func IsBlacklisted(ipaddr string, blacklist []string) bool {

	if len(blacklist) == 0 {
		return false
	}

	ip := strings.Split(ipaddr, ":")[0]

	for _, v := range blacklist {
		if v == ip {
			return true
		}
	}

	return false
}
