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
	//"time"
	//"math/rand"
	//"strings"
	//"log"
	//"fmt"
)
var (
	user_info = &url.Values{}
)

func loadRouter(){
	Router = gin.Default()
	Router.Static("/static","./static")
	Router.LoadHTMLGlob(Conf.Templates)
	Router.GET("/",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",gin.H{"user":user_info.Get("username")})
	})
	Router.GET("/codeimg",func(c *gin.Context){
		c.Header("Content-Type", "image/jpeg;charset=utf-8")
		var img []byte
		var buf [1024]byte
		err := ClientHttp(codeurl,"GET",200,nil,func(body io.Reader)error{
			//bufi := bufio.NewReader(body)
			for{
				n,err := body.Read(buf[0:])
				img= append(img,buf[:n]...)
				if err != nil {
					if err == io.EOF {
						return nil
					}else{
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
		c.Data(http.StatusOK,"image/jpeg;charset=utf-8",img)
	})

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

		err =  HandDBForBack(Conf.DbPath,func(db *sql.DB) error {
			row:= db.QueryRow("SELECT id FROM content WHERE id = ?",ids[0])
			var id_ int64
			err =  row.Scan(&id_)
			if err != nil {
				return err
			}
			err = SaveSiteDB(title,content,int64(dateTime))
			if err != nil {
				return err
			}
			sql_ := fmt.Sprintf("DELETE FROM content WHERE id in (%s) ",strings.Join(ids,","))
			_,err = db.Exec(sql_)
			if err != nil {
				panic(err)
			}

			return nil
		})
		if err != nil {
			c.JSON(http.StatusNotFound,gin.H{"msg":err})
			return
		}

		c.JSON(http.StatusOK,gin.H{"msg":"Success"})
		return


	})
	Router.POST("/verification",func(c *gin.Context){
		c.Header("Content-Type", "text/html; charset=utf-8")
		user_info.Set("username",c.DefaultPostForm("username",Conf.UserInfo.Get("username")))
		user_info.Set("password",c.DefaultPostForm("password",Conf.UserInfo.Get("password")))
		user_info.Set("randCode",c.PostForm("codeimg"))
		//fmt.Println(Conf.UserInfo.Encode())
		key := "window.location.href"
		read := "/cdform/cdmanage/vjsp/main.jsp"
		isOk := false
		err := ClientHttp(inurl,"POST",200,user_info,func(body io.Reader)error{
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
		//imgc := UpdateCode()
		c.HTML(http.StatusOK,
		"login.tmpl",
		gin.H{
		"username":Conf.UserInfo.Get("username"),
		"password":Conf.UserInfo.Get("password"),
		//"codeimg":strings.Replace(imgc,"\\","/",-1)})
		})
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
