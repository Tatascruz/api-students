package schemas

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name   string `json:"name"`
	CPF    string `json:"cpf" gorm:"uniqueIndex"`
	Email  string `json:"email" gorm:"uniqueIndex"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}

type StudentResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Active    bool      `json:"active"`
}

func NewStudentResponse(student Student) StudentResponse {
	return StudentResponse{
		ID:        int(student.ID),
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Name:      student.Name,
		CPF:       student.CPF,
		Email:     student.Email,
		Age:       student.Age,
		Active:    student.Active,
	}
}

func NewResponse(students []Student) []StudentResponse {
	studentsResponse := []StudentResponse{}

	for _, student := range students {
		studentsResponse = append(studentsResponse, NewStudentResponse(student))
	}
	return studentsResponse
}
