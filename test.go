package main
import(
	"regexp"
	"fmt"
)
func main(){
	var StyleReg *regexp.Regexp = regexp.MustCompile("\\d{1,2}-\\d{1,2}$")
	var reg *regexp.Regexp = regexp.MustCompile(`[^\p{Han}^\w^\pP]+`)
	title := reg.ReplaceAllString("龙泉驿区柏合镇广盛通通讯器材经营部 10-31","")
	title = StyleReg.ReplaceAllString(title,"sss")
	fmt.Println(title)
}
