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
	"github.com/garyburd/redigo/redis"
	//"github.com/name5566/leaf/log"
)

//判断set是否包含成员
func (this *RedisCacheClient) SContains(key string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	result, err := redis.Bool(conn.Do(OP_S_ISMEMBER, key, value))
	if err != nil {
		//log.Debug("%v", err)
		return false
	}
	return result
}

//Set添加数据
func (this *RedisCacheClient) SAddData(key string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_S_ADD, key, value)
	if err != nil {
		//log.Debug("%v", err)
		return false
	}
	return true
}

//随机Set中指定数量的数据   repeat:是否重复
func (this *RedisCacheClient) SRandMember(key string, value int, repeat bool) []int {
	conn := this.pool.Get()
	defer conn.Close()
	if repeat {
		value = -value
	}

	result, err := redis.Ints(conn.Do(OP_S_RANDMENBER, key, value))
	if err != nil {
		//log.Debug("%v",err)
	}
	return result
}

func (this *RedisCacheClient) SMembers(key string) []int {
	conn := this.pool.Get()
	defer conn.Close()

	result, _ := redis.Ints(conn.Do(OP_S_MEMBERS, key))
	return result
}

func (this *RedisCacheClient) SDelData(key string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_S_DEL, key, value)
	if err != nil {
		//log.Debug("%v", err)
		return false
	}
	return true
}