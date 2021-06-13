package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration)error{
	//
	serviceUpdateUrl,err :=url.Parse(r.ServiceUpdateURL)
	if err !=nil{
		return err
	}
	// 绑定给方法
	http.Handle(serviceUpdateUrl.Path,&serviceUpdateHanlder{})
	buf :=new(bytes.Buffer)
	enc :=json.NewEncoder(buf) // 编码
	err = enc.Encode(r)
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
type serviceUpdateHanlder struct{

}
func (suh *serviceUpdateHanlder)ServeHTTP(w http.ResponseWriter,r *http.Request)  {
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec :=json.NewDecoder(r.Body)
	var p patch
	err :=dec.Decode(&p)
	if err !=nil{
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}
func ShutdownService(url string)error{
	req,err :=http.NewRequest(http.MethodDelete,ServicesURL,bytes.NewBuffer([]byte(url)))
	if err !=nil{
		return err
	}
	req.Header.Add("content-type","text/plain")
	res,err:=http.DefaultClient.Do(req)
	if err!=nil{
		return err
	}
	if res.StatusCode != http.StatusOK{
		return fmt.Errorf("Failed to deregister service. registry service responded with code %v",res.StatusCode)
	}
	return nil
}

// 服务提供方
type providers struct{
	services map[ServiceName][]string // "服务名"：“URL”
	mutex *sync.RWMutex
}
//
func (p *providers)Update(pat patch)  {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// 增加的情况
	for _,patchEntry :=range pat.Added{
		if _,ok :=p.services[patchEntry.Name];!ok{
			p.services[patchEntry.Name]=make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}
	// 减少的情况
	for _,patchEntry :=range pat.Removed{
		if providersURLs,ok :=p.services[patchEntry.Name];ok{
			for i :=range providersURLs{
				if providersURLs[i]==patchEntry.URL{
					p.services[patchEntry.Name] = append(providersURLs[:i],providersURLs[i+1:]... )
				}
			}
		}
	}
}

func (p providers)get(name ServiceName)(string,error){
	providers,ok :=p.services[name]
	if !ok{
		return "",fmt.Errorf("No providers available for service %v", name)	
	}
	idx :=int(rand.Float32()*float32(len(providers)))
	return providers[idx],nil
}
// 对外暴露功能
func GetProvider(name ServiceName)(string,error){
return prov.get(name)
}
var prov =providers{
	services: make(map[ServiceName][]string),
	mutex:new(sync.RWMutex),
}