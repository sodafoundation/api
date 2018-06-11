package env

import(
  "os"
  "bufio"
  "io"
  "strings"
  "os/exec"
  "fmt"
)
/* 
func main(){
   str :=GetGopath()
   fmt.Println(str)
}
*/
func GetGopath() string{
    cmd:=exec.Command("/bin/sh", "./path.sh")
    bytes, err := cmd.Output()
    if err != nil {
        fmt.Println("cmd.Output: ", err)
        return ""
    }
    fmt.Println(string(bytes))
    //exec.Command("/bin/sh", "./path.sh")
     file,err := os.Open("./gopath.log")
     if err !=nil{
       panic(err)
     }
    defer file.Close()
    rd := bufio.NewReader(file)

   for {
        line, err := rd.ReadString('\n')
        
        if err != nil || io.EOF == err {
            break
        }
       //remove space
        line = strings.Replace(line, " ", "", -1) 
       //remove enter
       line = strings.Replace(line, "\n", "", -1)
       return line
     }      
    return ""
}
func GetGoroot() string{
    exec.Command("/bin/sh", "./path.sh")
     file,err := os.Open("goroot.log")
     if err !=nil{
       panic(err)
     }
    defer file.Close()
    rd := bufio.NewReader(file)

   for {
        line, err := rd.ReadString('\n')
        
        if err != nil || io.EOF == err {
            break
        }
       //remove space
        line = strings.Replace(line, " ", "", -1) 
       //remove enter
       line = strings.Replace(line, "\n", "", -1)
       return line
     }      
    return ""
}
