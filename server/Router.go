package server
import(
	"github.com/gin-gonic/gin"
	//"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"bufio"
	"io"
	"fmt"
	"strings"
	"strconv"
	//"path/filepath"
	//"path"
	"database/sql"
	"time"
	//"math/rand"
	//"strings"
	//"log"
	//"fmt"
)

func loadRouter(){
	Router = gin.Default()
	Router.Static("/static","./static")
	Router.LoadHTMLGlob(Conf.Templates)
	Router.GET("/",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",nil)
	})
	//Router.GET("/codeimg",func(c *gin.Context){
	//	c.Header("Content-Type", "image/jpeg;charset=utf-8")
	//	img:=""
	//	var buf [1024]byte
	//	err := ClientHttp(codeurl,"GET",200,nil,func(body io.Reader)error{
	//		for{
	//			n,err := body.Read(buf[0:])
	//			img+= string(buf[:n])
	//			if err != nil {
	//				if err == io.EOF {
	//					return nil
	//				}else{
	//					return err
	//				}
	//			}
	//		}
	//		return nil
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	c.String(http.StatusOK,img,nil)
	//})

	Router.POST("/savesite",func(c *gin.Context){
		title := c.PostForm("title")
		if title == "" {
			c.JSON(http.StatusNotFound,"title == nil")
			return
		}
		content:= c.PostForm("content")
		if content == "" {
			c.JSON(http.StatusNotFound,"content == nil")
			return
		}
		dateTime,err := strconv.Atoi(c.PostForm("date"))
		if err != nil {
			c.JSON(http.StatusNotFound,err.Error())
			return
		}
		ids := c.PostFormArray("ids[]")
		if len(ids) == 0 {
			c.JSON(http.StatusNotFound,"ids is nil")
			return
		}
		err = SaveSiteDB(title,content,int64(dateTime))
		if err != nil {
			c.JSON(http.StatusNotFound,err.Error())
			return
		}

		err =  HandDBForBack(Conf.DbPath,func(db *sql.DB) error {
			sql_ := fmt.Sprintf("DELETE FROM content WHERE id in (%s) ",strings.Join(ids,","))
			_,err = db.Exec(sql_)
			return err
		})

		c.JSON(http.StatusOK,gin.H{"msg":"Success"})
		return


	})
	Router.GET("/savetest",func(c *gin.Context){
		err := SaveSiteDB("test___","点击蓝框查看",time.Now().Unix())
		if err != nil {
			c.JSON(http.StatusOK,err.Error())
		}else{
			c.JSON(http.StatusOK,"over")
		}
	})
	Router.GET("/savedb",func(c *gin.Context){
		db := &url.Values{}
		db.Add("IMAGEPATH","")
		db.Add("ClassID","3002090507")
		db.Add("USERTYPE","1")
		db.Add("TITLE",fmt.Sprintf("title_%d",time.Now().Unix()))
		db.Add("Content","content")
		db.Add("ENDTIME",time.Now().Format("2006-01-02"))
		db.Add("id","")
		db.Add("ID","")
		db.Add("sw","")
		db.Add("p","")
		db.Add("UnitNo","")
		db.Add("TEL","")
		db.Add("EMAIL","")
		db.Add("ADDRESS","")
		db.Add("ID","")
		fmt.Println(db.Encode())
		save:="http://jcpt.chengdu.gov.cn/uycyw/SupplyAndDemand/save.jsp"
		h := Conf.Header
		h.Add("Referer","http://jcpt.chengdu.gov.cn/uycyw/SupplyAndDemand/edit.jsp?ClassID=3002090507&sw=&id=")
		err := ClientHttp_(save,"POST",200,db,h,func(body io.Reader)error{
			buf := bufio.NewReader(body)
			for{
				line,err := buf.ReadString('\n')
				if err != nil {
					fmt.Println(err)
					if err == io.EOF {
						break
					}
					return err
				}
				fmt.Println(len(line),line)
			}
			return nil
		})
		if err != nil {
			c.JSON(http.StatusNotFound,err.Error())
			return
		}
	})
	Router.POST("/verification",func(c *gin.Context){

		c.Header("Content-Type", "text/html; charset=utf-8")
		Conf.UserInfo.Set("username",c.DefaultPostForm("username",Conf.UserInfo.Get("username")))
		Conf.UserInfo.Set("password",c.DefaultPostForm("password",Conf.UserInfo.Get("password")))
		Conf.UserInfo.Set("randCode",c.PostForm("codeimg"))
		fmt.Println(Conf.UserInfo.Encode())
		key := "window.location.href"
		read := "/cdform/cdmanage/vjsp/main.jsp?fn=jclongquanyiqu/longquanjiedao"
		isOk := false
		err := ClientHttp(inurl,"POST",200,Conf.UserInfo,func(body io.Reader)error{
			buf := bufio.NewReader(body)
			for{
				line,err := buf.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				if strings.Contains(line,key){
					if strings.Contains(line,read){
						isOk = true
					}
					break
				}
			}
			return nil

		})
		if err != nil {
			c.String(http.StatusOK,err.Error(),nil)
			return
		}

		if isOk {
			err = LoginSite()
			if err != nil {
				c.String(http.StatusOK,err.Error(),nil)
			}else{
				c.String(http.StatusOK,"<!DOCTYPE html><html><script language='javascript' type='text/javascript'>window.location.href='/';</script></html>",nil)
			}
		}else{
			c.String(http.StatusOK,"<!DOCTYPE html><html><script language='javascript' type='text/javascript'>window.location.href='/login';</script></html>",nil)
		}
		return
	})
	Router.GET("/login",func(c *gin.Context){
		err := ClientHttp("http://jcpt.chengdu.gov.cn/cdform/index.jsp","GET",200,nil,nil)
		if err != nil {
			panic(err)
		}
		imgc := UpdateCode()
		c.HTML(http.StatusOK,
		"login.tmpl",
		gin.H{
		"username":Conf.UserInfo.Get("username"),
		"password":Conf.UserInfo.Get("password"),
		"codeimg":strings.Replace(imgc,"\\","/",-1)})
	})
	Router.GET("/show",func(c *gin.Context){

		for{
			select{
			case en:= <-EntryList:
				en.Update()
				if en.CheckContent() {
					EntryList <- en
					c.JSON(http.StatusOK,en)
					return
				}
			default:
				c.JSON(http.StatusOK,nil)
				return
			}
		}
		return

	})
	Router.GET("/test",func(c *gin.Context){

		c.String(http.StatusOK,"",nil)
		return
	})
	Router.Run(Conf.Port)
}
