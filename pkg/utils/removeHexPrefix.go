package utils

import "strings"

func RemoveHexPrefix(target string) string {
	if strings.HasPrefix(target, "0x") {
		return target[2:]
	}
	return target
}
