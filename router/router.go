package router

import(
        "fmt"
        "reflect"
)

type Router struct {
        //model string
        //action string
	//Content string
}

var logicMap = map[string]reflect.Value{}

func Register(model interface{}){
        //fmt.Printf("action type :%s", reflect.TypeOf(model).Name())
        logicMap[reflect.TypeOf(model).Name()] = reflect.ValueOf(model)
        //logicMap["User"].MethodByName("Test").Call([]reflect.Value{}) 
        //fmt.Println(logicMap["User"])
}

func CallLogicFunc(modelName string, funcName string, content string) {
        //fmt.Println(logicMap)
        callReflectType := logicMap[modelName].Type()
        call := reflect.New(callReflectType).Elem();
	//call.FieldByName("Content").SetString(content)
	fmt.Println("call=", call, ", address=", &call, ",type=", reflect.TypeOf(call), ",func=", funcName, ",fnObj=",call.MethodByName(funcName))
	fmt.Printf("call addr : %d", &call)
//	call.MethodByName(funcName).Call([]reflect.Value{})

	if call.CanAddr() && call.Kind() != reflect.Ptr{
		call = call.Addr()
	}

	if fn := call.MethodByName(funcName); fn.IsValid(){
		fn.Call([]reflect.Value{})
	}
}
