package shemas

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name   string `json:"name"`
	CPF    int    `json:"cpf"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}

type StudentResponse struct {
	ID        int       `json:"id"`
	CreateAt  time.Time `json:"createAt"`
	UpdateAt  time.Time `json:"updateAt"`
	DeletedAt time.Time `json:"deleteAt"`
	Name      string    `json:"name"`
	CPF       int       `json:"cpf"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
}

func NewResponse(students []Student) []StudentResponse {
	studentsResponse := []StudentResponse{}

	for _, student := range students {
		studentResponse := StudentResponse{
			ID:       int(student.ID),
			CreateAt: student.CreatedAt,
			UpdateAt: student.UpdatedAt,
			Name:     student.Name,
			CPF:      student.CPF,
			Email:    student.Email,
			Age:      student.Age,
			Active:   student.Active,
		}
		studentsResponse = append(studentsResponse, studentResponse)
	}
	return studentsResponse
}
