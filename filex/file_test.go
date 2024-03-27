package filex

import (
	"fmt"
	"testing"
)

func TestSqlFilter(t *testing.T) {
	names, _ := ListDir("dir")
	fmt.Printf("names:%v\n", names)

	paths, _ := GetDirAllFilePaths("dir/")
	for _, path := range paths {
		fmt.Println(path)
	}
}
