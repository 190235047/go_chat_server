package logic

import(
    "fmt"
    "router"
    "encoding/json"
    "net"
    "time"
    "log"
)

type User struct {
        router.Router
}

type userConnList struct{
	conn net.Conn
	time int64   //心跳更新时间
	next *userConnList
}

type userNode struct{
        uid int64
	time int64 //当前用户connList里最新的一次数据发送
        roomList *roomList
        next *userNode
}
//专门给userNode的room list
type roomList struct {
	roomid int64
	connList *userConnList
	next *roomList
}
//专门给roomNode的user list
type userList struct{
        uid int64
        uHash int16
        next *userList
}
type roomNode struct{
        roomid int64
        userList *userList
	userTail *userList
        next *roomNode
}

type RegisterJson struct{
	Roomid int64 `json:"roomid"`
	Uid    int64 `json:"uid"`
	Username string `json:"username"`
}

type SendMsgJson struct{
	Roomid int64 `json:"roomid"`
	Msg sring `json:"msg"`
	Username string `json:"username"`
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

func (user *userNode) addConn(conn net.Conn, my RegisterJson) bool{
	user.uid = my.Uid
	user.time = time.Now().Unix()
	if user.roomList == nil { //初始化
		user.roomList = new(roomList)
		user.roomList.roomid = my.Roomid
		fmt.Println("new")
	}
	rList := user.roomList
	var connList *userConnList
	for (rList != nil) {
		if rList.roomid == my.Roomid {
			connList = rList.connList
			fmt.Println("add Conn roomid=", my.Roomid)
			break;
		}
		if rList.next == nil {
			rList.next = new(roomList)
			rList.next.roomid = my.Roomid
			rList = rList.next
			//为什么不在这退出，让它再循环统一交给 roomList.roomid == my.Roomid 判断来退出
		}	
	}
	if connList == nil { //初始化
		connList = new(userConnList)
		rList.connList = connList
	} 
	for connList != nil {
		fmt.Println("uid : ", my.Uid)
		if connList.conn == conn {
			//write error log
			log.Println("[error] addConn conn equal, uid=" , my.Uid)
			return false
		}
		if (connList.next == nil) {
			connList.next = new(userConnList)
			connList = connList.next
			connList.conn = conn
			connList.time = time.Now().Unix()
			fmt.Println("addConn ok")
			return true
		}
		connList = connList.next			
	}
	return false
}

func (this RegisterJson) addUser(conn net.Conn){
	num := this.Uid % USER_BUCKET_NUM
	if (this.Uid == 0) {
		return
	}
	node := &userArr[num]
	for node != nil{
		if node.uid == 0{
			fmt.Println("addUser node.uid=0")
			node.addConn(conn, this)
			return
		}
		if node.uid == this.Uid {
			node.addConn(conn, this)
			return
		}
		if node.next == nil {
			node.next = new(userNode)
			node = node.next
		} else {
			node = node.next
		}
	}
}

//判断userNode 里面room id 是否在里面
func (user *userNode) userIsInRoom(roomid int64) bool {
	roomList := user.roomList
	for roomList != nil {
		if roomList.roomid == roomid {
			fmt.Println("kkkkkkkkkkkkkk  roomid=", roomList.roomid)
			return true
		}
		fmt.Println("bbbbbbbbbbb  roomid=", roomList.roomid)
		roomList = roomList.next
	}
	return false
}

func (room *roomNode) addUser(my RegisterJson) {
	num := my.Uid % USER_BUCKET_NUM
	user := &userArr[num]
	//如果userNode 的roomList 有当前Roomid说明之前已经添加过当前用户的uid了
	if !user.userIsInRoom(my.Roomid) {
		testGetRoomUserList(room.roomid)
		return
	}
	fmt.Println("ewrwer3455 addUser")
	if room.userTail == nil {
		room.userTail = new(userList)
		room.userList = room.userTail
                room.userTail.uid = my.Uid
                return
	}
	fmt.Println("roomNode  addUser new")
        newUser := new(userList)
        newUser.uid = my.Uid
        room.userTail.next = newUser
	return
}

func testGetRoomUserList(roomid int64) {
	num  := roomid % ROOM_BUCKET_NUM
	room := &roomArr[num]
	if (room == nil) {
		fmt.Println("no user list, roomid=", roomid)
	}
	for (room.userList != nil) {
		fmt.Println("testGetRoomUserList room userid ", room.userList.uid)
		room.userList = room.userList.next
	}
}

func (this RegisterJson) addRoom(){
	num := this.Roomid % ROOM_BUCKET_NUM
	if (this.Roomid <= 0) {
		return
	}
	node := &roomArr[num]
	for node != nil{
		if node.roomid == this.Roomid {
			node.addUser(this)
			return;	
		}
		if node.roomid == 0 {
			fmt.Println("add Room node.roomid=0")
			node.roomid = this.Roomid
			node.addUser(this)
			return
		}
		if node.next == nil {
			node.next = new(roomNode)
			node = node.next
			//为什么不在这里做return是因为让他进入下个循环在  node.roomid == 0 里面统一做 node.addUser(this.Uid) 的处理
		} else {
			node = node.next
		}
	}
}

func (this *User) Register() {

	//jsonData := make(map[string]interface{})
	var jsonData RegisterJson
	if json.Unmarshal([]byte(this.Content), &jsonData) != nil {
		fmt.Println("json decode false")
		return
	}
	fmt.Println("name", jsonData.Username)	
	fmt.Printf("register : %s\n",this.Content)
	f := []byte("asas")
	this.Conn.Write(f)
        jsonData.addRoom() 
        //必须在addUser的上面 因为在addRoom里面加入uid得去user数组链表查看
        //这个user有没有包含这个房间 如果没有那么说明这个服务器上这个房间没有注册这个uid
	jsonData.addUser(this.Conn)
}
func (this *User) SendMsg(){
	
}
func (this *User)Test(){

	fmt.Printf("sasas :")
}
