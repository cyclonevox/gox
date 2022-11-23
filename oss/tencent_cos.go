package oss

import (
	`context`
	`crypto/tls`
	`fmt`
	`net/http`
	`net/url`
	`strconv`
	`time`

	`github.com/storezhang/gox/cache`
	`github.com/storezhang/gox/core`
	`github.com/tencentyun/cos-go-sdk-v5`
)

type tencent struct {
	*cos.Client

	// 密钥id
	secretId string
	// 密钥key
	secretKey string
	// CDN加速地址
	cdn string
	// 过期时间
	expiration core.Expiration

	cache cache.Cache
}

func NewTencent(config core.OSSConfig) OSS {
	bu, _ := url.Parse(config.Endpoint)
	t := &tencent{
		Client: cos.NewClient(
			&cos.BaseURL{BucketURL: bu},
			&http.Client{
				Transport: &cos.AuthorizationTransport{
					SecretID:  config.AK,
					SecretKey: config.SK,
					// nolint:gosec
					Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
				},
			}),
		secretId:  config.AK,
		secretKey: config.SK,
		cdn:       config.CDN,

		expiration: config.Expiration,
	}

	if t.cdn != "" {
		t.cache = cache.NewLRUCache(config.DownloadCacheCap)
	}

	return t
}

func (t *tencent) GetUploadURL(req *GetUploadURLReq) (*GetUploadURLRsp, error) {
	var (
		err            error
		preassignedURL *url.URL
	)

	var query = make(url.Values)
	if req.UploadId != "" {
		query.Set("uploadId", req.UploadId)
	}
	if req.PartNumber != 0 {
		query.Set("partNumber", strconv.FormatInt(int64(req.PartNumber), 10))
	}

	if preassignedURL, err = t.Object.GetPresignedURL(
		context.Background(),
		http.MethodPut,
		req.Key,
		t.secretId, t.secretKey,
		time.Duration(t.expiration.Upload)*time.Minute,
		cos.ObjectPutHeaderOptions{
			ContentType:   req.ContentType,
			XOptionHeader: &http.Header{"Access-Control-Expose-Headers": []string{"ETag"}},
		},
	); err != nil {
		return nil, err
	}

	var s string
	if len(query) != 0 {
		s = "&" + query.Encode()
	}

	return &GetUploadURLRsp{URL: fixEscape(preassignedURL.String() + s)}, nil
}

func (t *tencent) GetDownloadURL(req *GetDownloadURLReq) (rsp *GetDownloadURLRsp, err error) {
	if t.cdn != "" {
		if val, ok := t.cache.Get(req.CacheKey()); ok {
			rsp = &GetDownloadURLRsp{URL: val.(string)}

			return
		}
	}

	defer func() {
		if t.cdn != "" && err == nil {
			t.cache.Set(req.CacheKey(), rsp.URL)
		}
	}()

	var disposition string
	if req.DownloadType == core.OSSDownloadTypeDownload {
		disposition = fmt.Sprintf("attachment; filename=%s;filename*=utf-8''%s", req.Filename, req.Filename)
	} else {
		disposition = fmt.Sprintf("inline; filename=%s;filename*=utf-8''%s", req.Filename, req.Filename)

		if req.ContentType == "" {
			if req.ContentType, err = t.getContentType(req.Key); err != nil {
				return nil, err
			}
		}
	}

	if t.cdn != "" {
		query := make(url.Values, 2)
		query.Add("response-content-disposition", disposition)
		if req.ContentType != "" {
			query.Add("response-content-type", req.ContentType)
		}

		rsp = &GetDownloadURLRsp{URL: fixEscape(fmt.Sprintf("%s/%s?%s", t.cdn, req.Key, query.Encode()))}

		return
	}

	expire := time.Duration(t.expiration.Download) * time.Minute
	if req.DownloadType == core.OSSDownloadTypeOpen {
		expire = time.Duration(t.expiration.Open) * time.Minute
	}

	var (
		preassignedURL *url.URL
		options        = &cos.ObjectGetOptions{}
	)

	options.ResponseContentDisposition = disposition
	if req.ContentType != "" {
		options.ResponseContentType = req.ContentType
	}

	// 获取预签名URL
	if preassignedURL, err = t.Object.GetPresignedURL(
		context.Background(),
		http.MethodGet,
		req.Key,
		t.secretId, t.secretKey,
		expire,
		options,
	); err != nil {
		return nil, err
	}

	rsp = &GetDownloadURLRsp{URL: fixEscape(preassignedURL.String())}

	return
}

func (t *tencent) DeleteObject(req *DeleteObjectReq) error {
	if _, err := t.Object.Delete(context.Background(), req.Key); err != nil {
		return err
	}

	return nil
}

func (t *tencent) InitiateMultipartUpload(req *InitiateMultipartUploadReq) (*InitiateMultipartUploadRsp, error) {
	var (
		err    error
		result *cos.InitiateMultipartUploadResult
	)

	if result, _, err = t.Object.InitiateMultipartUpload(
		context.Background(),
		req.Key,
		&cos.InitiateMultipartUploadOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: req.ContentType,
			},
		},
	); err != nil {
		return nil, err
	}

	rsp := &InitiateMultipartUploadRsp{
		Bucket:   result.Bucket,
		Key:      req.Key,
		UploadId: result.UploadID,
	}

	return rsp, nil
}

func (t *tencent) CompleteMultipartUpload(req *CompleteMultipartUploadReq) (*CompleteMultipartUploadRsp, error) {
	var (
		err     error
		result  *cos.CompleteMultipartUploadResult
		options = new(cos.CompleteMultipartUploadOptions)
	)

	for _, part := range req.Part {
		options.Parts = append(options.Parts, cos.Object{
			PartNumber: int(part.PartNumber),
			ETag:       part.Etag,
		})
	}

	if result, _, err = t.Object.CompleteMultipartUpload(
		context.Background(),
		req.Key,
		req.UploadId,
		options,
	); err != nil {
		return nil, err
	}

	rsp := &CompleteMultipartUploadRsp{
		Location: result.Location,
		Bucket:   result.Bucket,
		Key:      req.Key,
		Etag:     result.ETag,
	}

	return rsp, nil
}

func (t *tencent) AbortMultipartUpload(req *AbortMultipartUploadReq) error {
	if _, err := t.Object.AbortMultipartUpload(context.Background(), req.Key, req.UploadId); err != nil {
		return err
	}

	return nil
}

func (t *tencent) ListObject(req *ListObjectReq) (*ListObjectRsp, error) {
	var (
		err    error
		result *cos.BucketGetResult
	)

	if result, _, err = t.Bucket.Get(context.Background(), &cos.BucketGetOptions{
		Prefix:       req.Prefix,
		Delimiter:    req.Delimiter,
		EncodingType: req.EncodingType,
		Marker:       req.Marker,
		MaxKeys:      int(req.MaxKeys),
	}); err != nil {
		return nil, err
	}

	rsp := &ListObjectRsp{
		CommonPrefixes: result.CommonPrefixes,
		Contents:       make([]*Object, 0, len(result.Contents)),
		Delimiter:      result.Delimiter,
		EncodingType:   result.EncodingType,
		IsTruncated:    result.IsTruncated,
		Marker:         result.Marker,
		MaxKeys:        int32(result.MaxKeys),
		Name:           result.Name,
		NextMarker:     result.NextMarker,
		Prefix:         req.Prefix,
	}

	for i := range rsp.CommonPrefixes {
		rsp.CommonPrefixes[i] = rsp.CommonPrefixes[i]
	}

	for _, content := range result.Contents {
		rsp.Contents = append(rsp.Contents, &Object{
			Etag:         content.ETag,
			Key:          content.Key,
			LastModified: content.LastModified,
			Size:         content.Size,
			StorageClass: content.StorageClass,
			Owner: &Owner{
				Id:          content.Owner.ID,
				DisplayName: content.Owner.DisplayName,
			},
		})
	}

	return rsp, nil
}

func (t *tencent) CopyObject(req *CopyObjectReq) (*CopyObjectRsp, error) {
	var (
		err    error
		result *cos.ObjectCopyResult
	)

	req.Source = fmt.Sprintf("%s/%s", t.Client.BaseURL.BucketURL.Host, req.Source)
	if result, _, err = t.Object.Copy(context.Background(), req.Destination, req.Source, nil); err != nil {
		return nil, err
	}

	rsp := &CopyObjectRsp{
		Etag:         result.ETag,
		LastModified: result.LastModified,
	}

	return rsp, nil
}

func (t *tencent) getContentType(key string) (string, error) {
	var (
		err error
		rsp *cos.Response
	)

	if rsp, err = t.Object.Head(context.Background(), key, nil); err != nil {
		return "", err
	}

	return rsp.Header.Get("Content-Type"), nil
}
