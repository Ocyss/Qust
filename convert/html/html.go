package html

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type Data struct {
	*html.Node
}

func New(body *html.Node) *Data {
	return &Data{Node: body}
}

// Find 查找全部
func (h *Data) Find(expr string) (res []*Data) {
	list := htmlquery.Find(h.Node, expr)
	res = make([]*Data, len(list))
	for i, v := range list {
		res[i] = &Data{Node: v}
	}
	return
}

// FindOne 查找单个
func (h *Data) FindOne(expr string) *Data {
	node := htmlquery.FindOne(h.Node, expr)
	return &Data{Node: node}
}

// Text 获取文本数据
func (h *Data) Text() string {
	return htmlquery.InnerText(h.Node)
}

// Attr 获取属性
func (h *Data) Attr(name string) string {
	return htmlquery.SelectAttr(h.Node, name)
}
