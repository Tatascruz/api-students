package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Tatascruz/api-students/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, students)
}

func (api *API) createStudent(c echo.Context) error {
	student := db.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}
	if err := api.DB.AddStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Error to create student")
	}

	return c.String(http.StatusOK, "Create student")
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

	var recivedStudent db.Student
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
	id := c.Param("id")
	deleteStud := fmt.Sprintf("Delete %s student", id)
	return c.String(http.StatusOK, deleteStud)
}

func updateStudentInfo(recivedStudent, student db.Student) db.Student {
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
