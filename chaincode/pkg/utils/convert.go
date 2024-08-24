package utils

import (
	"fmt"
	"strconv"
)

func Int2Byte(x int32) []byte {
	return []byte(fmt.Sprintf("%d", x))
}

func Byte2Int(x []byte) int32 {
	var res int32
	fmt.Sscanf(string(x), "%d", &res)
	return res
}

func Int2Str(x int32) string {
	return strconv.FormatInt(int64(x), 10)
}
