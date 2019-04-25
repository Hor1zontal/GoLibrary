/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/3/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package util

import "time"

const (
	DURATION_DAY time.Duration = 24 * time.Hour
	SECOEND_OF_HOUR int64 = 60 * 60
	SECOEND_OF_DAY int64 = 24 * SECOEND_OF_HOUR
)

//获取当天开始时间
func GetTodayBegin() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetMontyBegin(hour int) time.Time {
	year, month, _ := time.Now().Date()
	return time.Date(year, month, 1, hour, 0, 0, 0, time.Local)
}

func GetTodayHourTime(hour int) time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

func GetEmptyTime() time.Time {
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
}

//获取丹田的结束时间
func GetTodayEnd() time.Time {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return tm1.Add(DURATION_DAY)
}

func RefreshTime(lastRefreshTime time.Time, currRefreshTime time.Time, interval int32) (int32, time.Time){
	duration := int32(currRefreshTime.Sub(lastRefreshTime).Seconds())
	if duration > interval {
		addCount := duration / interval
		remainTime := duration % interval
		return addCount, time.Unix(currRefreshTime.Unix()-int64(remainTime), 0)
	}
	return 0, lastRefreshTime
}


//获取当前的日志字符串打印
func GetCurrentDayStr() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func GetTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetDuraton(duration float64) time.Duration {
	return time.Duration(duration) * time.Second
}

