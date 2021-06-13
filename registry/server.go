package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
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
	err := r.sendRequiredServices(reg)

	// 服务变更通知
	r.notify(patch{
		Added: []patchEntry{
			{
				Name: reg.ServiceName,
				URL:  reg.ServiceUrl,
			},
		},
	})
	return err
}

func (r *registry) sendRequiredServices(reg Registration) error {
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
	return nil
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
			//
			r.notify(patch{
				Removed: []patchEntry{
					{
						Name: u.ServiceName,
						URL:  u.ServiceUrl,
					},
				},
			},
			)
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}

	}
	return fmt.Errorf("service at url %s not found", url)
}

//
func (r registry) notify(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}

				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}
// 心跳检测
func (r *registry)heartbeat(freq time.Duration)  {
	for {
		var wg sync.WaitGroup
		for _,reg :=range r.registrations{
			wg.Add(1)
			go func (reg Registration)  {
				defer wg.Done()
				success :=true
				// 重试
				for attemps :=0;attemps <3;attemps++{
					res,err := http.Get(reg.HeartbeatURL)
					if err != nil{
						log.Println(err)
					}else if res.StatusCode ==http.StatusOK {
						log.Printf("Heartbeat check passed for %v", reg.ServiceName)
						if !success{
							r.add(reg)
						}
						break;
						
					}
					log.Printf("Heartbeat check failed for %v", reg.ServiceName)
					if success{
						success = false
						r.remove(reg.ServiceUrl)
					}
					time.Sleep(time.Second*1)
				}
			}(reg)
			wg.Wait()
			time.Sleep(freq)
		}
	}
}

var once sync.Once
func SetupRegistryService()  {
	once.Do(func() {
		go reg.heartbeat(time.Second*3)
	})
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
