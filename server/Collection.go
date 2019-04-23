package server
import(
	"regexp"
	"fmt"
	"io"
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	//"database/sql"
	"strings"
	"strconv"
	//"net/http"
	"net/url"
	"encoding/json"
	"time"
	//"log"
	//"math"
	"math/rand"
)
var (
	reg *regexp.Regexp = regexp.MustCompile("\\s+")
	regTitle *regexp.Regexp = regexp.MustCompile(`[\p{Han}\w\pP]+`)
	//reg *regexp.Regexp = regexp.MustCompile(`[\s| ]+`)
	tagReg *regexp.Regexp = regexp.MustCompile("\\p{^Han}")
	StyleReg *regexp.Regexp = regexp.MustCompile("\\d{1,2}-\\d{1,2}$")
)

func WeixinUrlEnc(body io.Reader)error{
	key := "url+='"
	lenkey := len(key)
	buf := bufio.NewReader(body)
	var _url string
	for{
		line,err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line = reg.ReplaceAllString(line,"")
		//fmt.Println(line)
		n := strings.Index(line,key)
		if n != -1 {
			line = line[(n+lenkey):]
			_url += line[:strings.Index(line,"'")]
		}
		if len(line) < lenkey {
			continue
		}
	}
	return ClientHttp(_url,"GET",200,nil, func(body io.Reader)error {
		return ReadList(body)
	})

}
func ReadList(body io.Reader)error{
	//var err error
	//var tmpEntrys []*Entry
	key := "varmsgList="
	lenkey := len(key)
	buf := bufio.NewReader(body)
	for{
		line,err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		line = reg.ReplaceAllString(line,"")
		//fmt.Println(line)
		if len(line) < lenkey {
			continue
		}
		if strings.Contains(line[:lenkey],key) {
			db_ := map[string]interface{}{}
			line = line[lenkey:len(line)-1]
			err = json.Unmarshal([]byte(line),&db_)
			if err != nil {
				panic(err)
			}

			var LastEntTitle string
			err = ViewKvDB(Conf.DeduPath,func(b *bolt.Bucket)error{
				//fmt.Println(LastStr)
				val :=b.Get(LastStr)
				if len(val)>0{
					LastEntTitle = string(val)
				}
				return nil
			})
			if err != nil {
				panic(err)
				//fmt.Println(err)
			}
			//HandDB(Conf.DbPath,func(db *sql.DB){
			//	LastEnt = GetLastEntry(db)
			//})

			var Ens []*Entry
			for _,_li := range db_["list"].([]interface{}){
				li := _li.(map[string]interface{})
				msg :=li["app_msg_ext_info"].(map[string]interface{})
				info:=li["comm_msg_info"].(map[string]interface{})
				en :=&Entry{
				Title:msg["title"].(string),
				Url:"https://mp.weixin.qq.com"+strings.Replace(msg["content_url"].(string),"&amp;","&",-1),
				BaseTime:int64(info["datetime"].(float64)),
				BeginTime:time.Now().Unix(),
				EndTime:time.Now().Unix()}
				if LastEntTitle != "" {
					if en.Title == LastEntTitle {
						break
					}
				}
				Ens = append(Ens, en)
			}
			for i:= len(Ens)-1;i>=0;i--{
				//fmt.Println(Ens[i])
				Ens[i].HandContent()
			}
			for _,en := range Ens {
				EntryList <- en
			}
			fmt.Println("over")
			//GetNoneEntrys(db)
		}
	}
	return nil

}
func Collection() {

	//err := ClientDo("https://weixin.sogou.com/websearch/wexinurlenc_sogou_profile.jsp",nil)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	ur,err := url.Parse(Conf.WeixinUrl)
	if err != nil {
		panic(err)
	}
	key := ur.Query().Get("query")
	//key := res.Request.PostForm.Get("query")
	err = ClientHttp(Conf.WeixinUrl,"GET",200,nil,func(body io.Reader)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			//fmt.Println(err)
			return err
		}
		//fmt.Println(Conf.WeixinUrl)
		//fmt.Println(doc.Text())

		k:= "this.href.substr(a+4+parseInt("
		var str string
		var val int64
		doc.Find("script").EachWithBreak(func(i int,s *goquery.Selection)bool {
			str = s.Text()
			//fmt.Println(str)
			n := strings.Index(str,k)
			if n != -1{
				n+=len(k)+1
				//fmt.Println(str,n,str[n:(n+2)])
				val,err = strconv.ParseInt(str[n:(n+2)],10,64)
				if err != nil {
					panic(err)
				}
				//fmt.Println(str,val)
				return false
			}
			return true
		})
		var href_url string
		exit := false
		doc.Find(".news-list2 li").EachWithBreak(func(i int, s *goquery.Selection)bool {
			text := reg.ReplaceAllString(s.Find(".txt-box").Text(),"")
			//fmt.Println(text,key)
			if strings.Contains(text,key){
				href_url,exit =s.Find(".txt-box .tit a").Attr("href")
				if exit {
					rand.Seed(time.Now().UnixNano())
					b := rand.Intn(90)
					//a := href_url[5+23+b]
					href_url = fmt.Sprintf("https://weixin.sogou.com%s&k=%d&h=%s",href_url,b,string(href_url[strings.Index(href_url,"url=")+4+int(val)+b]))
					//err = ClientDo(href_url,func(body io.Reader,res *http.Response)error{
					header := Conf.Header
					header.Add("Referer",Conf.WeixinUrl)
					fmt.Println(href_url)
					err=ClientHttp_(href_url,"GET",200,nil,header, func(body io.Reader)error {
						return WeixinUrlEnc(body)
					})
					if err != nil {
						fmt.Println(err)
					}
					return false
				}
			}
			return true
		})
		if !exit {
			return io.EOF
			//err = Open(Conf.WeixinUrl)
			//if err != nil {
			//log.Println(err)
			//}
		}
		return err
	})
	if err != nil {
		if err != io.EOF {
			fmt.Println(err)
			time.Sleep(time.Second*3)
			Collection()
		}
		//err = Open(Conf.WeixinUrl)
		//if err != nil {
		//	log.Println(err)
		//}
	}
	fmt.Println("over Coll")

}
