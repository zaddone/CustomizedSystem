package server
import(
	"io"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"time"
	"math/rand"
	//"bufio"
)

func LoginSite() error {
	err := ClientHttp("http://jcpt.chengdu.gov.cn/cdform/cdmanage/frameset/welcome.jsp","GET",200,nil,func(body io.Reader)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return err
		}
		var hr string
		doc.Find(".work ul li h3 a").EachWithBreak(func(i int, s *goquery.Selection)bool {
			if s.Text() =="内容发布" {
				hr,_ = s.Attr("href")
				err = inSite1(hr)
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
	cdform_,err := url.Parse(_url)
	if err != nil {
		return err
	}
	q_ := cdform_.Query()
	con_,err := url.Parse(q_.Get("url"))
	if err != nil {
		return err
	}
	q:= con_.Query()
	//q.Set("fn",q_.Get("fn"))
	q.Set("fn","jclongquanyiqu/longquanjiedao")
	return  ClientHttp(con_.Scheme +"://"+con_.Host+"/"+con_.Path+"?"+q.Encode(),"GET",200,nil,func(body io.Reader)error{
		rand.Seed(time.Now().UnixNano())
		q.Set("r_",fmt.Sprintf("%.16f",rand.Float64()))
		return ClientHttp(con_.Scheme +"://"+con_.Host+"/uycyw/uymanage/tree/menutree.htm?"+q.Encode(),"GET",200,nil,nil)
	})
}
func SaveSiteDB(title string,content string,dateTime int64) error {

	//da :=time.Unix(dateTime,0).Add(time.Hour*24*7)
	//time.Now()
	da :=time.Now().Add(time.Hour*24*7)
	db := map[string]string{
		"IMAGEPATH":"",
		"ClassID":"3002090507",
		"USERTYPE":"1",
		"TITLE":StyleReg.ReplaceAllString(title,""),
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
	return ClientPost("http://jcpt.chengdu.gov.cn/uycyw/SupplyAndDemand/save.jsp","POST",200,db,Conf.Header,func(body io.Reader)error{
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
