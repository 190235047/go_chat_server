package logic

import(
    "fmt"
    "router"
    "encoding/json"
    "net"
    "time"
    "log"
    "sync"
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
	mutex sync.Mutex
}
//专门给userNode的room list
type roomList struct {
	roomid int64
	connList *userConnList
	next *roomList
	mutex sync.Mutex
}
//专门给roomNode的user list
type userList struct{
        uid int64
        uHash int16
        next *userList
	mutex sync.Mutex
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
	Msg string `json:"msg"`
	Username string `json:"username"`
	Uid int64 `json:"uid"`
}

const USER_BUCKET_NUM = 1024
const ROOM_BUCKET_NUM = 1024

//var roomArr = make([]roomNode, ROOM_BUCKET_NUM)
var roomArr [ROOM_BUCKET_NUM]*roomNode
//var userArr = make([]userNode, USER_BUCKET_NUM)
var  userArr [USER_BUCKET_NUM]*userNode
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
	user.roomList.mutex.Lock()
	defer user.roomList.mutex.Unlock()
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
		fmt.Println("connList equil nil")
		connList = new(userConnList)
		rList.connList = connList
		connList.conn = conn
		connList.time = time.Now().Unix()
		return true
	} 
	for connList != nil {
		if connList.conn == conn {
			connList.time = time.Now().Unix()
			//write error log
			fmt.Println("addConn conn equal, uid=", my.Uid)
			log.Println("[error] addConn conn equal, uid=" , my.Uid)
			return false
		}
		if (connList.next == nil) {
			connList.next = new(userConnList)
			connList = connList.next
			connList.conn = conn
			connList.time = time.Now().Unix()
			fmt.Println("addConn ok uid=", my.Uid)
			return true
		}
		connList = connList.next			
	}
	return false
}
func (this *userNode) closeConn(){
	if this == nil {
		return
	}
	fmt.Println("delete many...")
	roomList := this.roomList
	for roomList != nil {
		userConnList := roomList.connList
		for userConnList != nil {
			userConnList.conn.Close()
			userConnList = userConnList.next
		}
		roomList = roomList.next
	}
}
func (this RegisterJson) addUser(conn net.Conn){
	num := this.Uid % USER_BUCKET_NUM
	if (this.Uid == 0) {
		return
	}
	if userArr[num] == nil {
		userArr[num] = new(userNode)
	}
	//借助另外的实现队列
	userNodeHead := new(userNode)
	userNodeList := userNodeHead
	userNodeList.next = userArr[num]
	nowTime := time.Now().Unix()
	userNodeList.next.mutex.Lock()
	defer userNodeList.next.mutex.Unlock()
	for userNodeList.next != nil {
		
		//过期5分钟
		if userNodeList.next.time != 0 && nowTime - userNodeList.next.time > 300 {
			if userNodeList.next.next == nil {
				userNodeList.next.next = new(userNode)
			}
			fmt.Println("bbBBBBBBBBBBBBBBBBBBBB uid:", userNodeList.next.uid)
			userNodeList.next.closeConn()
			userNodeList.next = userNodeList.next.next
			continue
		}
		if userNodeList.next.uid == 0{
			fmt.Println("addUser userNodeList.uid = 0")
			userNodeList.next.time = time.Now().Unix()
			userNodeList.next.addConn(conn, this)
			break
		}
		if userNodeList.next.uid == this.Uid {
			fmt.Println("addUser userNodeList.uid != 0")
			userNodeList.next.time = time.Now().Unix()
			userNodeList.next.addConn(conn, this)
			break
		}
		if userNodeList.next.next == nil {
			userNodeList.next.next = new(userNode)
			userNodeList = userNodeList.next
		} else {
			userNodeList = userNodeList.next
		}
	}
	userArr[num] = userNodeHead.next
}

//判断userNode 里面room id 是否在里面
func (user *userNode) userIsInRoom(uid int64, roomid int64) bool {
	for user != nil {
		if user.uid == uid {
			roomList := user.roomList
			for roomList != nil {
				if roomList.roomid == roomid {
					fmt.Println("kkkkkkkkkkkkkk  roomid=", roomList.roomid)
					return true
				}
				fmt.Println("bbbbbbbbbbb  roomid=", roomList.roomid)
				roomList = roomList.next
			}
		}
		user = user.next
	}
	return false
}

func (room *roomNode) addUser(my RegisterJson) {
	num := my.Uid % USER_BUCKET_NUM
	user := userArr[num]
	//如果userNode 的roomList 有当前Roomid说明之前已经添加过当前用户的uid了
	if user != nil && user.userIsInRoom(my.Uid, my.Roomid) {
		fmt.Println("user is nil or user.userIsInRoom qqqqqqqqq")
		if room.userList == nil {
			fmt.Println("|||||||||||||||||||||||||")
		}
		//testGetRoomUserList(room.roomid)
		return
	}
	if room.userTail == nil {
		room.userTail = new(userList)
		room.userList = room.userTail
                room.userTail.uid = my.Uid
		fmt.Println("roomNode1  addUser new uid=", room.userList.uid)
                return
	}
	fmt.Println("roomNode2  addUser new uid=", my.Uid)
        newUser := new(userList)
        newUser.uid = my.Uid
        room.userTail.next = newUser
	room.userTail = newUser //???
	return
}

func testGetRoomUserList(roomid int64) {
	num  := roomid % ROOM_BUCKET_NUM
	room := roomArr[num]
	if (room == nil) {
		fmt.Println("no user list, roomid=", roomid)
	}
	for (room.userList != nil) {
		fmt.Println("testGetRoomUserList room userid ", room.userList.uid)
		//room.userList = room.userList.next
	}
}

func (this RegisterJson) addRoom(){
	num := this.Roomid % ROOM_BUCKET_NUM
	if (this.Roomid <= 0) {
		return
	}
	if roomArr[num] == nil {
		roomArr[num] = new(roomNode)
	}
	node := roomArr[num]
	for node != nil{
		if node.roomid == this.Roomid {
			fmt.Println("add Room node.roomid != 0")
			if node.userList == nil {
				fmt.Println("0000000yyyyyyyyyyybbbbbbbbbbb")
			}
			node.addUser(this)
			return;	
		}
		if node.roomid == 0 {
			fmt.Println("add Room node.roomid = 0")
			node.roomid = this.Roomid
			if node.userList == nil {
				fmt.Println("yyyyyyyyyyybbbbbbbbbbb")
			}
			node.addUser(this)
			if node.userList != nil {
			fmt.Println("kmsastttttttttttt uid=", node.userList.uid)
			}
			roomArr[num].userList = node.userList
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
	var jsonData RegisterJson
	if json.Unmarshal([]byte(this.Content), &jsonData) != nil {
		fmt.Println("Register json decode false")
		return
	}
	fmt.Printf("register : %s\n",this.Content)
        jsonData.addRoom() 
        //必须在addUser的上面 因为在addRoom里面加入uid得去user数组链表查看
        //这个user有没有包含这个房间 如果没有那么说明这个服务器上这个房间没有注册这个uid
	jsonData.addUser(this.Conn)
	if (roomArr[jsonData.Roomid%ROOM_BUCKET_NUM].userList == nil) {
		fmt.Println("{{{{{{{{{{{{{{{{ roomid=", jsonData.Roomid)
	}
}
func getRoomPtr(roomid int64) *roomNode{
	if (roomid <= 0) {
		return nil
	}
	num := roomid % ROOM_BUCKET_NUM
	node := roomArr[num]
	for node != nil{
		if node.roomid == roomid {
			fmt.Println("getRoomUserHead , roomid:", node.roomid)
			if (node.userList == nil) {
				fmt.Println("88888888888888888888")
			}
			return node
		}
		node = node.next
	}
	return nil
}
func getUserRoomListByUid(uid int64, roomid int64) *roomList{
	if (uid < 0) {
		return nil
	}
	fmt.Println("xxzzzzzzzzzzzzzzzzzzzzz uid:", uid)
	num := uid % USER_BUCKET_NUM
	node := userArr[num]
	for node != nil {
		if node.uid == uid {
			break
		}
		node = node.next
	}
	if node == nil {
		return nil
	}
	userRoomListNode := node.roomList
	for userRoomListNode != nil {
		if userRoomListNode.roomid == roomid {
			return userRoomListNode
		}
		userRoomListNode = userRoomListNode.next
	}
	return nil
}
func (this *User) SendMsg(){
	var jsonData SendMsgJson
        if json.Unmarshal([]byte(this.Content), &jsonData) != nil {
                fmt.Println("SendMsg json decode false")
                return
        }
	//应该把 this.Content 发送到kafka等队列然后用另外1个服务读取队列依次向其它前端机发送
	fmt.Println("send msg roomid : ", jsonData.Roomid)
	roomNodePtr := getRoomPtr(jsonData.Roomid)
	if roomNodePtr == nil || roomNodePtr.userList == nil {
		return
	}
	roomNodePtr.userList.mutex.Lock()
	defer roomNodePtr.userList.mutex.Unlock()
	userListHead := new(userList)
	userListForeach := userListHead
	userListForeach.next = roomNodePtr.userList
	for userListForeach.next != nil {
		fmt.Println("bbbbmmmmm")
		userConnListNode:= getUserRoomListByUid(userListForeach.next.uid, jsonData.Roomid)
		if userConnListNode == nil || userConnListNode.connList == nil {
			fmt.Println("kill kill kill uid:", userListForeach.next.uid)
			if userListForeach.next.next == nil {
				//说明是tail
				if userListForeach.next == roomNodePtr.userList {
					fmt.Println("userConnListNode == nil || userConnListNode.connList == nil, tail = nil")
					roomNodePtr.userTail = nil
				} else {
					//fmt.Println("userConnListNode == nil || userConnListNode.connList == nil, uid = ", userListForeach.next.next.uid)
					roomNodePtr.userTail = userListForeach.next
				}
			}
			//删除这个房间这个人的uid
			userListForeach.next = userListForeach.next.next
			continue
		}
		userConnListNode.mutex.Lock()
		defer userConnListNode.mutex.Unlock()
		//借助另外一个当header头来遍历，因为header有可能被删除
		userConnListHead := new(userConnList)
		foreachUserConnListNode := userConnListHead;
		foreachUserConnListNode.next = userConnListNode.connList;
		nowTime := time.Now().Unix()
		for foreachUserConnListNode.next != nil {
			//心跳时间大于5分钟算过期
			if nowTime - foreachUserConnListNode.next.time > 300 {
				fmt.Println("close QQQQQQQQ uid:", userListForeach.next.uid)
				foreachUserConnListNode.next.conn.Close()
				foreachUserConnListNode.next = foreachUserConnListNode.next.next
				continue
			}
			fmt.Println("send msg, uid:", userListForeach.next.uid)
			foreachUserConnListNode.next.conn.Write([]byte("send msg ok..."))
			foreachUserConnListNode = foreachUserConnListNode.next
		}
		userConnListNode.connList = userConnListHead.next
		userListForeach = userListForeach.next
	}
	roomNodePtr.userList = userListHead.next
}
func (this *User)Test(){

	fmt.Printf("sasas :")
}
