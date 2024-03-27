package oss

import (
	"github.com/AISHU-Technology/kw-go-core/utils"
	"os"
)

type OssConf struct {
	BucketName      string
	Type            string `json:",default=node,options=node|ali|huawei"`
	Endpoint        string `json:",optional"`
	AccessKeyID     string `json:",optional"`
	AccessKeySecret string `json:",optional"`
	MaxRetries      int    `json:",default=2,range=[0:5]"`
}

var globalObs OssConf

func InitOssConf(e OssConf) *OssConf {
	globalObs = e
	return &globalObs
}

type ossEnvConf struct {
	bucketName      string
	accessKeyID     string
	accessKeySecret string
}

func getOssEnvConf(e OssConf) *ossEnvConf {
	var accessKeyId = os.Getenv("AccessKeyID")
	if utils.IsBlank(accessKeyId) {
		accessKeyId = e.AccessKeyID
	}
	var accessKeySecret = os.Getenv("AccessKeySecret")
	if utils.IsBlank(accessKeySecret) {
		accessKeySecret = e.AccessKeySecret
	}
	var bucketName = os.Getenv("BucketName")
	if utils.IsBlank(bucketName) {
		bucketName = e.BucketName
	}
	return &ossEnvConf{
		bucketName:      bucketName,
		accessKeyID:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
}
