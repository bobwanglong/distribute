package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration)error{
	buf :=new(bytes.Buffer)
	enc :=json.NewEncoder(buf) // 编码
	err := enc.Encode(r)
	if err !=nil{
		return err
	}
	// 发送被注册服务到服务注册web服务器
	res,err :=http.Post(ServicesURL,"application/json",buf)
	if err !=nil{
		return err
	}
	if res.StatusCode != http.StatusOK{
		return fmt.Errorf("Failed to register service"+"responded with code %v",res.StatusCode)
	}
	return nil
}