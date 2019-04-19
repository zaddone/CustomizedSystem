package server
import(
	"regexp"
	"io"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/boltdb/bolt"
	"net/url"
	"time"
	"math/rand"
	"bufio"
	"strings"
	"encoding/json"
)
var (
	ClassID string
)

func LoginSite() error {

	u,err := url.Parse("http://jcpt.chengdu.gov.cn/cdform/cdmanage/frameset/welcome.jsp")
	if err != nil {
		panic(err)
	}
	err = ClientHttp(u.String(),"GET",200,nil,func(body io.Reader)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return err
		}
		var hr string
		doc.Find(".work ul li h3 a").EachWithBreak(func(i int, s *goquery.Selection)bool {
			if s.Text() =="内容发布" {
				hr,_ = s.Attr("href")
				err = inSite1("http://jcpt.chengdu.gov.cn/cdform/cdmanage/frameset/"+hr)
				//fmt.Println(w2)
				return false
			}
			return true
		})
		if hr == "" {
			err = fmt.Errorf("hr == nil")
		}
		return err
	})
	return err
}
func inSite1(_url string) error {

	//time.Sleep(time.Second*5)
	tu := "http://jcpt.chengdu.gov.cn/uycyw/front/back/notice.jsp?fn=null"

	return  ClientHttp(tu,"GET",200,nil,func(body io.Reader)error{
	return ClientHttp(_url,"GET",200,nil,func(body io.Reader)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return err
		}
		u,b := doc.Find("iframe").Attr("src")
		//fmt.Println(_url)
		fmt.Println(u,b)
		con_,err := url.Parse(u)
		if err != nil {
			return err
		}
		//fmt.Println(con_.RequestURI())
		qu := con_.Query()
		q := url.Values{}
		q.Set("Artifact",qu.Get("Artifact"))
		q.Set("fn",qu.Get("fn"))
		return  ClientHttp(u,"GET",200,nil,func(body io.Reader)error{
			rand.Seed(time.Now().UnixNano())
			q.Set("r_",fmt.Sprintf("%.16f",rand.Float64()))
			return ClientHttp(con_.Scheme +"://"+con_.Host+"/uycyw/uymanage/tree/menutree.htm?"+q.Encode(),"GET",200,nil,func(body io.Reader)error{
				buf := bufio.NewReader(body)
				key := "varzNodes="
				keyLen := len(key)
				for{
					line,err := buf.ReadString('\n')
					if err != nil {
						if err == io.EOF {
							break
						}
						panic(err)
					}
					line = reg.ReplaceAllString(line,"")
					if len(line) <= keyLen{
						continue
					}
					//fmt.Println(line)
					if strings.Contains(line[:keyLen],key){
						db_:=[]map[string]interface{}{}
						line = line[keyLen:len(line)-1]
						err = json.Unmarshal([]byte(line),&db_)
						if err != nil {
							panic(err)
						}
						for _,d := range db_{
							if d["name"].(string) == "招聘求职"{
								ClassID = d["id"].(string)
								//return nil
							}
						}
						if ClassID == ""{
							panic("--")
						}
						return nil
					}
				}
				return nil
			})
		})
	})
	})

	//if err != nil {
	//	return err
	//}
	//cdform_,err := url.Parse(_url)
	//if err != nil {
	//	return err
	//}
	//q_ := cdform_.Query()
	////con_ :=q_.Get("url")
	//con_,err := url.Parse(q_.Get("url"))
	//if err != nil {
	//	return err
	//}
	////con_.Query()
	//h := url.QueryEscape(q_.Get("fn"))
	//p_ := con_.String() +"&"+h+"&fn="+h
	//q := con_.Query()
	//q.Set("fn",q_.Get("fn"))
	////q.Set(q_.Get("fn"),"")
	////q.Set("fn","jclongquanyiqu/longquanjiedao")
	////con_.EscapedPath()
	////p_ := con_.Scheme +"://"+con_.Host+con_.Path+"?"+q.Encode()
	////p_ := con_.String()+"&"+v.Encode()
	//fmt.Println(p_)
	//return  ClientHttp(p_,"GET",200,nil,func(body io.Reader)error{
	//	rand.Seed(time.Now().UnixNano())
	//	q.Set("r_",fmt.Sprintf("%.16f",rand.Float64()))
	//	return ClientHttp(con_.Scheme +"://"+con_.Host+"/uycyw/uymanage/tree/menutree.htm?"+q.Encode(),"GET",200,nil,func(body io.Reader)error{
	//		buf := bufio.NewReader(body)
	//		key := "varzNodes="
	//		keyLen := len(key)
	//		for{
	//			line,err := buf.ReadString('\n')
	//			if err != nil {
	//				if err == io.EOF {
	//					break
	//				}
	//				panic(err)
	//			}
	//			line = reg.ReplaceAllString(line,"")
	//			if len(line) <= keyLen{
	//				continue
	//			}
	//			//fmt.Println(line)
	//			if strings.Contains(line[:keyLen],key){
	//				db_:=[]map[string]interface{}{}
	//				line = line[keyLen:len(line)-1]
	//				err = json.Unmarshal([]byte(line),&db_)
	//				if err != nil {
	//					panic(err)
	//				}
	//				for _,d := range db_{
	//					if d["name"].(string) == "招聘求职"{
	//						ClassID = d["id"].(string)
	//						//return nil
	//					}
	//				}
	//				if ClassID == ""{
	//					panic("--")
	//				}
	//				return nil
	//			}
	//		}
	//		return nil
	//	})
	//})
}
func SaveSiteDB(title string,content string,dateTime int64) error {

	isPost:= false
	_viewKvDB(UserKvj,func(b *bolt.Bucket)error{
		v := b.Get([]byte(strings.ToLower(user_info.Get("username"))))
		if len(v)>0 {
			isPost = true
		}
		return nil
	})
	if !isPost {
		return fmt.Errorf("post err")
	}
	//if !Conf.Post {
	//	return fmt.Errorf("post err")
	//}


	//da :=time.Unix(dateTime,0).Add(time.Hour*24*7)
	//time.Now()
	title = strings.Join(regTitle.FindAllString(title,-1),"")
	_reg := regexp.MustCompile("招聘|求职")
	k := _reg.FindAllString(title, -1)
	if len(k)==0{
		title = StyleReg.ReplaceAllString(title,"招聘信息")
	}else{
		title = StyleReg.ReplaceAllString(title,"")
	}
	da :=time.Now().Add(time.Hour*24*7)
	db := map[string]string{
		"IMAGEPATH":"",
		"ClassID":ClassID,
		"USERTYPE":"1",
		"TITLE":title,
		"Content":content,
		"ENDTIME":da.Format("2006-01-02"),
		"id":"",
		"ID":"",
		"sw":"",
		"p":"",
		"UnitNo":"",
		"TEL":"",
		"EMAIL":"",
		"ADDRESS":""}
	//fmt.Println(db)
	//return fmt.Errorf("==")
	//db.Add("IMAGEPATH","")
	//db.Add("ClassID","3002090507")
	//db.Add("USERTYPE","1")
	////db.Add("TITLE",StyleReg.ReplaceAllString(title,""))
	//db.Add("TITLE",title)
	//db.Add("Content",content)
	//db.Add("ENDTIME",time.Unix(dateTime,0).Format("2006-01-02"))
	//db.Add("id","")
	//db.Add("ID","")
	//db.Add("sw","")
	//db.Add("p","")
	//db.Add("UnitNo","")
	//db.Add("TEL","")
	//db.Add("EMAIL","")
	//db.Add("ADDRESS","")
	//db.Add("ID","")

	he := Conf.Header
	he.Add("Content-Type","application/x-www-form-urlencoded")
	return ClientPost("http://jcpt.chengdu.gov.cn/uycyw/SupplyAndDemand/save.jsp","POST",200,db,he,func(body io.Reader)error{
		//fmt.Println(db)
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return err
		}
		if len(doc.Find("title").Text()) >0 {
			return nil
		}
		return fmt.Errorf("error")
	})

}
