package oss

type OSS interface {
	// GetUploadURL 获取上传地址
	GetUploadURL(*GetUploadURLReq) (*GetUploadURLRsp, error)
	// GetDownloadURL 获取下载链接
	GetDownloadURL(*GetDownloadURLReq) (*GetDownloadURLRsp, error)

	// InitiateMultipartUpload 初始化分块上传
	InitiateMultipartUpload(*InitiateMultipartUploadReq) (*InitiateMultipartUploadRsp, error)
	// CompleteMultipartUpload 完成分块上传
	CompleteMultipartUpload(*CompleteMultipartUploadReq) (*CompleteMultipartUploadRsp, error)
	// AbortMultipartUpload 中止分块上传
	AbortMultipartUpload(*AbortMultipartUploadReq) error

	// DeleteObject 删除对象
	DeleteObject(*DeleteObjectReq) error
	// CopyObject 拷贝对象
	CopyObject(*CopyObjectReq) (*CopyObjectRsp, error)
	// ListObject 列出对象
	ListObject(*ListObjectReq) (*ListObjectRsp, error)
}
