/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/24
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

type IService interface {
	GetDesc() string //获取服务的描述信息
	GetID() string   //获取id
	SetID(id string) //设置id
	GetType() string
	SetType(serviceType string)
	Start() bool                              //启动服务
	Connect() bool                            //连接服务
	Close()									  //关闭服务
	Equals(other IService) bool               //比较服务
	IsLocal() bool                            //是否本机服务
	Request(interface{}) (interface{}, error) //服务请求
	AsyncRequest(interface{}) error  //异步请求
}

type IServiceFactory interface {
	CreateService(data []byte) IService //从数据中加载服务信息
}
