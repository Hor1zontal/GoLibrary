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
	LBS_STRATEGY_POLLING string = "polling" //轮询
	LBS_STRATEGY_IPHASH  string = "iphash"  //ip地址hash
	LBS_STRATEGY_WEIGHT  string = "weight"  //权重
)

func GetLBS(lbs string) LBStrategy {
	if lbs == LBS_STRATEGY_POLLING {
		return NewPollingLBS()
	} else {
		return NewPollingLBS()
	}
}

type LBStrategy interface {
	Init(nodes []string) //更新所有的负载节点信息

	AllocNode() string //分配可用节点

}
