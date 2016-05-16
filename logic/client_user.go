package logic

import(
    "fmt"
    "router"
)

type User struct {
        router.Router
}

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

func init(){
        router.Register(User{})
}

func Register() bool{
        fmt.Printf("register : %s\n", "asdasd")
        return true
}

func (this *User) Register() {
        fmt.Printf("register : %s\n","sadddee")
}
func (this *User)Test(){

        fmt.Printf("sasas")
}
