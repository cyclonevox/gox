package oss

import (
	`fmt`

	`github.com/storezhang/gox/core`
)

type (
	// GetUploadURLReq 获取上传链接的请求
	GetUploadURLReq struct {
		// 对象路径
		Key string `json:"key" validate:"required,max=65535,oss_key"`
		// 内容类型
		ContentType string `default:"application/x-www-form-urlencoded" json:"contentType" validate:"required,max=255"`

		// 分块上传id
		UploadId string `json:"uploadId" validate:"omitempty,max=255"`
		// 分块编号
		PartNumber int32 `json:"partNumber" validate:"omitempty,min=1,max=255"`
	}
	// GetUploadURLRsp 获取下载链接的请求
	GetUploadURLRsp struct {
		URL string `json:"url"`
	}
)

type (
	// GetDownloadURLReq 获取下载链接的请求
	GetDownloadURLReq struct {
		// 对象路径
		Key string `json:"key" validate:"required,oss_key"`
		// 内容类型，留空表示使用实际的contentType
		ContentType string `json:"contentType" validate:"omitempty,max=255"`

		// 下载类型
		DownloadType core.OSSDownloadType `json:"downloadType" validate:"required,oneof=1 2"`
		// 文件另存为名字
		Filename string `json:"filename" validate:"omitempty,max=255"`
	}
	// GetDownloadURLRsp 获取下载链接的返回
	GetDownloadURLRsp struct {
		URL string `json:"url"`
	}
)

type (
	// InitiateMultipartUploadReq 初始化分块上传的请求
	InitiateMultipartUploadReq struct {
		// 对象路径
		Key string `json:"key" validate:"required,max=65535,oss_key"`
		// 内容类型
		ContentType string `default:"application/x-www-form-urlencoded" json:"contentType" validate:"required,max=255"`
	}

	// InitiateMultipartUploadRsp 完成分块上传的请求
	InitiateMultipartUploadRsp struct {
		// 所在桶
		Bucket string `json:"bucket"`
		// 对象路径
		Key string `json:"key"`
		// 分块上传id
		UploadId string `json:"uploadId"`
	}
)

// Part 分块信息
type Part struct {
	// 校验码
	Etag string `json:"etag" validate:"required"`
	// 分块id
	PartNumber int32 `json:"partNumber" validate:"required,min=1"`
}

type (
	// CompleteMultipartUploadReq 完成分块上传的请求
	CompleteMultipartUploadReq struct {
		// 对象路径
		Key string `json:"key" validate:"required,max=65535,oss_key"`
		// 分块上传id
		UploadId string `json:"uploadId" validate:"omitempty,max=255"`
		// 分块信息列表
		Part []*Part `json:"part" validate:"required,dive"`
	}

	// CompleteMultipartUploadRsp 完成分块上传的返回
	CompleteMultipartUploadRsp struct {
		// 位置
		Location string `json:"location"`
		// 所在桶
		Bucket string `json:"bucket"`
		// 对象路径
		Key string `json:"key"`
		// 校验值
		Etag string `json:"etag"`
	}
)

// AbortMultipartUploadReq 中止分块上传的请求
type AbortMultipartUploadReq struct {
	// 对象路径
	Key string `json:"key" validate:"required,max=65535,oss_key"`
	// 分块上传id
	UploadId string `json:"uploadId" validate:"omitempty,max=255"`
}

// DeleteObjectReq 删除对象的请求
type DeleteObjectReq struct {
	// 对象路径
	Key string `json:"key" validate:"required,max=65535,oss_key"`
}

type (
	// CopyObjectReq 复制对象的请求
	CopyObjectReq struct {
		// 源对象路径
		Source string `json:"source" validate:"required,max=65535,oss_key"`
		// 目标路径
		Destination string `json:"destination" validate:"required,max=65535,oss_key,necsfield=Source"`
	}
	// CopyObjectRsp 复制对象的返回
	CopyObjectRsp struct {
		// 校验值
		Etag string `json:"etag"`
		// 最近更新时间
		LastModified string `json:"lastModified"`
	}
)

type Object struct {
	ChecksumAlgorithm []string `json:"checksumAlgorithm"`
	Etag              string   `json:"etag"`
	Key               string   `json:"key"`
	LastModified      string   `json:"lastModified"`
	Size              int64    `json:"size"`
	StorageClass      string   `json:"storageClass"`
	Owner             *Owner   `json:"owner,omitempty"`
}

type Owner struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type (
	ListObjectReq struct {
		Delimiter    string `json:"delimiter" query:"delimiter"`
		EncodingType string `json:"encodingType" query:"encodingType"`
		Marker       string `json:"marker" query:"marker"`
		MaxKeys      int32  `json:"maxKeys" query:"maxKeys" validate:"omitempty,max=1000"`
		Prefix       string `json:"prefix" query:"prefix" validate:"required"`
	}
	ListObjectRsp struct {
		CommonPrefixes []string  `json:"commonPrefixes"`
		Contents       []*Object `json:"contents"`
		Delimiter      string    `json:"delimiter"`
		EncodingType   string    `json:"encodingType"`
		IsTruncated    bool      `json:"isTruncated"`
		Marker         string    `json:"marker"`
		MaxKeys        int32     `json:"maxKeys"`
		Name           string    `json:"name"`
		NextMarker     string    `json:"nextMarker"`
		Prefix         string    `json:"prefix"`
	}
)

func (gr *GetDownloadURLReq) CacheKey() string {
	return fmt.Sprintf("%s-%d-%s-%s", gr.Key, gr.DownloadType, gr.ContentType, gr.Filename)
}
