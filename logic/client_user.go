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
	tail *userList
        next *roomNode
}

type RegisterJson struct{
	Roomid int64 `json:"roomid"`
	Uid    int64 `json:"uid"`
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
	roomList := user.roomList
	for (roomList != nil) {
		if roomList.roomid == my.Roomid {
			connlist := roomList.connList
			break;
		}
		if roomList.next == nil {
			roomList.next = new(roomList)
			roomList.next.roomid = my.Roomid
			roomList = roomList.next
			//为什么不在这退出，让它再循环统一交给 roomList.roomid == my.Roomid 判断来退出
		}	
	}
	if connList == nil { //初始化
		connList = new(userConnList)
	} 
	for connList != nil {
		fmt.Println("uid : ", my.Uid)
		if connlist.conn == conn {
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
	num = 3  //test test delete
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
			return true
		}
		roomList = roomList.next
	}
	return false
}

func (room *roomNode) addUser(uid int64, my RegisterJson) {
	num := uid % USER_BUCKET_NUM
	user := &userArr[num]
	//如果userNode 的roomList 有当前Roomid说明之前已经添加过当前用户的uid了
	if !user.userIsInRoom(my.Roomid) {
		return
	}
	if room.tail == nil {
		room.tail = new(userList)
		room.userList = room.tail
	}
	
	
}

func (this *RegisterJson) addRoom(){
	num := this.Roomid % ROOM_BUCKET_NUM
	if (this.Roomid <= 0) {
		return
	}
	node := &roomArr[num]
	for node != nil{
		if node.roomid == this.Roomid {
			return;	
		}
		if node.roomid == 0 {
			fmt.Println("add Room node.roomid=0")
			node.addUser(this.Uid)
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
	jsonData.addUser(this.Conn)
}
func (this *User)Test(){

	fmt.Printf("sasas :")
}
