package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tatascruz/api-students/db"
	"github.com/Tatascruz/api-students/schemas"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestAPI cria uma API com um banco temporário, só para o teste.
func newTestAPI(t *testing.T) *API {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")
	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("erro ao abrir o banco de teste: %v", err)
	}
	if err := database.AutoMigrate(&schemas.Student{}); err != nil {
		t.Fatalf("erro na migração: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			sqlDB.Close()
		}
	})

	api := &API{
		Echo: echo.New(),
		DB:   db.NewStudentHandler(database),
	}
	api.ConfigureRoutes()
	return api
}

func TestGetStudentsEmpty(t *testing.T) {
	api := newTestAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/students", nil)
	rec := httptest.NewRecorder()
	api.Echo.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("esperava status 200, veio %d", rec.Code)
	}
}

// doRequest faz um pedido à API de teste e devolve a resposta.
func doRequest(api *API, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	api.Echo.ServeHTTP(rec, req)
	return rec
}

const validStudent = `{"name":"Ana","email":"ana@email.com","cpf":"01234567890","age":20,"active":true}`

func TestCreateStudent(t *testing.T) {
	api := newTestAPI(t)

	rec := doRequest(api, http.MethodPost, "/students", validStudent)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperava status 201, veio %d", rec.Code)
	}

	var resp schemas.StudentResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("a resposta não é um JSON válido: %v", err)
	}
	if resp.ID == 0 {
		t.Error("esperava um ID diferente de zero")
	}
	if resp.CPF != "01234567890" {
		t.Errorf("o CPF perdeu o zero à esquerda: veio %q", resp.CPF)
	}
}

func TestCreateStudentDuplicate(t *testing.T) {
	api := newTestAPI(t)

	first := doRequest(api, http.MethodPost, "/students", validStudent)
	if first.Code != http.StatusCreated {
		t.Fatalf("o primeiro cadastro deveria dar 201, veio %d", first.Code)
	}

	second := doRequest(api, http.MethodPost, "/students", validStudent)
	if second.Code != http.StatusConflict {
		t.Fatalf("esperava status 409 para aluno repetido, veio %d", second.Code)
	}
}

func TestGetStudentInvalidID(t *testing.T) {
	api := newTestAPI(t)

	rec := doRequest(api, http.MethodGet, "/students/abc", "")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperava status 400 para id inválido, veio %d", rec.Code)
	}
}
