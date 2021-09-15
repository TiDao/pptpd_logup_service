package main
import(
    "fmt"
    "io"
    "os"
    "bufio"
    "strings"
    "net/http"
    "log"
    //"os/exec"
)

func check(e error){
    if e != nil{
        log.Fatal(e)
    }
}

func judgement(filename string,username string) bool{
    file,err := os.Open(filename)
    check(err)
    buffer := bufio.NewReader(file)
    for {
        a,_,c := buffer.ReadLine()
        if  c == io.EOF{
            break
        }
        name := strings.Split(string(a)," ")
        if strings.Replace(name[0],"\"","",-1) == username {
            return false
        }
    }
    file.Close()
    return true
}

func WriteFile(filename string,content string){
        file,err := os.OpenFile(filename,os.O_APPEND|os.O_WRONLY,0666)
        check(err)

        write := bufio.NewWriter(file)
        _,err = write.WriteString(content)
        check(err)
        write.Flush()
        file.Close()
}


func service(writer http.ResponseWriter,r *http.Request){ //the url is http://IP:8000/create/get?username=xxx&password=xxx
    file := "/etc/openvpn/server/user/psw-file"
	//file := "./psw-file"
    url := r.URL.Query()
	UserName,ok := url["username"]
	if !ok{
		fmt.Fprintf(writer,"请求URL为http://192.168.1.25/create/get?username=xxx&password=xxx,xxx为替换内容\n")
		return
	}else if len(UserName) > 1{
		fmt.Fprintf(writer,"username只能有一个\n")
	}

	PassWord,ok := url["password"]
	if !ok{
		fmt.Fprintf(writer,"请求URL为http://192.168.1.25/create/get?username=xxx&password=xxx,xxx为替换内容\n")
		return
	}else if len(PassWord) > 1 {
		fmt.Fprintf(writer,"password只能有一个\n")
	}

    if !judgement(file,UserName[0]) {
        fmt.Fprintf(writer,"用户名已经被使用")
    }else{
        input := UserName[0]+" "+PassWord[0]+"\n"
        WriteFile(file,input)
        fmt.Fprintf(writer,"注册成功,%s",input)
		log.Printf("注册成功 %s",UserName[0])
    }
}

func showURL(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"请求URL为http://192.168.1.25/create/get?username=xxx&password=xxx,xxx为替换内容\n")
}



func main(){
    http.HandleFunc("/create/",service)
	http.HandleFunc("/",showURL)
	log.Println("start litening 0.0.0.0:80")
    log.Fatal(http.ListenAndServe("0.0.0.0:80",nil))
}
