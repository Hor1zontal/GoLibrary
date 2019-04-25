/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/27
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package cache

import (
	"github.com/garyburd/redigo/redis"
	"aliens/log"
	"time"
)

//操作类型
const (
	PARAM_WITHSCORES string = "WITHSCORES"
	PARAM_LIMIT      string = "limit"
	PARAM_Z_MAX      string = "+inf"
	PARAM_Z_MIN      string = "-inf"

	OP_SELECT string = "SELECT"

	OP_S_ADD        string = "SADD"
	OP_S_RANDMENBER string = "SRANDMEMBER"
	OP_S_ISMEMBER   string = "SISMEMBER"
	OP_S_MEMBERS	string = "SMEMBERS"
	OP_S_DEL		string = "SREM"

	OP_DUMP    string = "DUMP"
	OP_RESTORE string = "RESTORE"

	OP_SET    string = "SET"
	OP_GET    string = "GET"
	OP_DEL    string = "DEL"
	OP_EXISTS string = "EXISTS"
	OP_SETEX  string = "SETEX"
	OP_SETNX  string = "SETNX"
	OP_INCR   string = "INCR"   //自增长
	OP_DECR   string = "DECR"   //自减
	OP_INCRBY string = "INCRBY"

	OP_EXPIRE string = "EXPIRE"

	OP_FLUSHALL string = "FLUSHALL"

	OP_H_SET     string = "HSET"
	OP_H_GET     string = "HGET"
	OP_H_GETALL  string = "HGETALL"
	OP_H_MGET    string = "HMGET"
	OP_H_MSET    string = "HMSET"
	OP_H_DEL     string = "HDEL"
	OP_H_HINCRBY string = "HINCRBY"
	OP_H_EXISTS  string = "HEXISTS"
	OP_H_LEN     string = "HLEN"

	OP_L_PUSH  string = "LPUSH"
	OP_R_PUSH  string = "RPUSH"
	OP_L_RANGE string = "LRANGE"
	OP_L_LEN   string = "LLEN"

	OP_Z_ADD             string = "ZADD"
	OP_Z_COUNT           string = "ZCOUNT"
	OP_Z_REM             string = "ZREM"
	OP_Z_RANGE           string = "ZRANGE"
	OP_Z_RANGEBYSCORE    string = "ZRANGEBYSCORE"
	OP_Z_REVRANGEBYSCORE string = "ZREVRANGEBYSCORE"
	OP_Z_REVRANGE        string = "ZREVRANGE"

	OP_Z_REVRANK string = "ZREVRANK"
	OP_Z_SCORE string = "ZSCORE"

	OP_PUBLISH string = "PUBLISH"
)

type RedisCacheClient struct {
	MaxIdle     int
	MaxActive   int
	Address     string
	Password    string
	IdleTimeout time.Duration //180 * time.Second
	pool        *redis.Pool
}

//redis.pool.maxActive=200  #最大连接数：能够同时建立的“最大链接个数”

//redis.pool.maxIdle=20     #最大空闲数：空闲链接数大于maxIdle时，将进行回收

//redis.pool.minIdle=5      #最小空闲数：低于minIdle时，将创建新的链接

//redis.pool.maxWait=3000    #最大等待时间：单位ms

//启动缓存客户端
func (this *RedisCacheClient) Start() {
	this.pool = &redis.Pool{
		MaxIdle:     this.MaxIdle,
		MaxActive:   this.MaxActive,
		IdleTimeout: this.IdleTimeout, //空闲释放时间
		Dial: func() (redis.Conn, error) {
			if this.Address == "" {
				this.Address = "127.0.0.1:6379"
			}
			c, err := redis.Dial("tcp", this.Address)
			if err != nil {
				log.Panic("start redis error : %v", err)
				return nil, err
			}
			if this.Password != "" {
				if _, err := c.Do("AUTH", this.Password); err != nil {
					c.Close()
					log.Panic("start redis error : %v", err)
					return nil, err
				}
			}
			return c, err
		},
	}
	//测试连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SET, "TestKey", "TestData")
	if err != nil {
		log.Panic("test redis error %v", err)
	}
}

//关闭缓存客户端
func (this *RedisCacheClient) Close() {
	if this.pool != nil {
		this.pool.Close()
	}
}

func (this *RedisCacheClient) SetNX(key string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	result, _ := redis.Int(conn.Do(OP_SETNX, key, value))
	return result == 1

}

//设置数据过期时间
func (this *RedisCacheClient) Expire(key string, seconds int) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_EXPIRE, key, seconds)
	if err != nil {
		return false
	}
	return true
}

//添加数据
func (this *RedisCacheClient) SetData(key string, value interface{}) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SET, key, value)
	if err != nil {
		log.Debug("%v", err)
		return false
	}
	return true
}

//添加数据
func (this *RedisCacheClient) SetExpireData(key string, value interface{}, seconds int) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SET, key, value, "EX", seconds)
	if err != nil {
		log.Debug("%v", err)
		return false
	}
	return true
}

func (this *RedisCacheClient) Incr(key string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_INCR, key))
}

func (this *RedisCacheClient) IncrBy(key string, add interface{}) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_INCRBY, key, add))
}

func (this *RedisCacheClient) Decr(key string) (int, error) {
	conn := this.pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do(OP_DECR, key))
}

func (this *RedisCacheClient) SelectDB(dbNumber int) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_SELECT, dbNumber)
	if err != nil {
		//log.Debug("%v", err)
		return false
	}
	return true
}

func (this *RedisCacheClient) GetDataInt32(key string) int {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Int(conn.Do(OP_GET, key))
	if err != nil {
		//log.Debug("%v", err)
		return 0
	}
	return value
}

func (this *RedisCacheClient) GetDataInt64(key string) int64 {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Int64(conn.Do(OP_GET, key))
	if err != nil {
		//log.Debug("%v", err)
		return 0
	}
	return value
}

//获取数据
func (this *RedisCacheClient) GetData(key string) string {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do(OP_GET, key))
	if err != nil {
		//log.Debug("%v", err)
		return ""
	}
	return value
}

func (this *RedisCacheClient) GetBytesData(key string) []byte {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.Bytes(conn.Do(OP_GET, key))
	if err != nil {
		//log.Debug("%v", err)
		return nil
	}
	return value
}

//导出数据
func (this *RedisCacheClient) Dump(key string) string {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do(OP_DUMP, key))
	if err != nil {
		return ""
	}
	return value
}

//导入数据
func (this *RedisCacheClient) Restore(key string, data string) string {
	conn := this.pool.Get()
	defer conn.Close()
	value, err := redis.String(conn.Do(OP_RESTORE, key, 0, data))
	if err != nil {
		return ""
	}
	return value
}

//是否存在数据
func (this *RedisCacheClient) Exists(key string) (bool, error) {
	conn := this.pool.Get()
	defer conn.Close()
	result, err := redis.Bool(conn.Do(OP_EXISTS, key))
	if err != nil {
		//log.Debug("%v", err)
		return false, err
	}
	return result, err
}

//添加数据
func (this *RedisCacheClient) DelData(key string) bool {
	conn := this.pool.Get()
	defer conn.Close()
	_, err := conn.Do(OP_DEL, key)
	if err != nil {
		//log.Debug("%v", err)
		return false
	}
	return true
}

//清除所有数据
func (this *RedisCacheClient) FlashAll() {
	conn := this.pool.Get()
	defer conn.Close()
	conn.Do(OP_FLUSHALL)
}

//// 存map
//func (this *RedisCacheClient)SetMap(key string ,value map[string]string) bool{
//	conn := this.pool.Get()
//	defer conn.Close()
//	// 转换成json
//	v,_ := json.Marshal(value)
//	// 存redis
//	_,err := conn.Do("SETNX",key, v)
//	if err != nil {
//		//log.Debug("%v",err)
//		return false
//	}
//	return true
//}
//
//// 取map
//func (this *RedisCacheClient)GetMap(key string) map[string]string {
//	conn := this.pool.Get()
//	defer conn.Close()
//	var imap map[string]string
//	value,err := redis.Bytes(conn.Do("GET",key))
//	if err != nil {
//		//log.Debug("%v",err)
//		return nil
//	}
//	// json转map
//	errShal := json.Unmarshal(value,&imap)
//	if errShal != nil {
//		//log.Debug("%v",errShal)
//		return nil
//	}
//	return imap
//}

//订阅数据变更
func (this *RedisCacheClient) Subscribe(callback func(channel, value string), channel ...interface{}) {
	//defer conn.Close()
	psc := redis.PubSubConn{Conn: this.pool.Get()}
	psc.Subscribe(channel...)
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				value, _ := redis.String(v.Data, nil)
				callback(v.Channel, value)
			case error:
				//log.Debug("error: %v\n", v)
				return
			}
		}
	}()
}

func (this *RedisCacheClient) PSubscribe(callback func(pattern, channel, value string), channel ...interface{}) {
	//defer conn.Close()
	psc := &redis.PubSubConn{Conn: this.pool.Get()}
	psc.PSubscribe(channel...)
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				value, _ := redis.String(v.Data, nil)
				callback(v.Pattern, v.Channel, value)
			case error:
				//log.Debug("error: %v\n", v)
				return
			}
		}
	}()
}

//发布数据
func (this *RedisCacheClient) Publish(channel, value interface{}) {
	conn := this.pool.Get()
	defer conn.Close()
	conn.Do(OP_PUBLISH, channel, value)
}
