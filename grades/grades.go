package grades

import (
	"fmt"
	"sync"
)

type Student struct{
	ID int
	FirstName string
	LastName string
	Grades []Grade
}
// 计算平均成绩
func (s Student)Average()float32{
	var result float32
	for _, grade := range s.Grades {
		result +=grade.Score
	}
	return result/float32(len(s.Grades))
}

type Students []Student

var (
	students Students
	studentsMutex sync.Mutex

)
func (ss Students)GetByID(id int)(*Student,error){
	for i:=range ss{
		if ss[i].ID==id{
			return &ss[i],nil
		}
	}
	return nil,fmt.Errorf("student with id: %d not found",id)
}
type GradeType string

const (
	GradeQuiz = GradeType("Quiz") // 小考
	GradeTest = GradeType("Test") // 单元测试
	GradeExam = GradeType("Exam") // 期末考试
)

type Grade struct{
	Title string
	Type GradeType
	Score float32
}
