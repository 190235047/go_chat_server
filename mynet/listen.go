package mynet
 
import (
    "errors"
    //"fmt"
    "log"
    "net"
)

func handleClose(conn net.Conn) {
	conn.Close()
}

type callFuncType func(net.Conn)
var handleFuncMap = make(map[string]callFuncType)

func HandleFuc(funcName string, callFuncName callFuncType) {
	handleFuncMap[funcName] = callFuncName
}

func handleConn(conn net.Conn, callFuncName callFuncType) {
    defer conn.Close();
    callFuncName(conn)
}
 
//start listens
func StartListen(addr string, callFuncName callFuncType) error {
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
        go handleConn(conn, callFuncName);
    }
    return errors.New("Too many listener.Accept() errors,listener stop")
}
