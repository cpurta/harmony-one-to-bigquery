package util

import (
	"strconv"
	"strings"
)

func HexToInt(hexString string) (int64, error) {
	str := strings.Replace(hexString, "0x", "", -1)
	str = strings.Replace(str, "0X", "", -1)

	return strconv.ParseInt(str, 16, 64)
}
