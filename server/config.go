package server
import(
	"net/http"
	//"net/url"
	"github.com/BurntSushi/toml"
	"os"
)
type Config struct {
	Proxy string
	Port string
	DbPath string
	KvDbPath string
	DeduPath string
	Templates string
	Static string
	Header http.Header
	WeixinUrl string
	Coll bool
	//UserInfo *url.Values
	//UserArr []string
	//Site []*SitePage
}
func (self *Config) Save(fileName string){
	fi,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,0777)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	e := toml.NewEncoder(fi)
	err = e.Encode(self)
	if err != nil {
		panic(err)
	}
}
func NewConfig(fileName string)  *Config {
	var c Config
	_,err := os.Stat(fileName)
	if err != nil {
		//c.UserInfo=&url.Values{
		//"username":[]string{""},
		//"password":[]string{""},
		//"randCode":[]string{""}}
		//c.UserArr=[]string{"lqylqjd","lqylxhsq","lqyyhsq","lqyjpc"}
		c.Coll = true
		c.Proxy = ""
		c.KvDbPath="MyKV.db"
		c.DeduPath="dedu.db"
		c.Static = "static"
		c.Port=":8080"
		c.DbPath = "foo.db"
		c.Templates = "./templates/*"
		c.WeixinUrl = "https://weixin.sogou.com/weixin?type=1&s_from=input&query=longquanjy&ie=utf8"
		c.Header = http.Header{
			"Content-Type":[]string{"application/x-www-form-urlencoded"},
			"Accept":[]string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"},
			"Connection":[]string{"keep-alive"},
			"Accept-Encoding":[]string{"gzip, deflate, sdch"},
			"Accept-Language":[]string{"zh-CN,zh;q=0.8"},
			"User-Agent":[]string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/58.0.3029.110 Chrome/58.0.3029.110 Safari/537.36"}}
		//c.Site = []*SitePage{
		//	NewSitePage(
		//		"http://www.ccgp-sichuan.gov.cn/CmsNewsController.do?method=recommendBulletinList&moreType=provincebuyBulletinMore&channelCode=sjcg1&rp=25&page=1",
		//		"page",
		//		".list-info .info li",
		//		".time.curr|text",
		//		"022006-01",
		//		"a|href",
		//		"a .title|text",
		//		560)}
		c.Save(fileName)
	}else{
		if _,err := toml.DecodeFile(fileName,&c);err != nil {
			panic(err)
		}
	}
	return &c
}
