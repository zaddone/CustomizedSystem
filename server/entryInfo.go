package server
import(
	//"fmt"
	//"bufio"
	"io"
	"os"
	"path/filepath"
)
const (
	indexurl string = "http://jcpt.chengdu.gov.cn/cdform/index.jsp"
	inurl string = "http://jcpt.chengdu.gov.cn/cdform/cdmanage/login/login.htm"
	codeurl string = "http://jcpt.chengdu.gov.cn/cdform/cdmanage/include/image.jsp"
	saveurl string = "http://jcpt.chengdu.gov.cn/uycyw/SupplyAndDemand/save.jsp"
)
func UpdateCode() string{
	co := filepath.Join(Conf.Static,"img","code.tmp")
	err := ClientHttp(codeurl,"GET",200,nil,func(body io.Reader)error{
		fi,err := os.OpenFile(co,os.O_CREATE|os.O_RDWR|os.O_SYNC,0777)
		if err != nil {
			panic(err)
		}
		_,err = io.Copy(fi,body)
		if err != nil {
			panic(err)
		}
		fi.Close()
		//fmt.Println(n)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return co

}

