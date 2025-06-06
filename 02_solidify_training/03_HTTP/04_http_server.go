package main

import (
	"io"
	"net/http"
)

type handler int

func (h handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	/*response设置响应头*/
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	homePicture := `<img src="https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=0&pn=0&spn=0&di=7490230549689139201&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3759128165%2C2653246275&os=4001073890%2C1746301819&simid=4117141386%2C466172178&adpicid=0&lpn=0&ln=0&fm=&sme=&cg=&bdtype=0&oriquery=&objurl=https%3A%2F%2Fwx3.sinaimg.cn%2Fmw690%2F88e90961ly1hwvqdknjo4j20u0140tav.jpg&fromurl=ippr_z2C%24qAzdH3FAzdH3Fojtk5_z%26e3Bv54AzdH3Fddlml0nmmcAzdH3FPm3dom75G&gsm=&islist=&querylist=&lid=10360685888006148096" alt=""/>`
	aboutPicture := `<img src=https://image.baidu.com/search/detail?ct=503316480&z=0&ipn=d&word=%E5%9B%BE%E7%89%87&hs=0&pn=1&spn=0&di=7482437761027276801&pi=0&rn=1&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cl=2&lm=-1&cs=3142980186%2C2358570891&os=1786321825%2C3288726602&simid=2668098%2C918978925&adpicid=0&lpn=0&ln=0&fm=&sme=&cg=&bdtype=0&oriquery=&objurl=https%3A%2F%2Fww2.sinaimg.cn%2Fmw690%2F005EUiO2ly1hxj8yk8u5oj30m81c37d9.jpg&fromurl=ippr_z2C%24qAzdH3FAzdH3Fojtk5_z%26e3Bv54AzdH3Fc8bc8nmdmdAzdH3FPl5NIrjJO&gsm=&islist=&querylist=&lid=10360685888006148096 alt=""/>`
	switch req.URL.Path {
	case "/":
		io.WriteString(resp, homePicture)
	case "/about":
		io.WriteString(resp, aboutPicture)
	}
}
func main() {
	var h handler
	/*这里不能添加http协议*/
	http.ListenAndServe("localhost:8080", h)
}
