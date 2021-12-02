package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func judgement(filename string, username string) bool {
	file, err := os.Open(filename)
	check(err)
	buffer := bufio.NewReader(file)
	for {
		line, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}
		name := strings.Split(line, " ")
		//if strings.Replace(name[0],"\"","",-1) == username {
		//    return false
		//}
		if name[0] == username {
			log.Printf("%s,%s", name[0], username)
			return false
		}
	}
	file.Close()
	return true
}

func AppendWriteFile(filename string, content string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	check(err)

	write := bufio.NewWriter(file)
	_, err = write.WriteString(content)
	check(err)
	write.Flush()
	file.Close()
}

func deleteWriteFile(filename string, username, password string) (uint8, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return 3, err
	}

	lines := strings.Split(string(fileContent), "\n")
	for i, line := range lines {
		lineSplit := strings.Split(line, " ")
		if lineSplit[0] == username {
			if lineSplit[1] == password {
				lines = append(lines[:i], lines[i+1:]...)
				output := strings.Join(lines, "\n")
				err = ioutil.WriteFile(filename, []byte(output), 0644)
				if err != nil {
					return 3, err
				}
				return 1, nil
			} else {
				return 2, nil
			}
		}
	}
	return 0, nil

}

func delete_service(writer http.ResponseWriter, r *http.Request) { //the url is http://IP:8000/create/get?username=xxx&password=xxx
	//file := "/etc/openvpn/server/user/psw-file"
	file := "./psw-file"
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request body error: %v\n", err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("json Decode error: %v\n", err)
		return
	} else {
		statusNumber, err := deleteWriteFile(file, user.UserName, user.PassWord)
		if err != nil {
			log.Printf("deleteWriteFile error: %v\n", err)
			return
		} else {
			switch statusNumber {
			case 0:
				{
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "用户名或密码不正确\n")
					log.Printf("用户名或密码不正确\n")
				}
			case 1:
				{
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "删除用户成功:%s\n", user.UserName)
					log.Printf("删除用户成功:%s\n", user.UserName)
				}
			case 2:
				{
					writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(writer, "用户名或密码不正确\n")
					log.Printf("用户名或密码不正确\n")
				}
			default:
				{
					log.Println("unknow error code")
				}

			}
		}
	}
}

func create_service(writer http.ResponseWriter, r *http.Request) { //the url is http://IP:8000/create/get?username=xxx&password=xxx
	//file := "/etc/openvpn/server/user/psw-file"
	file := "./psw-file"
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read body error: %v\n", err)
		fmt.Fprintf(writer, "read request body error: %v\n", err)
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("json Decode error: %v\n", err)
		fmt.Fprintf(writer, "json Decode error: %v\n", err)
		return
	}
	if user.UserName == "" || user.PassWord == "" {
		fmt.Fprintf(writer, "不允许使用空值作为用户名或密码\n")
		return
	}

	if !judgement(file, user.UserName) {
		fmt.Fprintf(writer, "用户名已经被使用,%s\n", user.UserName)
		log.Printf("用户名已经被使用,%s\n", user.UserName)
		return
	} else {
		input := user.UserName + " " + user.PassWord + "\n"
		AppendWriteFile(file, input)
		fmt.Fprintf(writer, "注册成功,%s\n", input)
		log.Printf("注册成功 %s\n", user.UserName)
	}
}

func download_service(writer http.ResponseWriter,r *http.Request) {
	//filename := "/root/openvpn_download/openvpn.zip"
	filename := "./openvpn.zip"
	file,err := os.Open(filename)
	if err != nil{
		log.Fatal("failed to open file: %v",err)
	}

	fileHeader := make([]byte,512)
	file.Read(fileHeader)
	fileStat,_ := file.Stat()

	writer.Header().Set("Content-Disposition", "attachment; filename=" + filename)
	writer.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	writer.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	file.Seek(0,0)
	io.Copy(writer,file)
}



func showURL(writer http.ResponseWriter, r *http.Request) {
	filename := "./index.html"
	html, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("failed to read html: %v\n", err)
	}

	fmt.Fprintf(writer, string(html))
}
