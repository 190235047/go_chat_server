package mynet
 
import (
    "errors"
    "fmt"
    "log"
    "net"
    //"sync"
    //"time"
    "encoding/binary"
    "bytes"
    "github.com/golang/protobuf/proto"
    "msgClient"
)
 
func handleConn(conn net.Conn) {
    defer conn.Close();
    headBuff := make([]byte, 2)
    var headNum int
    conn.Read(headBuff[headNum:])
    b_buf := bytes.NewBuffer(headBuff)  
    var x int16 
    binary.Read(b_buf, binary.BigEndian, &x)
    protoData := make([]byte, x);
    var bodyNum int
    conn.Read(protoData[bodyNum:])
    newData := &msgClient.Register{}
    proto.Unmarshal(protoData, newData)
    fmt.Printf("package length %d byte, name:%s, method:%s", x, newData.GetUsername(),newData.GetMethod());
}
 
//start listens
func StartListen(addr string) error {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return err
    }
    // if Errors accept arrive 100 .listener stop.
    for failures := 0; failures < 100; {
        conn, listenErr := listener.Accept()
        if listenErr != nil {
            log.Printf("number:%d,failed listening:%v\n", failures, listenErr)
            failures++
        }
        go handleConn(conn);
    }
    return errors.New("Too many listener.Accept() errors,listener stop")
}
