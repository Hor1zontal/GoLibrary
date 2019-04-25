package util

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
	"strings"
	"unsafe"
)

var RANDOM_ARRAY []byte = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

func GetRandomString(l int) string {
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, RANDOM_ARRAY[r.Intn(len(RANDOM_ARRAY))])
	}
	return string(result)
}

//随机验证码
func RandomVerifyCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

func RandInt32Scop(min int32, max int32) int32 {
	diff := max - min + 1
	return min + rand.Int31n(diff)
}

func Round(value float64) int32 {
	return int32(value + 0.5)
}

func Round64(value float64) int64 {
	return int64(value + 0.5)
}

func RoundFloat64(value float64) float64 {
	return float64(int64(value + 0.5))
}

func StringArray2Int32Array(array []string) []int32 {
	var result []int32
	for _, value := range array {
		new_value, _ := strconv.Atoi(value)
		result = append(result, int32(new_value))
	}
	return result
}

func ContainsInt32(value int32, array []int32) bool {
	if len(array) == 0 {
		return false
	}
	for _, member := range array {
		if member == value {
			return true
		}
	}
	return false
}

func StringToInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}

func StringToInt32(value string) int32 {
	result, _ := strconv.Atoi(value)
	return int32(result)
}

func StringToInt64(value string) int64 {
	result, _ := strconv.ParseInt(value, 10, 64)
	return result
}

func StringToFloat32(value string) float32 {
	result, _ := strconv.ParseFloat(value, 64)
	return float32(result)
}

func StringToFloat64(value string) float64 {
	result, _ := strconv.ParseFloat(value, 64)
	return result
}

func Int32ToString(value int32) string {
	return strconv.Itoa(int(value))
}


func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Int32Max(a int32, b int32) int32 {
	return int32(math.Max(float64(a), float64(b)))
}

func Int64Max(a int64, b int64) int64 {
	return int64(math.Max(float64(a), float64(b)))
}

func Int32Min(a int32, b int32) int32 {
	return int32(math.Min(float64(a), float64(b)))
}

func Int64Min(a int64, b int64) int64 {
	return int64(math.Min(float64(a), float64(b)))
}

func FirstToUpper(str string) string {
	length := len(str)
	return strings.ToUpper(str[0:1]) + str[1:length]
	//temp := strings.Split(str, "_")
	//var upperStr string
	//for y := 0; y < len(temp); y++ {
	//	vv := []rune(temp[y])
	//	if y != 0 {
	//		for i := 0; i < len(vv); i++ {
	//			if i == 0 {
	//				vv[i] -= 32
	//				upperStr += string(vv[i]) // + string(vv[i+1])
	//			} else {
	//				upperStr += string(vv[i])
	//			}
	//		}
	//	}
	//}
	//return temp[0] + upperStr
}

func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//类型转换  bytes to string
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//func ToStr(value interface{}, args ...int) (s string) {
//	switch v := value.(type) {
//	case bool:
//		s = strconv.FormatBool(v)
//	case float32:
//		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
//	case float64:
//		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
//	case int:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int8:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int16:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int32:
//		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
//	case int64:
//		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
//	case uint:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint8:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint16:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint32:
//		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
//	case uint64:
//		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
//	case string:
//		s = v
//	case []byte:
//		s = string(v)
//	default:
//		s = fmt.Sprintf("%v", v)
//	}
//	return s
//}