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
func Str2Int32(x string) int32 {
	res, _ := strconv.ParseInt(x, 10, 32)
	return int32(res)
}
func Str2Int64(x string) int64 {
	res, _ := strconv.ParseInt(x, 10, 64)
	return res
}
