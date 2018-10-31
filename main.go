package main
import(
	"fmt"
	"github.com/zaddone/CustomizedSystem/server"
	//"os/exec"
	//"runtime"
)
//var commands = map[string]string{
//	"windows": "explorer",
//	"darwin":  "open",
//	"linux":   "xdg-open",
//}
//func Open(uri string) error {
//	run,ok := commands[runtime.GOOS]
//	if !ok {
//		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
//	}
//	cmd := exec.Command(run,uri)
//	return cmd.Start()
//}
func main(){
	fmt.Println("customized system")

	//server.InSys()
	err :=server.Open("http://127.0.0.1"+server.Conf.Port+"/login")
	if err != nil {
		fmt.Println(err)
	}
	var cmd string
	for{
		fmt.Scanf("%s",&cmd)
		fmt.Println(cmd)
		cmd = ""
	}
}
