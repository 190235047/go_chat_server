package logic

import(
    "msgClient"
    "fmt"
)

type userNode struct{
        uid int64
        socket int
        next *userNode
}

type userList struct{
        uid int64
        uHash int16
        next *userList
}
type roomNode struct{
        roomid int64
        userList *userList
        next *roomNode
}

const USER_BUCKET_NUM = 1024
const ROOM_BUCKET_NUM = 1024

var roomArr = make([]roomNode, ROOM_BUCKET_NUM)
var userArr = make([]userNode, USER_BUCKET_NUM)

/*
type functionRegisterType func(*msgClient.Register) bool

var mapFunc = map[string]functionRegisterType {
                 "register" : register,
                }
*/
func Register(clientData *msgClient.Register) bool{
        fmt.Printf("register : %s\n",clientData.GetUsername())
        return true
}

