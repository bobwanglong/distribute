package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// http handler
func RegisterHandlers(){
	handler :=new(studentsHandler)
	http.Handle("/students",handler) // 查看集合资源
	http.Handle("/students/",handler) //查看单个资源

}

type studentsHandler struct{}
// /students
// /students/{id}
// /students/{id}/grades
func (sh studentsHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	pathSegments :=strings.Split(r.URL.Path,"/")
	switch len(pathSegments){
	case 2:
		sh.getAll(w,r)
	case 3:
		id,err :=strconv.Atoi(pathSegments[2])
		if err !=nil{
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w,r,id)
	case 4:
		id,err :=strconv.Atoi(pathSegments[2])
		if err !=nil{
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w,r,id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh studentsHandler)getAll(w http.ResponseWriter, r *http.Request){
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data,err :=sh.toJson(students)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("content-type","application/json")
	w.Write(data)
}

func (sh studentsHandler)getOne(w http.ResponseWriter, r *http.Request,id int){
	studentsMutex.Lock()
	defer 	studentsMutex.Unlock()
	// for _, student := range students {
	// 	if student.ID == id{
	// 		data,err :=sh.toJson(student)
	// 		if err !=nil{
	// 			w.WriteHeader(http.StatusInternalServerError)
	// 			log.Println(err)
	// 			return
	// 		}
	// 		w.Header().Add("content-type","application/json")
	// 		w.Write(data)
	// 	}
	// }
	 student,err :=students.GetByID(id)
	 if err != nil{
		 w.WriteHeader(http.StatusNotFound)
		 log.Println(err)
		 return
	 }
	 data,err :=sh.toJson(student)
	 if err != nil{
		 w.WriteHeader(http.StatusInternalServerError)
		 log.Println(err)
		 return
	 }
	 w.Header().Add("content-type","application/json")
	 w.Write(data)
}

func (sh studentsHandler)addGrade(w http.ResponseWriter, r *http.Request,id int){
	studentsMutex.Lock()
	defer 	studentsMutex.Unlock()
	student,err :=students.GetByID(id)
	 if err != nil{
		 w.WriteHeader(http.StatusNotFound)
		 log.Println(err)
		 return
	 }
	 var g Grade
	 dec := json.NewDecoder(r.Body)
	 err =dec.Decode(&g)
	 if err !=nil{
		 w.WriteHeader(http.StatusBadRequest)
		 log.Println(err)
		 return
	 }
	 student.Grades = append(student.Grades, g)
	 w.WriteHeader(http.StatusCreated) //201
	 data,err :=sh.toJson(student)
	 if err != nil{
		 w.WriteHeader(http.StatusInternalServerError)
		 log.Println(err)
		 return
	 }
	 w.Header().Add("content-type","application/json")
	 w.Write(data)
}

func (sh studentsHandler)toJson(obj interface{})([]byte,error){
	var b bytes.Buffer
	enc :=json.NewEncoder(&b) // 编码器
	err :=enc.Encode(obj)
	if err != nil{
		return nil,fmt.Errorf("Failed to serialize students: %q",err)
	}
	return b.Bytes(),nil
}