package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

func aliUpload(e OssConf, key string, reader io.Reader) (bool, error) {
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		log.Print("Error:", err)
		return false, err
	}
	c := getOssEnvConf(e)
	client, err := oss.New(e.Endpoint, c.accessKeyID, c.accessKeySecret, oss.SetCredentialsProvider(&provider))
	if err != nil {
		log.Print("Error:", err)
		return false, err
	}
	bucket, err := client.Bucket(c.bucketName)
	if err != nil {
		log.Print(err)
		return false, err
	}
	err = bucket.PutObject(key, reader)
	if err != nil {
		log.Print("===PutObject=====Error:", err)
		return false, err
	}
	return true, nil
}
