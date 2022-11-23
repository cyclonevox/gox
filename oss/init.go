package oss

import (
	`strings`

	`github.com/storezhang/gox/core`
)

var client OSS

func InitOSS(config core.OSSConfig) {
	switch config.Type {
	case core.OSSTypeTencentCOS:
		client = NewTencent(config)
	default:
		client = NewHuawei(config)
	}
}

func GetOSS() OSS {
	if client == nil {
		panic("nil client")
	}

	return client
}

func fixEscape(s string) string {
	// 解决Golang JSON序列化时的HTML Escape
	s = strings.Replace(s, "\\u003c", "<", -1)
	s = strings.Replace(s, "\\u003e", ">", -1)
	s = strings.Replace(s, "\\u0026", "&", -1)

	return s
}

// 桶名带点时，不使用桶的域名，使用通用域名跟上桶名路径的形式
// 例：https://xx.name.obs.myhuaweicloud.com/test/1/flower.jpg
// 改为：https://obs.myhuaweicloud.com/xx.name/test/1/flower.jpg
func fixHuaweiOBSURLForBucketNameWithPeriod(url string, bucketName string) string {
	if !strings.Contains(bucketName, ".") {
		return url
	}

	// 拆分成https://和obs.myhuaweicloud.com/test/1/flower.jpg
	slices := strings.SplitN(url, bucketName+".", 2)

	// 拼接https://与obs.myhuaweicloud.com/xx.name/test/1/flower.jpg
	return slices[0] + strings.Replace(slices[1], "/", "/"+bucketName+"/", 1)
}
