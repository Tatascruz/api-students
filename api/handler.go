package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tatascruz/api-students/schemas"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// getStudents godoc
//
// @Summary          Get a list of students
// @Description      Retrieve students details
// @Tags             students
// @Accpept          jsom
// @Produce          json
// @Param            register path    int false   "Registration"
// @Success          200 {object} schemas.StudentResponse
// @Failure          404
// @Router           /students [get]
func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		log.Error().Err(err).Msg("[api] error getting students")
		return c.String(http.StatusInternalServerError, "Failed to get students")
	}

	listOfStudents := map[string][]schemas.StudentResponse{"students": schemas.NewResponse(students)}

	return c.JSON(http.StatusOK, listOfStudents)
}

// createStudent godoc
//
// @Summary          Create student
// @Description      Create student
// @Tags             students
// @Accpept          jsom
// @Produce          json
// @Success          200 {object} schemas.StudentResponse
// @Failure          400
// @Router           /students [post]
func (api *API) createStudent(c echo.Context) error {
	studentReq := StudentRequest{}
	if err := c.Bind(&studentReq); err != nil {
		return err
	}

	if err := studentReq.Validate(); err != nil {
		log.Error().Err(err).Msgf("[api] error validating struct")
		return c.String(http.StatusBadRequest, "Error validating student")
	}

	student := schemas.Student{
		Name:   studentReq.Name,
		Email:  studentReq.Email,
		CPF:    studentReq.CPF,
		Age:    studentReq.Age,
		Active: *studentReq.Active,
	}

	if err := api.DB.AddStudent(&student); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.String(http.StatusConflict, "CPF or email already registered")
		}
		return c.String(http.StatusInternalServerError, "Error to create student")
	}
	return c.JSON(http.StatusCreated, schemas.NewStudentResponse(student))
}

// getStudents godoc
//
// @Summary          Get student by ID
// @Description      Retrieve students details using ID
// @Tags             students
// @Accpept          jsom
// @Produce          json
// @Success          200 {object} schemas.StudentResponse
// @Failure          404
// @Failure          500
// @Router           /students/{id} [get]
func (api *API) getStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student")
	}

	return c.JSON(http.StatusOK, schemas.NewStudentResponse(student))
}

// updateStudents godoc
//
// @Summary          Update Student
// @Description      Update Students details
// @Tags             students
// @Accpept          jsom
// @Produce          json
// @Success          200 {object} schemas.StudentResponse
// @Failure          404
// @Failure          500
// @Router           /students/{id} [put]
func (api *API) updateStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid student ID")
	}

	updatingStudent, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student")
	}

	var recivedStudent StudentRequest
	if err := c.Bind(&recivedStudent); err != nil {
		return c.String(http.StatusBadRequest, "Invalid body")
	}

	student := updateStudentInfo(recivedStudent, updatingStudent)

	if err := api.DB.UpdateStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to save student")
	}

	return c.JSON(http.StatusOK, schemas.NewStudentResponse(student))
}

// deleteStudents godoc
//
// @Summary          Delete Student
// @Description      Delete Students details
// @Tags             students
// @Accpept          jsom
// @Produce          json
// @Success          200 {object} schemas.StudentResponse
// @Failure          404
// @Failure          500
// @Router           /students/{id} [delete]
func (api *API) deleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid student ID")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student")
	}

	if err := api.DB.DeleteStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete student")
	}

	return c.JSON(http.StatusOK, schemas.NewStudentResponse(student))
}

func updateStudentInfo(recivedStudent StudentRequest, student schemas.Student) schemas.Student {
	if recivedStudent.Name != "" {
		student.Name = recivedStudent.Name
	}
	if recivedStudent.CPF != "" {
		student.CPF = recivedStudent.CPF
	}
	if recivedStudent.Email != "" {
		student.Email = recivedStudent.Email
	}
	if recivedStudent.Age > 0 {
		student.Age = recivedStudent.Age
	}
	if recivedStudent.Active != nil {
		student.Active = *recivedStudent.Active
	}
	return student
}
