package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// 注册服务的web service

//两个常量
const ServerPort = ":30000"

// 服务注册 web服务的地址，通过该地址可以查询到那些服务注册了
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 添加依赖
	r.sendRequiredServices(reg)
	return nil
}

func (r *registry) sendRequiredServices(reg Registration) {
	r.mutex.RLock() // 读锁
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceUrl,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {

	}
}

//
func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

//
func (r *registry) remove(url string) error {
	for i, u := range r.registrations {
		if u.ServiceUrl == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}

	}
	return fmt.Errorf("service at url %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with URL:%s\n", r.ServiceName, r.ServiceUrl)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodGet:
		b, err := json.Marshal(reg.registrations)
		if err != nil {
			log.Fatalln(err)
		}
		n, _ := w.Write(b)
		fmt.Println(n)
	case http.MethodDelete:
		pload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = reg.remove(string(pload))
		log.Printf("remove servicewith URL:%s\n", string(pload))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
