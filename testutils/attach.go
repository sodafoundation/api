package main

import(
  "os"
  "strconv"
  "attach"
  "fmt"
)

func main(){
  conn := make(map[string]interface{})
  conn["attachment"]= os.Args[1]
  conn["discard"],_=strconv.ParseBool(os.Args[2])
  conn["targetDiscovered"],_=strconv.ParseBool(os.Args[3])
  conn["targetIQN"]=os.Args[4]
  conn["targetLun"],_=strconv.ParseFloat(os.Args[5],64)
  conn["targetPortal"]=os.Args[6]
  fmt.Println("conn main:",conn)
  var isc = &attach.Iscsi{}
  isc.Attach(conn)
}
