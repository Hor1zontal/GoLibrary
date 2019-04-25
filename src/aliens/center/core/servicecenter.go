/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 *
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

//服务中心，处理服务的调度和查询
import (
	"aliens/center/core/lbs"
	"encoding/json"
	"aliens/log"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

const SERVICE_NODE_NAME string = "service"

const DEFAULT_LBS string = lbs.LBS_STRATEGY_POLLING

type ServiceCenter struct {
	sync.RWMutex
	zkCon            *zk.Conn
	zkName           string
	serviceRoot      string

	serviceListeners  map[string]ServiceListener  //服务监听
	serviceContainer map[string]*serviceCategory //服务容器 key 服务名 订阅的服务
	factorys         map[string]IServiceFactory  //服务工厂容器 key 服务名


	localServices []IService //本地启动的服务

	lbs string //default polling
}

//启动服务中心客户端
func (this *ServiceCenter) Connect(address string, timeout int, zkName string) {
	this.ConnectCluster([]string{address}, timeout, zkName)
}

func (this *ServiceCenter) ConnectCluster(addressArray []string, timeout int, zkName string) {
	this.zkName = zkName
	//this.serviceFactory = serviceFactory
	c, _, err := zk.Connect(addressArray, time.Duration(timeout)*time.Second)
	if err != nil {
		panic(err)
	}
	this.serviceListeners = make(map[string]ServiceListener)
	this.serviceContainer = make(map[string]*serviceCategory)
	this.factorys = make(map[string]IServiceFactory)
	this.localServices = []IService{}
	this.serviceRoot = NODE_SPLIT + this.zkName + NODE_SPLIT + SERVICE_NODE_NAME
	this.zkCon = c
	this.confirmNode(NODE_SPLIT + this.zkName)
	this.confirmNode(this.serviceRoot)
}

func (this *ServiceCenter) SetLBS(lbs string) {
	this.lbs = lbs
}

//新增服务监听
func (this *ServiceCenter) AddServiceListener(listener ServiceListener) {
	this.Lock()
	defer this.Unlock()
	this.serviceListeners[listener.GetServiceType()] = listener
}

//设置指定服务类型使用的服务工厂类
func (this *ServiceCenter) AddServiceFactory(serviceType string, serviceFactory IServiceFactory) {
	this.Lock()
	defer this.Unlock()
	this.factorys[serviceType] = serviceFactory
}

func (this *ServiceCenter) IsConnect() bool {
	return this.zkCon != nil
}

func (this *ServiceCenter) assert() {
	if this.zkCon == nil {
		panic("mast start service center first")
	}
}

//关闭服务中心
func (this *ServiceCenter) Close() {
	for _, localService := range this.localServices {
		this.ReleaseService(localService)
	}
	if this.zkCon != nil {
		this.zkCon.Close()
	}
}

//更新服务
//func (this *ServiceCenter) UpdateService(service IService) {
//	this.Lock()
//	defer this.Unlock()
//	if this.serviceContainer[service.GetType()] == nil {
//		this.serviceContainer[service.GetType()] = NewServiceCategory(service.GetType(), this.lbs, service.GetDesc())
//	}
//	this.serviceContainer[service.GetType()].updateService(service)
//}

//根据服务类型获取一个空闲的服务节点
func (this *ServiceCenter) AllocService(serviceType string) IService {
	this.RLock()
	defer this.RUnlock()
	//TODO 后续要优化，考虑负载、空闲等因素
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.allocService()
}

//
func (this *ServiceCenter) GetMasterService(serviceType string) IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getMaster()
}

func (this *ServiceCenter) CanHandle(serviceType string, seq int32) bool {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return false
	}
	return serviceCategory.canHandle(seq)
}

func (this *ServiceCenter) GetService(serviceType string, serviceID string) IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.services[serviceID]
	////节点没有取第一个
	//if (service == nil) {
	//	serviceCategory.allocService()
	//}
	//return service
}

func (this *ServiceCenter) GetAllService(serviceType string) []IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getAllService()
	////节点没有取第一个
	//if (service == nil) {
	//	serviceCategory.allocService()
	//}
	//return service
}


func (this *ServiceCenter) GetAllServiceIgnoreID(serviceType string, serviceID string) []IService {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getAllServiceIgnoreID(serviceID)
	////节点没有取第一个
	//if (service == nil) {
	//	serviceCategory.allocService()
	//}
	//return service
}

func (this *ServiceCenter) GetServiceInfo(serviceType string) []string {
	this.RLock()
	defer this.RUnlock()
	serviceCategory := this.serviceContainer[serviceType]
	if serviceCategory == nil {
		return nil
	}
	return serviceCategory.getNodes()
}

//订阅服务  能实时更新服务信息
//func (this *ServiceCenter) SubscribeServices(serviceTypes ...string) {
//	this.assert()
//	for _, serviceType := range serviceTypes {
//		this.SubscribeService(serviceType)
//	}
//}

func (this *ServiceCenter) SubscribeService(serviceType string) {
	this.Lock()
	defer this.Unlock()
	serviceFactory := this.factorys[serviceType]
	if serviceFactory == nil {
		log.Info("service %v factory not register", serviceType)
		return
	}

	path := this.serviceRoot + NODE_SPLIT + serviceType
	desc := this.confirmContentNode(path)
	serviceIDs, _, ch, err := this.zkCon.ChildrenW(path)
	if err != nil {
		log.Info("subscribe service %v error: %v", path, err)
		return
	}
	oldContainer := this.serviceContainer[serviceType]
	serviceCategory := NewServiceCategory(serviceType, this.lbs, desc)

	listener := this.serviceListeners[serviceType]

	for _, serviceID := range serviceIDs {
		if oldContainer != nil {
			oldService := oldContainer.takeoutServiceByID(serviceID)
			if oldService != nil {
				//oldService.SetID(service.GetID())
				serviceCategory.updateService(oldService)
				continue
			}
		}

		data, _, err := this.zkCon.Get(path + NODE_SPLIT + serviceID)
		service := serviceFactory.CreateService(data)
		if service == nil {
			log.Info("%v unmarshal json error : %v", path, err)
			continue
		}
		service.SetID(serviceID)
		service.SetType(serviceType)
		//新服务需要连接上才能更新
		if service.Connect() {
			if listener != nil {
				listener.AddNode(service.GetID())
			}
			serviceCategory.updateService(service)
		}
	}
	if oldContainer != nil {
		for id, service := range oldContainer.services {
			service.Close()
			if listener != nil {
				listener.RemoveNode(id)
			}
		}
	}

	this.serviceContainer[serviceType] = serviceCategory
	go this.openListener(serviceType, path, ch)
}


func (this *ServiceCenter) openListener(serviceType string, path string, ch <-chan zk.Event) {
	event, _ := <-ch
	//更新服务节点信息
	if event.Type == zk.EventNodeChildrenChanged {
		this.SubscribeService(serviceType)
	}
}

//
func (this *ServiceCenter) confirmNode(path string, flags ...int32) bool {
	_, err := this.zkCon.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	return err == nil
}

func (this *ServiceCenter) confirmContentNode(path string, flags ...int32) string {
	_, err := this.zkCon.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		data, _, _ := this.zkCon.Get(path)
		return string(data)
	}
	return ""
}

func (this *ServiceCenter) confirmDataNode(path string, data string) bool {
	byteData := []byte(data)
	_, err := this.zkCon.Create(path, byteData, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		this.zkCon.Set(path, byteData, -1)
	}
	return err == nil
}


func (this *ServiceCenter) ReleaseService(service IService) {
	this.assert()
	servicePath := this.serviceRoot + NODE_SPLIT + service.GetType() + NODE_SPLIT + service.GetID()
	err := this.zkCon.Delete(servicePath, -1)
	if err != nil {
		log.Error("release service %v error : %v", servicePath, err)
	} else {
		log.Info("release %v success", servicePath)
	}
}

//发布服务
func (this *ServiceCenter) PublicService(service IService) bool {
	this.assert()
	if !service.IsLocal() {
		log.Info("service info is invalid")
		return false
	}
	//path string, data []byte, version int32
	data, err := json.Marshal(service)
	if err != nil {
		log.Info("marshal json service data error : %v", err)
		return false
	}
	servicePath := this.serviceRoot + NODE_SPLIT + service.GetType()
	this.confirmDataNode(servicePath, service.GetDesc())
	//id, err := this.zkCon.Create(servicePath + NODE_SPLIT + service.GetType(), data,
	//	zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	id, err := this.zkCon.Create(servicePath+NODE_SPLIT+service.GetID(), data,
		zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Info("create service error : %v", err)
		return false
	}
	this.localServices = append(this.localServices, service)
	log.Info("public %v success : %v-%v", service.GetType(), id, string(data))
	//服务注册在容器
	//this.UpdateService(service)
	return true
}
