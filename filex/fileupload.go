package filex

import (
	"github.com/AISHU-Technology/kw-go-core/date"
	"github.com/AISHU-Technology/kw-go-core/errorx"
	"github.com/AISHU-Technology/kw-go-core/oss"
	"github.com/AISHU-Technology/kw-go-core/utils"
	"log"
	"mime/multipart"
	"net/http"
	"path"
	"strings"
)

type UploadConf struct {
	MaxImageSize int64  `json:",default=33554432,config=MAX_IMAGE_SIZE"`
	MaxVideoSize int64  `json:",default=1073741824,config=MAX_VIDEO_SIZE"`
	MaxAudioSize int64  `json:",default=33554432,config=MAX_AUDIO_SIZE"`
	MaxOtherSize int64  `json:",default=10485760,config=MAX_OTHER_SIZE"`
	MaxModelSize int64  `json:",default=104857600,config=MAX_MODEL_SIZE"`
	ServerURL    string `json:",optional"`
}
type UploadInfo struct {
	// File name | 文件名称
	Name string `json:"name"`
	// File path | 文件路径
	Url string `json:"url"`
}

const (
	overSizeError   string = " The file is over size"
	wrongTypeError  string = "The file type is illegal"
	parseFormFailed string = "Failed to parse multiform data"
	failed          string = "failed"
)

type FileInfo struct {
	FileName   string
	FileSuffix string
	File       multipart.File
	Size       int64
	FileType   string
}

func Upload(r *http.Request, e UploadConf, userId int8) (resp *UploadInfo, err error) {
	file, err := UploadValidate(r, e)
	if err != nil {
		log.Printf("the file is upload error,userId=%d=", userId)
		return nil, err
	}
	defer file.File.Close()
	storeFilePath := path.Join(file.FileType, date.GetNowDate(), utils.ToString(userId), GetNewFileName(file.FileSuffix))
	//提交对象存储
	ok, err := oss.Upload(storeFilePath, file.File)
	if !ok {
		return nil, err
	}
	return &UploadInfo{Name: file.FileName, Url: e.ServerURL + storeFilePath}, nil
}

func CheckOverSize(e UploadConf, fileType string, size int64) error {
	if fileType == ImageStr && size > e.MaxImageSize {
		return errorx.NewCodeErrorMsg(overSizeError)
	} else if fileType == VideoStr && size > e.MaxVideoSize {
		return errorx.NewCodeErrorMsg(overSizeError)
	} else if fileType == AudioStr && size > e.MaxAudioSize {
		return errorx.NewCodeErrorMsg(overSizeError)
	} else if fileType == ModelStr && size > e.MaxModelSize {
		return errorx.NewCodeErrorMsg(overSizeError)
	} else if size > e.MaxOtherSize {
		return errorx.NewCodeErrorMsg(overSizeError)
	}
	return nil
}

func UploadValidate(r *http.Request, e UploadConf) (filex *FileInfo, err error) {
	err = r.ParseMultipartForm(e.MaxVideoSize)
	if err != nil {
		log.Print("fail to parse the multipart form")
		return nil, errorx.NewCodeErrorMsg(parseFormFailed)
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Print("the value of file cannot be found")
		return nil, errorx.NewCodeErrorMsg(parseFormFailed)
	}
	pointIndex := GetFileLastIndex(handler.Filename, FilePoint)
	// 判断是否有后缀文件
	if pointIndex == -1 {
		defer file.Close()
		log.Print("reject the file which does not have suffix")
		return nil, errorx.NewCodeErrorMsg(wrongTypeError)
	}
	fileName, fileSuffix := handler.Filename[:pointIndex], handler.Filename[pointIndex+1:]
	// 获取文件类型
	fileType := strings.Split(handler.Header.Get("Content-Type"), FileCharacter)[0]
	if fileType != ImageStr && fileType != VideoStr && fileType != AudioStr && fileType != ModelStr {
		fileType = OtherStr
	}
	err = CheckOverSize(e, fileType, handler.Size)
	if err != nil {
		defer file.Close()
		log.Printf("the file is over size type=%s,size=%d,fileName=%s", fileType, handler.Size, fileName)
		return nil, err
	}
	return &FileInfo{
		FileName:   fileName,
		Size:       handler.Size,
		FileSuffix: fileSuffix,
		File:       file,
		FileType:   fileType,
	}, nil
}
