package core

const (
	// OSSTypeTencentCOS 腾讯云COS
	OSSTypeTencentCOS OSSType = 1
	// OSSTypeHuaweiOBS 华为云OBS
	OSSTypeHuaweiOBS OSSType = 2
)

const (
	// OSSDownloadTypeDownload 直接下载
	OSSDownloadTypeDownload OSSDownloadType = 1
	// OSSDownloadTypeOpen 打开
	OSSDownloadTypeOpen OSSDownloadType = 2
)

type (
	// OSSType 对象存储服务提供商类型
	OSSType int8
	// OSSDownloadType 对象存储对象下载方式
	OSSDownloadType int8
)

// Expiration 对象存储过期时间配置
type Expiration struct {
	// 上传
	Upload int `default:"720" yaml:"upload" validate:"required,max=4320"`
	// 下载
	Download int `default:"1440" yaml:"download" validate:"required,max=4320"`
	// 打开
	Open int `default:"1440" yaml:"open" validate:"required,max=4320"`
}

// OSSConfig 对象存储配置
type OSSConfig struct {
	// 环境
	Environment string `yaml:"environment"`
	// 服务名
	ServiceName string `yaml:"serviceName"`

	// 对象存储提供商
	Type OSSType `yaml:"type"`
	// 桶访问域名
	Endpoint string `yaml:"endpoint"`

	AK string `yaml:"ak"`
	SK string `yaml:"sk"`

	// 桶名
	Bucket string `yaml:"bucket"`
	// cdn域名
	CDN string `yaml:"cdn"`

	// 获取下载链接时的缓存大小，使用lru策略进行缓存，默认32767
	DownloadCacheCap int `default:"32767" yaml:"downloadCacheCap"`

	// 过期时间
	Expiration Expiration `yaml:"expiration"`
}
