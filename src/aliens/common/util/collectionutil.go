/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/4/19
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import (
	"encoding/json"
	"math/rand"
)

func MargeMap(a map[int32]int32, b map[int32]int32) {
	for key, value := range b {
		a[key] = value
	}
}


//从a中过滤掉b存在的数据
func FilterArray(a []int32, b []int32) []int32 {
	if b == nil || len(b) == 0 {
		return a
	}
	result := make([]int32, 0)
	for _, value := range a {
		if !ContainsInt32(value, b) {
			result = append(result, value)
		}
	}
	return result
}

func CopyMap(a map[int32]int32) map[int32]int32 {
	results := make(map[int32]int32)
	for key, value := range a {
		results[key] = value
	}
	return results
}

func JSONCopy(marshaler interface{}, unMarshaler interface{}) error {
	data, error := json.Marshal(marshaler)
	if error != nil {
		return error
	}
	return json.Unmarshal(data, unMarshaler)
}

func RandomMultiWeight(weightMapping map[int32]int32, count int) []int32 {
	results := []int32{}
	for i := 0; i < count; i++ {
		result := RandomWeight(weightMapping)
		if result == 0 {
			return results
		}
		results = append(results, result)
		delete(weightMapping, result)
	}
	return results
}

func RandIntervalN(b1, b2 int32, n int) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = int(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := 0; i < n; i++ {
		v := int32(rand.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}


func RandomWeight(weightMapping map[int32]int32) int32 {
	var totalWeight int32 = 0
	for _, weight := range weightMapping {
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0
	}
	randomValue := rand.Int31n(totalWeight) + 1
	var currentValue int32 = 0
	for id, weight := range weightMapping {
		currentValue += weight
		if currentValue >= randomValue {
			return id
		}
	}
	return 0
}

type WeightData interface {
	GetWeight() int32
}

func RandomWeightData(weightMapping map[int32]WeightData) WeightData {
	var totalWeight int32 = 0
	for _, weightData := range weightMapping {
		weight := weightData.GetWeight()
		if weight <= 0 {
			continue
		}
		totalWeight += weight
	}
	if (totalWeight <= 0) {
		return nil
	}
	randomValue := rand.Int31n(totalWeight) + 1
	var currentValue int32 = 0
	for _, weightData := range weightMapping {
		weight := weightData.GetWeight()
		if weight <= 0 {
			continue
		}
		currentValue += weight
		if currentValue >= randomValue {
			return weightData
		}
	}
	return nil
}



func RandomFloat32Weight(weightMapping map[float32]int32) float32 {
	var totalWeight int32 = 0
	for _, weight := range weightMapping {
		totalWeight += weight
	}
	if totalWeight <= 0 {
		return 0
	}
	randomValue := rand.Int31n(totalWeight) + 1
	var currentValue int32 = 0
	for id, weight := range weightMapping {
		currentValue += weight
		if currentValue >= randomValue {
			return id
		}
	}
	return 0
}

