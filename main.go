package main

import (
	"fmt"
	"net/http"

	"github.com/Tatascruz/api-students/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger()) // use the default RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	// Routes
	e.GET("/students", getStudents)
	e.POST("/students", createStudents)
	e.GET("/students/:id", getStudent)
	e.PUT("/students/:id", updateStudent)
	e.DELETE("/students/:id", deleteStudent)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler
func getStudents(c echo.Context) error {
	students, err := db.GetStudents()
	if err != nil {
		return c.String(http.StatusNotFound, "Failed to get students")
	}

	return c.JSON(http.StatusOK, students)
}

func createStudents(c echo.Context) error {
	student := db.Student{}
	if err := c.Bind(&student); err != nil {
		return err
	}
	if err := db.AddStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Error to create student")
	}

	return c.String(http.StatusOK, "Create student")
}

func getStudent(c echo.Context) error {
	id := c.Param("id")
	getStud := fmt.Sprintf("Get %s student", id)
	return c.String(http.StatusOK, getStud)
}

func updateStudent(c echo.Context) error {
	id := c.Param("id")
	UpdateStud := fmt.Sprintf("Update %s student", id)
	return c.String(http.StatusOK, UpdateStud)
}

func deleteStudent(c echo.Context) error {
	id := c.Param("id")
	deleteStud := fmt.Sprintf("Delete %s student", id)
	return c.String(http.StatusOK, deleteStud)
}
