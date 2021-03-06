package main
 
import (
    "flag"
    "fmt"
    "log"
    "os"
    "net"
    "runtime"
    "mynet"
    "encoding/binary"
    "github.com/golang/protobuf/proto"
    "bytes"
  _ "logic"
    "router"
    "myproto"
)
 
var (
    Port           = flag.String("i", ":12345", "IP port to listen on")
    logFileName    = flag.String("log", "cServer.log", "Log file name")
    //configFileName = flag.String("configfile", "config.ini", "General configuration file")
)
/*
var (
    configFile = flag.String("configfile", "config.ini", "General configuration file")
)
*/ 


func handleConn(conn net.Conn){
    log.Println("net conn")
    var headNum int
    var bodyNum int
    var x int16
    for {
            headBuff := make([]byte, 2)
            _, err := conn.Read(headBuff[headNum:])
            if (err != nil) {
                break
            }
            b_buf := bytes.NewBuffer(headBuff)
            binary.Read(b_buf, binary.BigEndian, &x)
            protoData := make([]byte, x);
            _, err =  conn.Read(protoData[bodyNum:])
            if (err != nil) {
                break
            }
            newData := &myproto.Client{}
            proto.Unmarshal(protoData, newData)
            //fmt.Printf("package length %d byte, mod:%s, action:%s\n", x, newData.GetModel(),newData.GetAction());
            flag := router.CallLogicFunc(newData.GetModel(),newData.GetAction(),newData.GetContent(), conn)
            if (!newData.GetIsKeep() || !flag) {
                break
            }
    }
	fmt.Println("hahha main close hahah")
    conn.Close()	
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    flag.Parse()
 
    //set logfile Stdout
    logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
    if logErr != nil {
        fmt.Println("Fail to find", *logFile, "cServer start Failed")
        os.Exit(1)
    }
    log.SetOutput(logFile)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    //set logfile Stdout End
    //start listen
    listenErr := mynet.StartListen(*Port, handleConn)
    if listenErr != nil {
        log.Fatalf("Server abort! Cause:%v \n", listenErr)
    }
}
