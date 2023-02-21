// @Time : 2023/2/21 5:09 PM
// @Author : zhangguangqiang
// @File : main_test.go
// @Software: GoLand

package initEnv

import (
	"fmt"
	"testing"
)

func TestInfo_GetTagArticleByPage(t *testing.T) {

	srv := NewInfo()
	srv.TagUrl = "https://www.v2ex.com/go/go"
	got, err := srv.GetTagArticleByPage(1)
	if err != nil {
		fmt.Printf(" error: %s \n", err.Error())
		return
	}
	fmt.Println(got)

}
