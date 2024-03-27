package config

import (
	"flag"
	"github.com/AISHU-Technology/kw-go-core/filex"
	"log"
	"os"
)

// InitEnvConf
// 文件参数传递顺序 dev(开发)、test(测试)、pro(生产)
func InitEnvConf(paths ...string) *string {
	var configPath *string
	configEnv := os.Getenv("GO_ENV")
	switch configEnv {
	case "dev":
		configPath = flag.String("f", paths[0], "the config dev file")
	case "test":
		configPath = flag.String("f", paths[1], "the config test file")
	case "prod":
		configPath = flag.String("f", paths[2], "the config prod file")
	default:
		configPath = flag.String("f", paths[0], "the config dev file")
	}
	return configPath
}

func InitConf(path string) *string {
	return &path
}

func ConfLoad(path string, t any) {
	if err := ReadLoad(path, t); err != nil {
		log.Fatalf("error: config file %s, %s", path, err.Error())
	}
}

func ResourceLoad(dir string) (map[string]interface{}, error) {
	arr, err := filex.GetDirAllFilePaths(dir)
	if err != nil {
		log.Fatalf("error: file dir error  %s, %s", dir, err.Error())
		return nil, err
	}
	var maps map[string]interface{}
	for _, a := range arr {
		maps, err = ReadLoadMaps(a, maps)
		if err != nil {
			log.Fatalf("error: config file path=%s, %s", a, err.Error())
			return nil, err
		}
	}
	return maps, nil
}
