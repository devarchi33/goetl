package infra

import (
	"fmt"
	"time"
)

// Parse8BitsDate 根据标准时间格式的字符串转换为8位的日期字符串（如：20190823），时区为空的话，转换后的日期字符串默认为北京时间
func Parse8BitsDate(datetime string, local *time.Location) string {
	date, err := time.Parse("2006-01-02T15:04:05Z", datetime)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if local == nil {
		zone := "Asia/Shanghai"
		local, _ = time.LoadLocation(zone)
	}

	return date.In(local).Format("20060102")
}
