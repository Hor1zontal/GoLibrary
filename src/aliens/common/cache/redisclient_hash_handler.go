/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"aliens/log"
	"reflect"
	"time"
)

//自增指定数量
func (this *RedisCacheClient) HIncrby(key interface{}, field string, increment int) int {
	conn := this.pool.Get()
	defer conn.Close()
	result, err := redis.Int(conn.Do(OP_H_HINCRBY, key, field, increment))
	if err != nil {
		//log.Debug("%v",	err)
		return 0
	}
	return result
}

func (this *RedisCacheClient) HSet(key interface{}, field string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_H_SET, key, field, value)
	if err != nil {
		//log.Debug("%v",	err)
		return false
	}
	return true
}

func (this *RedisCacheClient) HGetBytes(key interface{}, field string) []byte {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Bytes(conn.Do(OP_H_GET, key, field))
	if err != nil {
		return nil
	}
	return value
}

func (this *RedisCacheClient) HGet(key interface{}, field string) string {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do(OP_H_GET, key, field))
	if err != nil {
		//log.Debug("%v",err)
		return ""
	}
	return value
}


func (this *RedisCacheClient) HLen(key interface{}) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_H_LEN, key))
}


func (this *RedisCacheClient) HDel(key interface{}, field string) error {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_H_DEL, key, field)
	return err
}

//提取结构体的注解，写入redis
func (this *RedisCacheClient) HSetFieldData(key interface{}, fieldPrefix string, data interface{}) {
	conn := this.pool.Get()
	defer conn.Close()
	dataValue := reflect.ValueOf(data).Elem()
	dataType := reflect.TypeOf(data).Elem()
	for i := 0; i < dataValue.NumField(); i++ {
		fieldValue := dataValue.Field(i)
		fieldType := dataType.Field(i)
		tag := fieldType.Tag.Get("rorm")
		if tag == "" {
			continue
		}
		tag = fieldPrefix + tag
		if fieldValue.Kind() == reflect.Struct {
			if timeValue, ok := fieldValue.Interface().(time.Time); ok {
				conn.Do(OP_H_SET, key, tag, timeValue.Unix())
			}
		} else {
			conn.Do(OP_H_SET, key, tag, fieldValue.Interface())
		}
	}
}

//获取redis数据，注入结构体
func (this *RedisCacheClient) HGetFieldData(key interface{}, fieldPrefix string, data interface{}) {
	conn := this.pool.Get()
	defer conn.Close()
	dataValue := reflect.ValueOf(data).Elem()
	dataType := reflect.TypeOf(data).Elem()
	for i := 0; i < dataValue.NumField(); i++ {
		fieldValue := dataValue.Field(i)
		field := dataType.Field(i)
		tag := field.Tag.Get("rorm")
		if tag == "" {
			continue
		}
		tag = fieldPrefix + tag
		switch fieldValue.Kind() {
		case reflect.String:
			value, err := redis.String(conn.Do(OP_H_GET, key, tag))
			if err == nil {
				fieldValue.SetString(value)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, err := redis.Int64(conn.Do(OP_H_GET, key, tag))
			if err == nil {
				fieldValue.SetInt(value)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value, err := redis.Uint64(conn.Do(OP_H_GET, key, tag))
			if err == nil {
				fieldValue.SetUint(value)
			}
		case reflect.Bool:
			value, err := redis.Bool(conn.Do(OP_H_GET, key, tag))
			if err == nil {
				fieldValue.SetBool(value)
			}
		case reflect.Float32, reflect.Float64:
			value, err := redis.Float64(conn.Do(OP_H_GET, key, tag))
			if err == nil {
				fieldValue.SetFloat(value)
			}
		case reflect.Struct:
			if _, ok := fieldValue.Interface().(time.Time); ok {
				value, err := redis.Int64(conn.Do(OP_H_GET, key, tag))
				if err == nil {
					fieldValue.Set(reflect.ValueOf(time.Unix(value, 0)))
				}
			}
		default:
			log.Debug("unsupport redis type: %v", fieldValue.Kind())
		}

	}
}

func (this *RedisCacheClient) HSetData(key interface{}, data interface{}) {
	this.HSetFieldData(key, "", data)
}

//获取redis数据，注入结构体
func (this *RedisCacheClient) HGetData(key interface{}, data interface{}) {
	this.HGetFieldData(key, "", data)
}

func (this *RedisCacheClient) HGetBool(key interface{}, field string) bool {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Bool(conn.Do(OP_H_GET, key, field))
	if err != nil {
		////log.Debug("%v",err)
		return false
	}
	return value
}

func (this *RedisCacheClient) HGetInt32(key interface{}, field string) int32 {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do(OP_H_GET, key, field))
	if err != nil {
		//log.Debug("%v",err)
		return 0
	}
	return int32(value)
}

func (this *RedisCacheClient) HGetInt64(key interface{}, field string) int64 {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Int64(conn.Do(OP_H_GET, key, field))
	if err != nil {
		//log.Debug("%v",err)
		return 0
	}
	return value
}

func (this *RedisCacheClient) HGetFloat64(key interface{}, field string) float64 {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Float64(conn.Do(OP_H_GET, key, field))
	if err != nil {
		//log.Debug("%v",err)
		return 0
	}
	return value
}

//判断hash字段是否存在
func (this *RedisCacheClient) HFieldExists(key interface{}, field string) bool {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Bool(conn.Do(OP_H_EXISTS, key, field))
	if err != nil {
		//log.Debug("%v",err)
		return false
	}
	return value
}

//获取所有的hash数据
func (this *RedisCacheClient) HGetAllInt(key interface{}) map[string]int {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.IntMap(conn.Do(OP_H_GETALL, key))
	if err != nil {
		//log.Debug("%v",err)
		return nil
	}
	return value
}

func (this *RedisCacheClient) HGetAllInt64(key interface{}) map[string]int64 {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Int64Map(conn.Do(OP_H_GETALL, key))
	if err != nil {
		//log.Debug("%v",err)
		return nil
	}
	return value
}

func (this *RedisCacheClient) HGetAll(key interface{}) map[string]string {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.StringMap(conn.Do(OP_H_GETALL, key))
	if err != nil {
		//log.Debug("%v",err)
		return nil
	}
	return value
}

//批量添加字段
func (this *RedisCacheClient) HMSet(key interface{}, fields map[interface{}]interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()

	params := []interface{}{key}

	for key, value := range fields {
		params = append(params, key, value)
	}

	_, err := conn.Do(OP_H_MSET, params...)
	if err != nil {
		//log.Debug("%v",err)
		return false
	}
	return true
}

func (this *RedisCacheClient) HMGet(key interface{}, fieldNames ...interface{}) map[string]string {
	conn := this.pool.Get()
	defer conn.Close()

	params := []interface{}{key}
	params = append(params, fieldNames...)

	values, err := redis.Strings(conn.Do(OP_H_MGET, params...))
	if err != nil {
		//log.Debug("%v",err)
		return nil
	}
	if len(values) != len(fieldNames) {
		return nil
	}
	results := make(map[string]string)
	for index, value := range values {
		results[fieldNames[index].(string)] = value
	}
	return results
}

func (this *RedisCacheClient) HMarshal(key interface{}, object interface{}) error {
	conn := this.pool.Get()
	defer conn.Close()
	t := reflect.ValueOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("object must be struct")
	}

	params := []interface{}{key}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Type().Name()
		value := field.Interface()
		params = append(params, key, value)
	}
	_, err := conn.Do(OP_H_MSET, params...)
	return err
}
