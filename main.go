package main

import (
	"github.com/labstack/echo"
	"github.com/tylerb/graceful"
	"net/http"
	"time"
	"toon-sensitive/trie"
)

type respMeta struct {
	Code    int    `json:"code"`
	Message string `json:"error,omitempty"`
}

type respData struct {
	Keywords []string `json:keywords,omitempty`
	Text     string   `json,omitempty`
}

type resp struct {
	Meta respMeta `json:meta`
	Data respData `json:data,omitempty`
}

func main() {
	e := echo.New()

	e.GET("/", queryWords)
	e.Server.Addr = ":1323"

	//noise := trie.NoiseWords()
	//for char, val := range noise.StopWords {
	//	fmt.Sprintf("%d, %d\n", char, val)
	//}

	graceful.ListenAndServe(e.Server, 5*time.Second)
}

func queryWords(c echo.Context) error {
	word := c.QueryParam("q")

	meta := respMeta{}
	data := respData{}
	res := resp{}

	if word == "" {
		meta.Code = 1001
		meta.Message = "参数q不能为空"
		res.Meta = meta

		return c.JSON(http.StatusOK, res)
	}

	ok, keywords, newText := trie.BlackTrie().Query(word)

	if ok {
		meta.Code = 1002
		meta.Message = "存在敏感字"
		data.Keywords = keywords
		data.Text = newText
	} else {
		meta.Code = 0
		meta.Message = "OK"
	}

	res.Meta = meta
	res.Data = data

	return c.JSON(http.StatusOK, res)
}
