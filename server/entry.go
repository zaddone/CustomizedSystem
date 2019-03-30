package server
import(
	"reflect"
	"database/sql"
	"fmt"
	"io"
	"strings"
	"github.com/PuerkitoBio/goquery"
	//"bufio"
	"github.com/boltdb/bolt"
	"net/http"
	"time"
	//"sync"
)
var (
	EntryList chan *Entry = make(chan *Entry,1000)
	//EntrySync  sync.Mutex
)
type PostDB struct {
	title string
	content string
	dateTime int64
}
type Content struct {
	Id int64
	Db string `json:"db"`
	Tag string `json:"tag"`
	Entryid int64 `json:"entryid"`
	Style int64 `json:"style"`
}
func (self *Content) LoadDB_(row _row) (err error) {
	return row.Scan(
		&self.Id,
		&self.Db,
		&self.Tag,
		&self.Entryid,
		&self.Style)
}

type Entry struct {
	Id	 int64
	Title	string `json:"title"`
	Url string `json:"url"`
	BaseTime int64 `json:"baseTime"`
	BeginTime int64 `json:"beginTime"`
	EndTime int64 `json:"endTime"`
	Site	int64 `json:"site"`
	Con []*Content
}

func (self *Entry) DelDB_(db *sql.DB)(err error){
	_,err = db.Exec("DELETE FROM content WHERE entryid = ? ",self.Id)
	if err != nil {
		return
		//panic(err)
	}
	_,err = db.Exec("DELETE FROM entry WHERE id = ?",self.Id)
	return
}
func (self *Entry) DelDB()(err error){

	HandDB(Conf.DbPath,func(db *sql.DB){
		err = self.DelDB_(db)
	})
	return

}
func (self *Entry) CheckContent() bool {

	for _,c := range self.Con {
		if c.Style > 0 {
			return true
		}
	}
	return false

}

func (self *Entry) SaveDB(){

	HandDB(Conf.DbPath,func(db *sql.DB){
		res := StructSaveForDB(db,"entry",self)
		id,err := res.LastInsertId()
		if err != nil {
			panic(err)
		}
		self.Id = int64(id)
		tx,err := db.Begin()
		if err != nil {
			panic(err)
		}
		stmt,err := tx.Prepare("INSERT INTO `content` (`db`,`tag`,`entryid`,`style`) values (?,?,?,?)")
		if err != nil {
			panic(err)
		}
		for _,co := range self.Con {
			//co.Entryid = self.Id
			_,err =stmt.Exec(co.Db,co.Tag,self.Id,co.Style)
			if err != nil {
				panic(err)
			}
		}
		err = tx.Commit()
		if err != nil {
			panic(err)
		}
	})
	err := UpdateKvDB(Conf.DeduPath,func(b *bolt.Bucket)error{
		b.Put(LastStr,[]byte(self.Title))
		return nil
	})
	if err != nil {
		panic(err)
	}
	//EntryList <- self
}
func ReadAllEntrys(db *sql.DB,hand func(*Entry)) {
	d, _ := time.ParseDuration("-24h")
	rows,err := db.Query("SELECT id,title,url,baseTime,beginTime,endTime,site FROM entry WHERE baseTime > ? ORDER BY id DESC",time.Now().Add(d*7).Unix())
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		en := &Entry{}
		err = en.LoadDB_(rows,db)
		if err != nil {
			panic(err)
		}
		hand(en)
		//TmpEntrys <- en
	}
	rows.Close()
}

func LoadEntryChan() {
	//var err error
	HandDB(Conf.DbPath,func(db *sql.DB){
		var tmpen []*Entry
		ReadAllEntrys(db,func(en *Entry){
			//fmt.Println(en)
			if !en.CheckContent() {
				tmpen = append(tmpen,en)
			}else{
				EntryList <- en
			}
		})
		for _,en := range tmpen {
			en.DelDB_(db)
		}
	})
}

func GetLastEntry(db *sql.DB) (en *Entry) {
	row := db.QueryRow("SELECT id,title,url,baseTime,beginTime,endTime,site FROM entry order by id desc limit 1")
	en = &Entry{}
	err := row.Scan(
		&en.Id,
		&en.Title,
		&en.Url,
		&en.BaseTime,
		&en.BeginTime,
		&en.EndTime,
		&en.Site)
	if err != nil {
		//panic(err)
		fmt.Println(err)
		return nil
	}
	return en
}

type info struct {
	title string
	content string
}

func (self *Entry) HandContent(){

	fmt.Println(self)
	err := ClientDo(self.Url,func(body io.Reader,res *http.Response)error{
		doc,err := goquery.NewDocumentFromReader(body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var styleId int64 = 0
		doc.Find("#page-content p").EachWithBreak(func(i int, s *goquery.Selection) bool {
			text :=s.Text()
			str := reg.ReplaceAllString(text,"")
			str = strings.TrimSpace(str)
			if str == "" {
				styleId = 1
				return true
			}

			//switch styleId {
			//case 1:
			//	styleId = 2
			//default:
			//	fi := StyleReg.FindAllString(str,-1)
			//	if len(fi)==1 {
			//		styleId = 1
			//	}
			//}
			if str == "关注" {
				return false
			}
			var tag []string
			for _,t := range strings.Split(tagReg.ReplaceAllString(text," ")," ") {
				if t !="" {
					tag = append(tag,t)
				}
			}
			self.Con = append(self.Con,&Content{
				Style:styleId,
				Db:str,
				Tag:strings.Join(tag," ")})
			if styleId == 1 {
				styleId ++
			}
			//fmt.Println(tag)
			return true

		})

		if len(self.Con) >0 {
			self.SaveDB()
			fmt.Println(len(self.Con))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
func (self *Entry) Update() {

	HandDB(Conf.DbPath,func(db *sql.DB){
		row := db.QueryRow("SELECT id,title,url,baseTime,beginTime,endTime,site FROM entry WHERE id = ?",self.Id)
		err := self.LoadDB_(row,db)
		if err != nil {
			panic(err)
		}
	})
	return

}
func (self *Entry) LoadDB_(row _row,db *sql.DB) (err error) {
	if row == nil {
		return fmt.Errorf("row == nil")
	}
	err = row.Scan(
		&self.Id,
		&self.Title,
		&self.Url,
		&self.BaseTime,
		&self.BeginTime,
		&self.EndTime,
		&self.Site)
	if err != nil {
		return err
	}
	rows,err := db.Query("SELECT id,db,tag,entryid,style FROM content WHERE entryid = ?",self.Id)
	if err != nil {
		return err
	}
	self.Con = nil
	for rows.Next(){
		co:=&Content{}
		err = co.LoadDB_(rows)
		if err != nil {
			panic(err)
		}
		self.Con = append(self.Con,co)
	}
	rows.Close()

	return nil
}
func (self *Entry) LoadDB(id int64,db *sql.DB) (err error) {

	row := db.QueryRow("SELECT id,title,url,baseTime,beginTime,endTime,site FROM entry WHERE id = ?",id)
	err = row.Scan(
		&self.Id,
		&self.Title,
		&self.Url,
		&self.BaseTime,
		&self.BeginTime,
		&self.EndTime,
		&self.Site)
	if err != nil {
		return err
	}
	return nil

}

func ReadEntry(hand func(*Entry) error,where string,val ...interface{}) error {
	sql_ := "SELECT id,title,url,baseTime,beginTime,endTime,site FROM entry " + where + ";"
	var en *Entry
	return HandDBForBack(Conf.DbPath,func(db *sql.DB)error{
		row,err := db.Query(sql_,val...)
		if err != nil {
			return err
		}
		for row.Next() {
			en = &Entry{}
			err = row.Scan(
				&en.Id,
				&en.Title,
				&en.Url,
				&en.BaseTime,
				&en.BeginTime,
				&en.EndTime,
				&en.Site)
			if err != nil {
				return err
			}
			err = hand(en)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func StructSaveForDB(db *sql.DB,TableName string,st interface{}) sql.Result {
	re := reflect.TypeOf(st).Elem()
	va := reflect.ValueOf(st).Elem()
	var fi []string
	var x []string
	var val []interface{}
	for i:=0;i<re.NumField();i++ {
		str :=re.Field(i).Tag.Get("json")
		if str != "" {
			v:= va.Field(i).Interface()
			if v  != nil{
				x = append(x,"?")
				fi = append(fi,str)
				val = append(val,v)
			}
		}
	}
	sql_ := fmt.Sprintf("INSERT INTO %s (%s) values (%s)",TableName,strings.Join(fi,","),strings.Join(x,","))
	//fmt.Println(sql_)
	res,err := db.Exec(sql_,val...)
	if err != nil {
		panic(err)
	}
	return res
	//__id,__err :=m.RowsAffected()

	//fmt.Println(_id,_err,__id,__err)
}

func StructUpdateForDB(db *sql.DB,TableName string,st interface{},keyname string ,key interface{}) sql.Result {
	re := reflect.TypeOf(st).Elem()
	va := reflect.ValueOf(st).Elem()
	var fi []string
	var val []interface{}
	for i:=0;i<re.NumField();i++ {
		str :=re.Field(i).Tag.Get("json")
		if str != "" {
			v:= va.Field(i).Interface()
			if v  != nil{
				fi = append(fi,str+" = ?")
				val = append(val,v)
			}
		}
	}
	val = append(val,key)
	sql_ := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",TableName,strings.Join(fi,","),keyname)
	if db == nil {
		fmt.Println(sql_)

		panic("db == nil")
	}
	res,err := db.Exec(sql_,val...)
	if err != nil {
		panic(err)
	}
	return res

}
