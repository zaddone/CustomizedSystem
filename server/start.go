package server
import(
	"database/sql"
	"net/http/cookiejar"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
	//"golang.org/x/net/publicsuffix"
	"github.com/boltdb/bolt"

	//"mime/multipart"
	//"bytes"

	"flag"
	"io"
	"net/http"
	"net/url"
	"compress/gzip"
	"compress/flate"
	"time"
	"os"
	"log"
	"fmt"
	//"path/filepath"
	"strings"
	"sync"
)

const(
	DBType string = "sqlite3"
)

var (
	DBMutex  sync.Mutex
	FileName   = flag.String("c", "conf.log", "config log")
	Router *gin.Engine
	Conf *Config
	Jar *cookiejar.Jar
	LastStr []byte = []byte("LastEntry")
)
type _row interface{
	Scan(dest ...interface{}) error
}
func init(){
	flag.Parse()
	Conf = NewConfig(*FileName)
	//KvDB, err := bolt.Open(KvDB, 0600, nil)
	//if err != nil {
	//	panic(err)
	//}
	_,err :=  os.Stat(Conf.DbPath)
	if err != nil {
		createDB()
	}
	//Jar,err := cookiejar.New(&cookiejar.Options{PublicSuffixList:publicsuffix.List})
	//TmpEntrys = make(chan *Entry,1000)
	Jar,_ = cookiejar.New(nil)
	LoadEntryChan()
	go loadRouter()




	//go runColl()
}

func createDB(){
	_sql :=`
	CREATE TABLE entry (
	id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	title	TEXT NOT NULL,
	url	TEXT,
	baseTime	INTEGER,
	beginTime	INTEGER NOT NULL,
	endTime		INTEGER NOT NULL,
	site INTEGER
	);
	CREATE TABLE content (
	id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	db	TEXT NOT NULL,
	tag	TEXT,
	entryid INTEGER  NOT NULL,
	style	INTEGER
	);
	CREATE TABLE style (
	id	INTEGER NOT NULL UNIQUE,
	name	TEXT NOT NULL UNIQUE,
	PRIMARY KEY(id)
	);
	`
	HandDB(Conf.DbPath,func(db *sql.DB){
		_,err := db.Exec(_sql)
		if err != nil {
			panic(err)
		}
		tx,err := db.Begin()
		if err != nil {
			panic(err)
		}
		stmt,err := tx.Prepare("INSERT INTO `style` (`id`,`name`) values (?,?)")
		if err != nil {
			panic(err)
		}
		_,err = stmt.Exec(0,"none")
		if err != nil {
			panic(err)
		}
		_,err = stmt.Exec(1,"title")
		if err != nil {
			panic(err)
		}
		_,err = stmt.Exec(2,"content")
		if err != nil {
			panic(err)
		}
		_,err = stmt.Exec(3,"contact")
		if err != nil {
			panic(err)
		}
		_,err = stmt.Exec(4,"addr")
		if err != nil {
			panic(err)
		}
		tx.Commit()
	})

}

func ClientPost(path string,ty string,statu int,body map[string]string,h http.Header, hand func(io.Reader)error) error {
	//r := &bytes.Buffer{}
	//if body != nil {
	//	writer := multipart.NewWriter(r)
	//	for k,v := range body {
	//		err := writer.WriteField(k,v)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//	writer.Close()
	//}
	var r io.Reader
	if body != nil {
		d := &url.Values{}
		for k,v := range body {
			d.Add(k,v)
		}
		r = strings.NewReader(d.Encode())
	}

	Req, err := http.NewRequest(ty,path,r)
	//Req.Close = true
	if err != nil {
		return err
	}
	for k,v := range h {
		for _,_v := range v{
			Req.Header.Add(k,_v)
		}
	}
	Req.Header.Add("Content-Type","application/x-www-form-urlencoded")
	Req.Header.Add("Content-Type","multipart/form-data")
	Cli := &http.Client{Jar:Jar}
	res, err := Cli.Do(Req)

	//var reqstr string
	//for{
	//	var req [1024]byte
	//	n,err := res.Request.Body.Read(req[0:])
	//	reqstr+=string(req[:n])
	//	if err != nil {
	//		fmt.Println(err)
	//		break
	//	}
	//}
	//fmt.Println(reqstr)

	if err != nil {
		return err
		//log.Println(err)
		//time.Sleep(time.Second*5)
		//return ClientHttp(path,ty,statu,body,hand)
	}
	if res.StatusCode != statu {
		var da [1024]byte
		n,err := res.Body.Read(da[0:])
		res.Body.Close()
		return fmt.Errorf("status code %d %s %s", res.StatusCode, path,string(da[:n]),err)
	}

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
		//defer reader.Close()
	default:
		reader = res.Body
	}
	if hand != nil {
		err = hand(reader)
	}
	reader.Close()
	return err

}
func ClientHttp_(path string,ty string,statu int,body *url.Values,h http.Header, hand func(io.Reader)error) error {
	var r io.Reader
	if body == nil {
		r = nil
	}else{
		r = strings.NewReader(body.Encode())
	}

	Req, err := http.NewRequest(ty,path,r)
	//Req.Close = true
	if err != nil {
		return err
	}
	for k,v := range h {
		for _,_v := range v{
			Req.Header.Add(k,_v)
		}
	}
	//Jar.SetCookies(u,[]*http.Cookie{&http.Cookie{}})
	//Req.Header = h
	//Req.AddCookie(&http.Cookie{})
	Cli := &http.Client{Jar:Jar}
	res, err := Cli.Do(Req)
	if err != nil {
		return err
		//log.Println(err)
		//time.Sleep(time.Second*5)
		//return ClientHttp(path,ty,statu,body,hand)
	}

	if res.StatusCode != statu {
		var da [1024]byte
		n,err := res.Body.Read(da[0:])
		res.Body.Close()
		return fmt.Errorf("status code %d %s %s", res.StatusCode, path,string(da[:n]),err)
	}
	//for _,co := range Req.Cookies() {
	//	fmt.Println("domain" ,co.Domain)
	//	fmt.Println("expires" ,co.Expires)
	//	fmt.Println("name" ,co.Name)
	//	fmt.Println("path" ,co.Path)
	//	fmt.Println("value",co.Value)
	//}
	//fmt.Println(res.Header.Get("Content-Encoding"))
	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
		//defer reader.Close()
	default:
		reader = res.Body
	}
	if hand != nil {
		err = hand(reader)
	}
	reader.Close()
	return err

}
func ClientHttp(path string,ty string,statu int,body *url.Values, hand func(io.Reader)error) error {
	return ClientHttp_(path,ty,statu,body,Conf.Header,hand)
}
func ClientDo(path string, hand func(io.Reader,*http.Response)error) error {

	Req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	Req.Header = Conf.Header
	var Client http.Client
	res, err := Client.Do(Req)
	if err != nil {
		log.Println(err)
		time.Sleep(time.Second*5)
		return ClientDo(path,hand)
	}
	if res.StatusCode != 200 {
		var da [1024]byte
		n,err := res.Body.Read(da[0:])
		res.Body.Close()
		return fmt.Errorf("status code %d %s %s", res.StatusCode, path,string(da[:n]),err)
	}
	//fmt.Println(res.Header.Get("Content-Encoding"))
	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
		//defer reader.Close()
	default:
		reader = res.Body
	}
	if hand != nil {
		err = hand(reader,res)
	}
	reader.Close()
	return err

}

func ViewKvDB(dbfile string,handle func(*bolt.Bucket)error) error {
	KvDB, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = KvDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("user"))
		if b== nil {
			return nil
		}
		return handle(b)
	})
	KvDB.Close()
	return err
}
func UpdateKvDB(dbfile string,handle func(*bolt.Bucket) error ) error {
	KvDB, err := bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	err = KvDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			return err
		}
		return handle(b)
	})
	KvDB.Close()
	return err
}

func HandDB(dbfile string,handle func(*sql.DB)){

	DBMutex.Lock()
	DB,err := sql.Open(DBType,dbfile)
	if err != nil {
		panic(err)
	}
	handle(DB)
	DB.Close()
	DBMutex.Unlock()

}

func HandDBForBack(dbfile string,handle func(*sql.DB) error) (err error){

	DB,err := sql.Open(DBType,dbfile)
	if err != nil {
		return err
	}
	err = handle(DB)
	DB.Close()
	return

}
func runColl () {
	for  {
		log.Println("run coll")
		go Collection()
		<-time.Tick(time.Hour)

	}
}
//func runEntry(){
//
//}
