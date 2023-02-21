package initEnv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/antchfx/htmlquery"
	"github.com/duke-git/lancet/v2/convertor"
)

type (
	Info struct {
		// 标签
		Tags TagS
		// 当前标签页
		TagUrl string
		// 当前文章页
		ArticleUrl string
		// 当前分页
		Page int
		// 当前评论分页
		CommentPage int
		// 当前打开的文章数组
		ArticleS ArticleS
		// 当前打开的文章
		ArticleInfo *ArticleInfo
	}
	ArticleS []*ArticleInfo
	TagS     []*TagInfo

	ArticleInfo struct {
		// 标题
		Title string
		// 内容
		Content string
		// 标签
		Tag string
		// 发布时间
		CreateTime string
		// 作者
		Author string
		// 评论数
		CommentNum int
		// 评论
		CommentS []*CommentInfo
		// 连接
		Url string
	}
	CommentInfo struct {
		// 作者
		Author string
		// 发布时间
		CreateTime string
		// 内容
		Content string
	}
	TagInfo struct {
		Topics  int      `json:"topics"`
		Aliases []string `json:"aliases"`
		Name    string   `json:"name"`
		Title   string   `json:"title"`
		Url     string   `json:"url,omitempty"`
	}
)

func (t *ArticleS) ToArray() []string {
	list := make([]string, 0)
	for _, v := range *t {
		str := fmt.Sprintf(`标题: %s 
	[作者: %s 时间: %s 评论数: %d]`, v.Title, v.Author, v.CreateTime, v.CommentNum)
		list = append(list, str)
	}
	return list
}

func (t *ArticleS) CheckTitle(title string) *ArticleInfo {
	for _, v := range *t {
		str := fmt.Sprintf(`标题: %s 
	[作者: %s 时间: %s 评论数: %d]`, v.Title, v.Author, v.CreateTime, v.CommentNum)
		if str == title {
			return v
		}
	}
	return nil
}

func (t *TagS) ToArray() []string {
	list := make([]string, 0)
	for _, v := range *t {
		list = append(list, v.Title)
	}
	return list
}
func (t *TagS) CheckTitle(title string) string {
	for _, v := range *t {
		if v.Title == title {
			return v.Url
		}
	}
	return ""
}

func NewInfo() *Info {
	tags := make([]*TagInfo, 0)
	// articleS := make([]*ArticleInfo, 0)
	return &Info{
		Tags:        tags,
		Page:        1,
		CommentPage: 1,
		// ArticleS:    articleS,
	}
}

func (i *Info) Process() {
	// 获取全部标签
	err := i.InitTags()
	if err != nil {
		fmt.Printf("获取全部标签异常: %s \n", err.Error())
		return
	}

	// 选择标签
	i.SelectTag()

	// 打开标签页文章列表(分页)
	err = i.GetTagArticleByPage(1)
	if err != nil {
		fmt.Printf("打开标签页异常: %s \n", err.Error())
		return
	}
	// 选择文章
	i.SelectArticle()

	// 打开文章详情
	err = i.GetArticle()
	if err != nil {
		fmt.Printf("打开文章详情异常: %s \n", err.Error())
		return
	}

	// 选择文章
	i.SelectArticle()
	// 获取评论列表(分页)
	// i.GetArticleAndCommentByPage(1)

}

func (i *Info) GetArticle() error {
	doc, err := htmlquery.LoadURL(i.ArticleUrl)
	if err != nil {
		return err
	}

	a := htmlquery.FindOne(doc, `//*[@id="Main"]/div[2]/div[2]`)
	if a != nil {
		i.ArticleInfo.Content = htmlquery.InnerText(a)
	}
	fmt.Println(fmt.Sprintf(`%s
        [作者: %s 时间: %s 评论数: %d]
        
		%s
`, i.ArticleInfo.Title, i.ArticleInfo.Author, i.ArticleInfo.CreateTime, i.ArticleInfo.CommentNum, i.ArticleInfo.Content))
	return nil
}

func (i *Info) GetArticleAndCommentByPage(page int) (*ArticleInfo, error) {
	Commentlist := make([]*CommentInfo, 0)
	return &ArticleInfo{
		Title:      "",
		Content:    "",
		Tag:        "",
		CreateTime: "",
		Author:     "",
		CommentS:   Commentlist,
	}, nil
}

func (i *Info) GetTagArticleByPage(page int) error {
	url := fmt.Sprintf("%s?p=%d", i.TagUrl, page)

	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return err
	}

	tr := htmlquery.Find(doc, `//*[@id="TopicsNode"]/div`)
	for _, row := range tr {
		articleInfo := &ArticleInfo{}

		a := htmlquery.FindOne(row, `//td[3]/span[2]/strong[1]/a`)
		if a != nil {

			articleInfo.Author = htmlquery.InnerText(a)
		}

		a = htmlquery.FindOne(row, `//td[3]/span[2]/span`)
		if a != nil {

			articleInfo.CreateTime = htmlquery.InnerText(a)
		}

		a = htmlquery.FindOne(row, `//td[4]/a`)
		if a != nil {

			commentNum, _ := convertor.ToInt(htmlquery.InnerText(a))
			articleInfo.CommentNum = int(commentNum)
		}

		a = htmlquery.FindOne(row, "//td[3]/span[1]/a")
		if a != nil {
			articleInfo.Title = htmlquery.InnerText(a)
			articleInfo.Url = fmt.Sprintf(`https://www.v2ex.com/%s`, htmlquery.SelectAttr(a, "href"))
		}
		if articleInfo.Title == "" {
			continue
		}

		// fmt.Println(articleInfo)
		i.ArticleS = append(i.ArticleS, articleInfo)
	}

	return nil
}

func (i *Info) InitTags() error {
	url := "https://www.v2ex.com/api/nodes/list.json?fields=name,title,topics,aliases&sort_by=topics&reverse=1"
	method := "GET"

	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("authority", "www.v2ex.com")
	// req.Header.Add("cookie", "_ga=GA1.2.334168341.1588164896; A2=\"2|1:0|10:1675145417|2:A2|56:MjkzZGRjMzEyNGVhYmJkY2MwYmFmZjdlYTYyNTQyNDI3MjNmNWQxNQ==|3c5d17adff64c68ad47b40935e0f828f7e23ce629b8a12dde7f1d7a9b391a379\"; V2EX_LANG=zhcn; PB3_SESSION=\"2|1:0|10:1676858537|11:PB3_SESSION|36:djJleDoyMTYuMjQuMTg3LjY4OjI2NjQ3OTA2|68613d1281fe8ff364a6d8bd6d84157ea356ad5f21695b6de3df04f871a445cc\"; _gid=GA1.2.540327019.1676858540; V2EX_REFERRER=\"2|1:0|10:1676858555|13:V2EX_REFERRER|12:QmlybEdveQ==|b7d5814ec7c037e9e9b9ed96a9b646ddd62f8fd033937e6e6c25de22e96a847b\"; V2EX_TAB=\"2|1:0|10:1676949213|8:V2EX_TAB|8:cGxheQ==|978afe3e23003e8616e21a61b40791f4f168e01f188fabaf60ca9a51ffddd756\"; _gat=1")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("User-Agent", "apifox/1.0.0 (https://www.apifox.cn)")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(body))

	var resp []*TagInfo
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	for _, v := range resp {
		v.Url = fmt.Sprintf(`https://www.v2ex.com/go/%s`, v.Name)
		i.Tags = append(i.Tags, v)
	}
	return nil
}

func (i *Info) SelectTag() {
	action := ""
	prompt := &survey.Select{
		Message: "请指定标签[可搜索或者选择]",
		Options: i.Tags.ToArray(),
	}
	err := survey.AskOne(prompt, &action)
	if err != nil {
		return
	}

	if action == "" {
		fmt.Println("标签不能为空")
		return
	}
	url := i.Tags.CheckTitle(action)
	if url == "" {
		fmt.Println("标签不存在,请重新选择")
		return
	}

	i.TagUrl = url
}

func (i *Info) SelectArticle() {
	action := ""
	tmpArr := make([]string, 0)
	tmpArr = append(tmpArr, "返回上级")
	tmpArr = append(tmpArr, i.ArticleS.ToArray()...)

	prompt := &survey.Select{
		Message: "请指定帖子[可搜索或者选择]",
		Options: tmpArr,
	}
	err := survey.AskOne(prompt, &action)
	if err != nil {
		return
	}

	if action == "" {
		fmt.Println("帖子不能为空")
		return
	}

	if action == "返回上级" {
		i.SelectTag()
		// 打开标签页文章列表(分页)
		err = i.GetTagArticleByPage(1)
		if err != nil {
			fmt.Printf("打开标签页异常: %s \n", err.Error())
			return
		}
		// 选择文章
		i.SelectArticle()
		return
	}

	res := i.ArticleS.CheckTitle(action)
	if res == nil {
		fmt.Println("帖子不存在,请重新选择")
		return
	}

	i.ArticleUrl = res.Url
	i.ArticleInfo = res
}
