package convert

import (
	"encoding/json"
	"log"
	"reflect"
	"unsafe"
)

/*
工具包
类型转换
*/

//MustMarshall 直接返回JSON序列化后的对象(节省error)
func MustMarshall(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	return bytes
}

//StringToByte 注意 - 转换完毕的对象不可更改
func StringToByte(src string) []byte {
	str := (*reflect.StringHeader)(unsafe.Pointer(&src))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: str.Data, Len: str.Len, Cap: str.Len}))
}

//ByteToString 注意 - 转换完毕的对象不可更改
func ByteToString(src []byte) string {
	str := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: str.Data, Len: str.Len}))
}
