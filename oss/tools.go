package oss

import (
	"errors"
	"io"
)

const (
	ali    string = "ali"
	huawei string = "huawei"
	none   string = "none"
)

func Upload(key string, reader io.Reader) (bool, error) {
	switch globalObs.Type {
	case ali:
		return aliUpload(globalObs, key, reader)
	case huawei:
		return huaweiUpload(globalObs, key, reader)
	case none:
		// TODO内部对象存储
		return true, nil
	}
	return false, errors.New("invalid param is not IsExist")
}
