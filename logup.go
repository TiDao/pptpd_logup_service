package main
import(
    "fmt"
    "io"
    "os"
    "bufio"
    "strings"
    "net/http"
    "log"
    "os/exec"
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

func service(writer http.ResponseWriter,r *http.Request){ //the url is http://IP:8000/username=XXXX_password=XXXXXX
    file := "/etc/ppp/chap-secrets"
    url := strings.Split(r.URL.Path,"_")
    user := strings.Split(url[0],"=")
    var Username string
    var Password string 
    if user[0] == "/username" {
        Username = user[1]
    }else{
        fmt.Fprintf(writer,"请输入正确的链接1,%s",user[0])
    }
    pass := strings.Split(url[1],"=")
    if pass[0] == "password" {
        Password = pass[1]
    }else{
        fmt.Fprintf(writer,"请输入正确的链接2,%s",pass[0])
    }
    if !judgement(file,Username) {
        fmt.Fprintf(writer,"用户名已经被使用")
    }else{
        input := "\""+Username+"\""+" "+"pptpd"+" "+"\""+Password+"\""+" "+"*"+"\n"
        WriteFile(file,input)
        exec.Command("bash","-c","service pptpd restart")
        fmt.Fprintf(writer,"注册成功,%s",input)
    }
}



func main(){
    http.HandleFunc("/",service)
    log.Fatal(http.ListenAndServe("0.0.0.0:8000",nil))
}
