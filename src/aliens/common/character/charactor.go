package character

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
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

func CopyMap(source map[int32]float32) map[int32]float32 {
	target := make(map[int32]float32)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyFloat64Map(source map[int32]float64) map[int32]float64 {
	target := make(map[int32]float64)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyInt32Map(source map[int32]int32) map[int32]int32 {
	target := make(map[int32]int32)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyInt64Map(source map[int32]int64) map[int32]int64 {
	target := make(map[int32]int64)
	for key, value := range source {
		target[key] = value
	}
	return target
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

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

//随机验证码
func RandomVerifyCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
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

func LeftString(str string, index int) string {
	if index <= 0 || index >= len(str) {
		return ""
	}
	strLeft := string([]byte(str)[:index])
	return strLeft
}

func RightString(str string, index int) string {
	if index <= 0 || index >= len(str) {
		return ""
	}
	strLeft := string([]byte(str)[len(str)-3:])
	return strLeft
}

func AppendArray(array []int32, appendArray []int32) []int32 {
	for _, member := range appendArray {
		array = append(array, member)
	}
	return array
}

func CopyArray(src []int32) []int32 {
	result := make([]int32, len(src))
	for index, value := range src {
		result[index] = value
	}
	return result
}

//获取数组1在数组2中不存在的元素
func GetArrayDeff(array1 []int32, array2 []int32) []int32 {
	result := []int32{}
	for _, member := range array1 {
		if !ContainsInt32(member, array2) {
			result = append(result, member)
		}
	}
	return result
}

