package router

import(
        //"fmt"
        "reflect"
	"net"
)

type Router struct {
        //model string
        //action string
	Content string
	Conn net.Conn
}

var logicMap = map[string]reflect.Value{}

func Register(model interface{}){
        //fmt.Printf("action type :%s", reflect.TypeOf(model).Name())
        logicMap[reflect.TypeOf(model).Name()] = reflect.ValueOf(model)
        //logicMap["User"].MethodByName("Test").Call([]reflect.Value{}) 
}

func CallLogicFunc(modelName string, funcName string, content string, conn net.Conn) bool{
        callReflect, ok := logicMap[modelName]
        if !ok {
                return false
        }
        callReflectType := callReflect.Type()
        call := reflect.New(callReflectType).Elem();
        call.FieldByName("Content").SetString(content)
	call.FieldByName("Conn").Set(reflect.ValueOf(conn))
	//fmt.Println("call=", call, ", address=", &call, ",type=", reflect.TypeOf(call), ",func=", funcName, ",fnObj=",call.MethodByName(funcName))
	//fmt.Printf("call addr : %d", &call)

	if call.CanAddr() && call.Kind() != reflect.Ptr{
		call = call.Addr()
	}
        //fmt.Printf("call addr : %d", &call)
	if fn := call.MethodByName(funcName); fn.IsValid(){
		fn.Call([]reflect.Value{})
                return true
	}
        return false;
}
