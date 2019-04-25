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

const (
	DefaultVirualSpots        = 40
	DefaultHashKey     string = "key"
)

func NewHashLBS(hashkey string) *HashLBS {
	if hashkey == "" {
		hashkey = DefaultHashKey
	}
	return &HashLBS{
		NewHashRing(0),
		hashkey,
	}
}

//hash 权重负载策略
type HashLBS struct {
	hash *HashRing
	key  string
}

func (this *HashLBS) Init(nodes map[string]int) {
	this.GetHash().UpdateNodes(nodes)
}

func (this *HashLBS) AllocNode() string {
	return this.GetHash().GetNode(this.key)
}

func (this *HashLBS) GetHash() *HashRing {
	if this.hash == nil {
		this.hash = NewHashRing(0)
	}
	return this.hash
}
