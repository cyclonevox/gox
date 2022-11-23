package oss

import (
	`fmt`
	`net/url`
	`strconv`

	`github.com/huaweicloud/huaweicloud-sdk-go-obs/obs`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/cache`
	`github.com/storezhang/gox/core`
)

type huawei struct {
	*obs.ObsClient

	// 桶名
	bucket string
	// CDN加速地址
	cdn string
	// 过期时间
	expiration core.Expiration

	cache cache.Cache
}

func NewHuawei(config core.OSSConfig) OSS {
	o, _ := obs.New(config.AK, config.SK, config.Endpoint, obs.WithSslVerify(false))

	h := &huawei{
		ObsClient: o,

		bucket:     config.Bucket,
		expiration: config.Expiration,
		cdn:        config.CDN,
	}

	if h.cdn != "" {
		h.cache = cache.NewLRUCache(config.DownloadCacheCap)
	}

	return h
}

func (h *huawei) GetUploadURL(req *GetUploadURLReq) (*GetUploadURLRsp, error) {
	var (
		err    error
		output *obs.CreateSignedUrlOutput
	)

	var query = make(map[string]string)
	if req.UploadId != "" {
		query["uploadId"] = req.UploadId
	}
	if req.PartNumber != 0 {
		query["partNumber"] = strconv.FormatInt(int64(req.PartNumber), 10)
	}

	if output, err = h.ObsClient.CreateSignedUrl(&obs.CreateSignedUrlInput{
		Method:  obs.HttpMethodPut,
		Bucket:  h.bucket,
		Key:     req.Key,
		Expires: h.expiration.Upload * 60,
		Headers: map[string]string{
			obs.HEADER_CONTENT_TYPE_CAML: req.ContentType,
			"x-amz-acl":                  string(obs.AclPublicRead),
		},
		QueryParams: query,
	}); err != nil {
		return nil, err
	}

	return &GetUploadURLRsp{URL: fixHuaweiOBSURLForBucketNameWithPeriod(fixEscape(output.SignedUrl), h.bucket)}, nil
}

func (h *huawei) GetDownloadURL(req *GetDownloadURLReq) (rsp *GetDownloadURLRsp, err error) {
	if h.cdn != "" {
		if val, ok := h.cache.Get(req.CacheKey()); ok {
			rsp = &GetDownloadURLRsp{URL: val.(string)}

			return
		}
	}

	defer func() {
		if h.cdn != "" && err == nil {
			h.cache.Set(req.CacheKey(), rsp.URL)
		}
	}()

	var disposition string
	if req.DownloadType == core.OSSDownloadTypeDownload {
		disposition = fmt.Sprintf("attachment; filename=%s;filename*=utf-8''%s", req.Filename, req.Filename)
	} else {
		disposition = fmt.Sprintf("inline; filename=%s;filename*=utf-8''%s", req.Filename, req.Filename)

		if req.ContentType == "" {
			if req.ContentType, err = h.getContentType(req.Key); err != nil {
				return nil, err
			}
		}
	}

	if h.cdn != "" {
		query := make(url.Values, 2)
		query.Add("response-content-disposition", disposition)
		if req.ContentType != "" {
			query.Add("response-content-type", req.ContentType)
		}

		rsp = &GetDownloadURLRsp{URL: fixEscape(fmt.Sprintf("%s/%s?%s", h.cdn, req.Key, query.Encode()))}

		return
	}

	expire := h.expiration.Download * 60
	if req.DownloadType == core.OSSDownloadTypeOpen {
		expire = h.expiration.Open * 60
	}

	var (
		output      *obs.CreateSignedUrlOutput
		queryParams = make(map[string]string)
	)

	queryParams[obs.PARAM_RESPONSE_CONTENT_DISPOSITION] = disposition
	if req.ContentType != "" {
		queryParams[obs.PARAM_RESPONSE_CONTENT_TYPE] = req.ContentType
	}

	if output, err = h.ObsClient.CreateSignedUrl(&obs.CreateSignedUrlInput{
		Method:      obs.HttpMethodGet,
		Bucket:      h.bucket,
		Key:         req.Key,
		Expires:     expire,
		QueryParams: queryParams,
	}); err != nil {
		return nil, err
	}

	rsp = &GetDownloadURLRsp{
		URL: fixHuaweiOBSURLForBucketNameWithPeriod(fixEscape(output.SignedUrl), h.bucket),
	}

	return
}

func (h *huawei) DeleteObject(req *DeleteObjectReq) error {
	if _, err := h.ObsClient.DeleteObject(&obs.DeleteObjectInput{
		Bucket: h.bucket,
		Key:    req.Key,
	}); err != nil {
		return err
	}

	return nil
}

func (h *huawei) InitiateMultipartUpload(req *InitiateMultipartUploadReq) (*InitiateMultipartUploadRsp, error) {
	var (
		err    error
		output *obs.InitiateMultipartUploadOutput
	)

	if output, err = h.ObsClient.InitiateMultipartUpload(&obs.InitiateMultipartUploadInput{
		ObjectOperationInput: obs.ObjectOperationInput{
			ACL:    obs.AclPublicRead,
			Bucket: h.bucket,
			Key:    req.Key,
		},
		ContentType: req.ContentType,
	}); err != nil {
		return nil, err
	}

	rsp := &InitiateMultipartUploadRsp{
		Bucket:   output.Bucket,
		Key:      req.Key,
		UploadId: output.UploadId,
	}

	return rsp, nil
}

func (h *huawei) CompleteMultipartUpload(req *CompleteMultipartUploadReq) (*CompleteMultipartUploadRsp, error) {
	var (
		err    error
		output *obs.CompleteMultipartUploadOutput
	)

	parts := make([]obs.Part, 0, len(req.Part))
	for _, part := range req.Part {
		parts = append(parts, obs.Part{
			PartNumber: int(part.PartNumber),
			ETag:       part.Etag,
		})
	}

	if output, err = h.ObsClient.CompleteMultipartUpload(&obs.CompleteMultipartUploadInput{
		Bucket:   h.bucket,
		Key:      req.Key,
		UploadId: req.UploadId,
		Parts:    parts,
	}); err != nil {
		return nil, err
	}

	rsp := &CompleteMultipartUploadRsp{
		Location: output.Location,
		Bucket:   output.Bucket,
		Key:      req.Key,
		Etag:     output.ETag,
	}

	return rsp, err
}

func (h *huawei) AbortMultipartUpload(req *AbortMultipartUploadReq) error {
	if _, err := h.ObsClient.AbortMultipartUpload(&obs.AbortMultipartUploadInput{
		Bucket:   h.bucket,
		Key:      req.Key,
		UploadId: req.UploadId,
	}); err != nil {
		return err
	}

	return nil
}

func (h *huawei) ListObject(req *ListObjectReq) (*ListObjectRsp, error) {
	var (
		err    error
		output *obs.ListObjectsOutput
	)

	if output, err = h.ObsClient.ListObjects(&obs.ListObjectsInput{
		ListObjsInput: obs.ListObjsInput{
			Prefix:       req.Prefix,
			MaxKeys:      int(req.MaxKeys),
			Delimiter:    req.Delimiter,
			EncodingType: req.EncodingType,
		},
		Bucket: h.bucket,
		Marker: req.Marker,
	}); err != nil {
		return nil, err
	}

	rsp := &ListObjectRsp{
		CommonPrefixes: output.CommonPrefixes,
		Contents:       make([]*Object, 0, len(output.Contents)),
		Delimiter:      output.Delimiter,
		EncodingType:   output.EncodingType,
		IsTruncated:    output.IsTruncated,
		Marker:         output.Marker,
		MaxKeys:        int32(output.MaxKeys),
		Name:           output.Name,
		NextMarker:     output.NextMarker,
		Prefix:         output.Prefix,
	}

	for i := range rsp.CommonPrefixes {
		rsp.CommonPrefixes[i] = rsp.CommonPrefixes[i]
	}

	for _, content := range output.Contents {
		rsp.Contents = append(rsp.Contents, &Object{
			Etag:         content.ETag,
			Key:          content.Key,
			LastModified: content.LastModified.Format(gox.DefaultTimeLayout),
			Size:         content.Size,
			StorageClass: string(content.StorageClass),
			Owner: &Owner{
				Id:          content.Owner.ID,
				DisplayName: content.Owner.DisplayName,
			},
		})
	}

	return rsp, nil
}

func (h *huawei) CopyObject(req *CopyObjectReq) (*CopyObjectRsp, error) {
	var (
		err    error
		output *obs.CopyObjectOutput
	)

	if output, err = h.ObsClient.CopyObject(&obs.CopyObjectInput{
		ObjectOperationInput: obs.ObjectOperationInput{
			Bucket: h.bucket,
			Key:    req.Destination,
		},
		CopySourceBucket: h.bucket,
		CopySourceKey:    req.Source,
	}); err != nil {
		return nil, err
	}

	rsp := &CopyObjectRsp{
		Etag:         output.ETag,
		LastModified: output.LastModified.Format(gox.DefaultTimeLayout),
	}

	return rsp, nil
}

func (h *huawei) getContentType(key string) (string, error) {
	var (
		err error
		rsp *obs.GetObjectMetadataOutput
	)

	if rsp, err = h.ObsClient.GetObjectMetadata(&obs.GetObjectMetadataInput{
		Bucket: h.bucket,
		Key:    key,
	}); err != nil {
		return "", err
	}

	return rsp.ContentType, nil
}
