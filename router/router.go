package router

import(
        "fmt"
        "reflect"
)

type Router struct {
        method string
        action string
}

var logicMap = map[string]reflect.Value{}

func Register(action interface{}){
        fmt.Printf("action type :%s", reflect.TypeOf(action).Name())
        logicMap[reflect.TypeOf(action).Name()] = reflect.ValueOf(action)
        //logicMap["User"].MethodByName("Test").Call(nil) 
        // []reflect.Value{}
        fmt.Println(logicMap["User"])
}

func CallLogicFunc(modelName string, funcName string) {
        //fmt.Println(logicMap["User"])
        //logicMap["User"].MethodByName(funcName).Call([]reflect.Value{})
        callReflectType := logicMap[modelName].Type()
        call := reflect.New(callReflectType).Elem();
        call.MethodByName(funcName).Call([]reflect.Value{})
}
