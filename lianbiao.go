package main

import (
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

func main(){
	roomArr[0].roomid = 34234234
	fmt.Printf("sasa %d", roomArr[0].roomid)
	if roomArr[0].next == nil && roomArr[0].userList == nil {
		fmt.Printf("asdadasd")
	}
}
