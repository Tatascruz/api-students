package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Tatascruz/api-students/shemas"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	listOfStudents := map[string][]shemas.StudentResponse{"students": shemas.NewResponse(students)}

	return c.JSON(http.StatusOK, listOfStudents)
}

func (api *API) createStudent(c echo.Context) error {
	studentReq := StudentRequest{}
	if err := c.Bind(&studentReq); err != nil {
		return err
	}

	if err := studentReq.Validate(); err != nil {
		log.Error().Err(err).Msgf("[api] error validating struct")
		return c.String(http.StatusBadRequest, "Error validating student")
	}

	student := shemas.Student{
		Name:   studentReq.Name,
		Email:  studentReq.Email,
		CPF:    studentReq.CPF,
		Age:    studentReq.Age,
		Active: *studentReq.Active,
	}

	if err := api.DB.AddStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Error to create student")
	}

	return c.JSON(http.StatusOK, student)
}

func (api *API) getStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student")
	}

	return c.JSON(http.StatusOK, student)
}

func (api *API) updateStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	updatingStudent, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	var recivedStudent shemas.Student
	if err := c.Bind(&recivedStudent); err != nil {
		return c.String(http.StatusBadRequest, "Invalid body")
	}

	student := updateStudentInfo(recivedStudent, updatingStudent)

	if err := api.DB.UpdateStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save student ID")
	}

	return c.JSON(http.StatusOK, student)
}

func (api *API) deleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	if err := api.DB.DeleteStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete student")
	}

	return c.JSON(http.StatusOK, student)

}

func updateStudentInfo(recivedStudent, student shemas.Student) shemas.Student {
	if recivedStudent.Name != "" {
		student.Name = recivedStudent.Name
	}
	if recivedStudent.CPF > 0 {
		student.CPF = recivedStudent.CPF
	}
	if recivedStudent.Email != "" {
		student.Email = recivedStudent.Email
	}
	if recivedStudent.Age > 0 {
		student.Age = recivedStudent.Age
	}
	if recivedStudent.Active != student.Active {
		student.Active = recivedStudent.Active
	}
	return student
}
