package oss

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"io"
	"log"
)

func huaweiUpload(e OssConf, key string, reader io.Reader) (bool, error) {
	c := getOssEnvConf(e)
	obsClient, err := obs.New(c.accessKeyID, c.accessKeySecret, e.Endpoint, obs.WithMaxRetryCount(e.MaxRetries))
	if err != nil {
		log.Printf("Create obsClient error, errMsg: %s", err.Error())
		return false, err
	}
	input := &obs.PutObjectInput{}
	// 指定存储桶名称
	input.Bucket = c.bucketName
	// 指定下载对象，此处以 example/objectname 为例。
	input.Key = key
	input.Body = reader
	// 流式上传本地文件
	output, err := obsClient.PutObject(input)
	if err == nil {
		log.Printf("Put object(%s) under the bucket(%s) successful!StorageClass:%s, ETag:%s\n", input.Key, input.Bucket, output.StorageClass, output.ETag)
		return true, nil
	}
	log.Printf("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
	return false, err
	//fmt.Printf("Put object(%s) under the bucket(%s) fail!\n", input.Key, input.Bucket)
	//if obsError, ok := err.(obs.ObsError); ok {
	//	fmt.Println("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
	//	return false, obsError
	//} else {
	//	fmt.Println("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.")
	//	return false, err
	//}
}
