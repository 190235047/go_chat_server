package main
import (
    "fmt"
    "net"
    "bytes"
    "encoding/binary"
    "github.com/golang/protobuf/proto"
    "msgClient"
)
func main() {
    conn, err := net.Dial("tcp", ":12345")
    if err != nil {
        panic(err)
    }
    test := &msgClient.Register{  // 使用辅助函数设置域的值
        Uid: proto.Int64(9999),
        Username: proto.String("Danger"),
        Roomid : proto.Int64(8888),
	Method : proto.String("register"),
    }    // 进行编码
    data, err := proto.Marshal(test)
    fmt.Println(len(data))
    //var num int16
    //conn.Write([]byte{0x7f, 0xff})
    //conn.Close()
	buf := new(bytes.Buffer)
	var len1 int
	len1 = len(data)
	err = binary.Write(buf, binary.LittleEndian, int16(len1))
	sendBuf := make([]byte, 2)
        sendBuf = buf.Bytes();
	sendBuf[0], sendBuf[1] = sendBuf[1], sendBuf[0];
	
	var endSendBuff []byte
	endSendBuff = append(endSendBuff, sendBuf...)
	endSendBuff = append(endSendBuff, data...)
	conn.Write(endSendBuff)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
fmt.Printf("% x", buf.Bytes())
    fmt.Printf("ok");
}
