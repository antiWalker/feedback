package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangshizebin/jiebago"
	"net/http"
)

type MessageInfo struct {
	Message string `json:"message" form:"message"`
}

func main() {
	r := gin.Default()
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		return
	}
	JieBaGo := jiebago.NewJieBaGo()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/feedback", func(c *gin.Context) {
		message := &MessageInfo{}
		if err := c.ShouldBind(&message); err == nil {
			message.Message = autoDealMessage(message.Message, *JieBaGo)
			c.JSON(200, message)
		} else {
			c.JSON(200, gin.H{
				"err": err.Error(),
			})
		}
	})
	err = r.Run()
	if err != nil {
		return
	}
}

// 自动化回复消息 三个分类。两个大类和一个其他分类。分别是物流类，商品类，默认类。
func autoDealMessage(sentence string, baGo jiebago.JieBaGo) string {
	// 提取带权重的关键词，即Tag标签
	keywordsWeight := baGo.ExtractKeywordsWeight(sentence, 20)

	//fmt.Println("提取带权重的关键词：", keywordsWeight)
	ExpressSet := make(map[string]struct{})
	GoodsSet := make(map[string]struct{})
	//从数据库取出来的预先设置好的两个大分类。暂时用数组模拟数据，分别是物流类，商品类，
	ExpressArr := [5]string{"物流", "快递", "包裹", "自提柜", "驿站"}
	for _, value := range ExpressArr {
		ExpressSet[value] = struct{}{}
	}
	GoodsArr := [5]string{"新鲜", "烂了", "有味", "商品", "太小"}
	for _, value := range GoodsArr {
		GoodsSet[value] = struct{}{}
	}
	hitWord := false
	message := ""
	catID := 0
	for _, val := range keywordsWeight {
		if val.Weight > 0.5 {
			//在两个分类里分别查找关键词
			if _, exists := ExpressSet[val.Word]; exists {
				message = fmt.Sprintf(ExpressText, val.Word)
				catID = ExpressId
				hitWord = true
				break
			}
			if _, exists := GoodsSet[val.Word]; exists {
				message = fmt.Sprintf(GoodsText, val.Word)
				catID = GoodsId
				hitWord = true
				break
			}
		}
	}
	if hitWord == false {
		message = DefaultText
		catID = DefaultId
	}
	//message store into database func
	err := messageToDB(message, catID)
	if err != nil {
		return ""
	}
	return message
}

// 把采集到的反馈内容和分类insert到bd
func messageToDB(message string, catID int) error {
	//todo ... insert into db
	return nil
}
