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
	"net/http"
	"encoding/json"
	"time"
)
var (
	reg *regexp.Regexp = regexp.MustCompile("\\s+")
	tagReg *regexp.Regexp = regexp.MustCompile("\\p{^Han}")
	StyleReg *regexp.Regexp = regexp.MustCompile("\\d{1,2}-\\d{1,2}$")
)
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
		if len(line) < lenkey {
			continue
		}
		if strings.Contains(line[:lenkey],key) {
			db_ := map[string]interface{}{}
			line = line[lenkey:len(line)-1]
			err = json.Unmarshal([]byte(line),&db_)
			if err != nil {
				fmt.Println(line)
				panic(err)
			}

			var LastEntTitle string
			err = ViewKvDB(Conf.KvDbPath,func(b *bolt.Bucket)error{
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
			fmt.Println("over")
			//GetNoneEntrys(db)
		}
	}
	return nil

}
func Collection() {

	err := ClientDo(Conf.WeixinUrl,func(body io.Reader,res *http.Response)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		key := res.Request.PostForm.Get("query")

		doc.Find(".news-list2 li").EachWithBreak(func(i int, s *goquery.Selection)bool {

			text := reg.ReplaceAllString(s.Find(".txt-box").Text(),"")
			fmt.Println(text,key)
			if strings.Contains(text,key){
				href_url,exit :=s.Find(".txt-box .tit a").Attr("href")
				if exit {
					fmt.Println(href_url,exit)
					ClientDo(href_url,func(body io.Reader,res *http.Response)error{
						err = ReadList(body)
						if err != nil {
							panic(err)
						}
						return nil
					})
					return false
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		panic(err)
	}

}
