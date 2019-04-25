/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/3
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 * Desc:
 *     Load Balance Strategy
 *******************************************************************************/
package lbs

import (
	"sync"
)

func NewPollingLBS() *PollingLBS {
	return &PollingLBS{}
}

//轮询负载策略
type PollingLBS struct {
	sync.RWMutex //锁托管给外部
	nodes        []string
	index        int
}

func (this *PollingLBS) Init(nodes []string) {
	this.Lock()
	defer this.Unlock()
	this.nodes = nodes
}

func (this *PollingLBS) AllocNode() string {
	this.Lock()
	defer this.Unlock()
	if this.nodes == nil {
		return ""
	}
	len := len(this.nodes)
	if len == 0 {
		return ""
	}
	if this.index >= len {
		this.index = 0
	}
	node := this.nodes[this.index]
	this.index++
	return node
}
