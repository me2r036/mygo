package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
	"regexp"
)

var (
	staticDomain = "dx916rfs0fg2c.cloudfront.net"
	mediaDomain  = "media.hpstorethailand.com"
	pageName     = "elitex3"
)

func getBody(doc *html.Node) (*html.Node, error) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "link" {
			for i, a := range n.Attr {
				if a.Key == "href" {
					n.Attr[i].Val = "https://" + staticDomain + "/landingpages/" + pageName + "/" + a.Val
					fmt.Println(n.Attr[i].Val)
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "script" {
			for i, a := range n.Attr {
				if a.Key == "src" {
					n.Attr[i].Val = "https://" + staticDomain + "/landingpages/" + pageName + "/" + a.Val
					fmt.Println(n.Attr[i].Val)
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			for i, a := range n.Attr {
				if a.Key == "src" {
					n.Attr[i].Val = "{{media url=\"wysiwyg/" + pageName + "/" + a.Val + "\"}}"
					fmt.Println(n.Attr[i].Val)
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "div" {
			for i, a := range n.Attr {
				if a.Key == "style" {
					re := regexp.MustCompile("background:url\\((.*)\\)")
					// re := regexp.MustCompile(`background:url\((.*)\)`)
					n.Attr[i].Val = re.ReplaceAllString(a.Val, "background:url({{media url=\"wysiwyg/" + pageName + "/" + "${1}\"}})")
					//n.Attr[i].Val = re.ReplaceAllString(a.Val, "background:url()")
					fmt.Println(n.Attr[i].Val)
				}
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return doc, nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func main() {
	doc, _ := html.Parse(strings.NewReader(htm))
	bn, err := getBody(doc)
	if err != nil {
		return
	}
	body := renderNode(bn)
	fmt.Println(html.UnescapeString(body))
}

var htm = `<!DOCTYPE html>
<html>
<head>
    <title></title>
</head>
<body>
    body content
    <p>more content</p>
    <a href="something.png">click</a>
</body>
</html>`