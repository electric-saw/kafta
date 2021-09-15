package kafka

import "strings"

func removePort(address string) string {
	if i := strings.Index(address, ":"); i > 0 {
		return address[:i]
	}
	return address
}
