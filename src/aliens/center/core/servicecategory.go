/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/4/28
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package core

import "aliens/center/core/lbs"
import "sort"
import "encoding/json"

func NewServiceCategory(category string, lbsStr string, desc string) *serviceCategory {
	seqs := []int32{}
	json.Unmarshal([]byte(desc), &seqs)
	seqMaps := make(map[int32]struct{})
	for _, seq := range seqs {
		seqMaps[seq] = struct{}{}
	}
	return &serviceCategory{
		category: category,
		lbs:      lbs.GetLBS(lbsStr),
		services: make(map[string]IService),
		nodes:    []string{},
		seqs:     seqMaps,
	}
}

type serviceCategory struct {
	category string
	lbs      lbs.LBStrategy      //负载均衡策略
	services map[string]IService //服务节点名,和服务句柄
	nodes    []string
	seqs     map[int32]struct{} //能够处理的消息编号
}

//分配一个可用服务
func (this *serviceCategory) allocService() IService {
	nodeName := this.lbs.AllocNode()
	if nodeName == "" {
		return nil
	}
	return this.services[nodeName]
}

func (this *serviceCategory) canHandle(messageSeq int32) bool {
	_, ok := this.seqs[messageSeq]
	return ok
}

//初始化lbs节点信息
func (this *serviceCategory) initLBSNode() {
	nodes := []string{}
	for node, _ := range this.services {
		nodes = append(nodes, node)
	}
	sort.Strings(nodes)
	this.nodes = nodes
	this.lbs.Init(this.nodes)
}

//更新服务
func (this *serviceCategory) updateService(service IService) {
	this.services[service.GetID()] = service
	this.initLBSNode()
}

//取出相同的服务
func (this *serviceCategory) takeoutService(serviceConfig IService) IService {
	//服务地址信息没有变，不需要再连接
	for key, service := range this.services {
		if service.Equals(serviceConfig) {
			delete(this.services, key)
			this.initLBSNode()
			return service
		}
	}
	return nil
}

func (this *serviceCategory) takeoutServiceByID(id string) IService {
	service, ok := this.services[id]
	if ok {
		delete(this.services, id)
	}
	return service
	////服务地址信息没有变，不需要再连接
	//for key, service := range this.services {
	//	if service.Equals(serviceConfig) {
	//		delete(this.services, key)
	//		this.initLBSNode()
	//		return service
	//	}
	//}
	//return nil
}

func (this *serviceCategory) getNodes() []string {
	return this.nodes
}

func (this *serviceCategory) getAllService() []IService {
	results := []IService{}
	for _, service := range this.services {
		results = append(results, service)
	}
	return results
}

func (this *serviceCategory) getAllServiceIgnoreID(serviceID string) []IService {
	results := []IService{}
	for _, service := range this.services {
		if service.GetID() != serviceID {
			results = append(results, service)
		}
	}
	return results
}

func (this *serviceCategory) getMaster() IService {
	//TODO 后续要加一套master-salve机制
	if len(this.nodes) == 0 {
		return nil
	}
	return this.services[this.nodes[0]]
}
