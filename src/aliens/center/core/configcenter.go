/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 *
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import (
	"aliens/log"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

const CONFIG_NODE_NAME string = "config"
const NODE_SPLIT string = "/"

type ConfigCenter struct {
	zkCon      *zk.Conn
	zkName     string
	configRoot string
	listener   *ConfigListener
}

//启动服务中心客户端
func (this *ConfigCenter) Start(address string, timeout int, zkName string) {
	this.StartCluster([]string{address}, timeout, zkName)
}

func (this *ConfigCenter) StartCluster(addressArray []string, timeout int, zkName string) {
	this.zkName = zkName
	c, _, err := zk.Connect(addressArray, time.Duration(timeout)*time.Second)
	if err != nil {
		panic(err)
	}
	this.configRoot = NODE_SPLIT + this.zkName + NODE_SPLIT + CONFIG_NODE_NAME
	this.zkCon = c
	this.confirmNode(NODE_SPLIT + this.zkName)
	this.confirmNode(this.configRoot)
	this.listener = NewConfigListener()
}

func (this *ConfigCenter) assert() {
	if this.zkCon == nil {
		panic("mast start config center first")
	}
}

//关闭服务中心
func (this *ConfigCenter) Close() {
	if this.zkCon != nil {
		this.zkCon.Close()
	}
}


//订阅服务  能实时更新服务信息
//func (this *ConfigCenter) SubscribeConfigs(configTypes ...string) {
//	this.assert()
//	//this.listener = listener
//	for _, configType := range configTypes {
//		this.subscribeConfig(configType)
//	}
//}


//func (this *ConfigCenter) SubscribeConfig(configType string, handler func(data []byte)) {
//	this.assert()
//
//	this.subscribeConfig(configType)
//}

func (this *ConfigCenter) SubscribeConfig(configType string, handler func(data []byte)) {
	this.assert()
	path := this.configRoot + NODE_SPLIT + configType
	this.confirmNode(path)

	content, _, ch, err := this.zkCon.GetW(path)
	if !this.listener.Contains(configType) {
		go func(){
			for {
				event, _ := <-ch
				//更新配置节点信息
				if event.Type == zk.EventNodeDataChanged {
					//content, _, err := this.zkCon.Get(path)
					content, _, chw, err := this.zkCon.GetW(path)
					if err == nil {
						log.Debug("config %v update", configType)
						this.listener.ConfigChange(configType, content)
					}
					ch = chw
				}
			}

		}()
	}
	this.listener.AddHandler(configType, handler)
	if this.listener != nil {
		this.listener.ConfigChange(configType, content)
	}
	if err != nil {
		log.Info("subscribe config %v error: %v", path, err)
		return
	}
	//go this.openListener(configType, path, ch)
}

//监听配置变更
//func (this *ConfigCenter) openListener(configType string, path string, ch <-chan zk.Event) {
//	event, _ := <-ch
//	//更新配置节点信息
//	if event.Type == zk.EventNodeDataChanged {
//		this.subscribeConfig(configType)
//	}
//}

//
func (this *ConfigCenter) confirmNode(path string, flags ...int32) bool {
	_, err := this.zkCon.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	return err == nil
}

//发布配置信息
func (this *ConfigCenter) PublicConfig(configType string, configContent []byte) bool {
	this.assert()
	if configType == "" {
		log.Info("config type con not be empty")
		return false
	}

	configPath := this.configRoot + NODE_SPLIT + configType
	this.confirmNode(configPath)
	_, err := this.zkCon.Set(configPath, configContent, -1)
	if err != nil {
		log.Info("public config %v  err : %v", configType, err)
		return false
	}
	log.Info("public config %v success", configType)
	return true
}
