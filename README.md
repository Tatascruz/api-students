
# api-students

API REST de cadastro de alunos, feita em Go como projeto de estudo (curso "Golang do Zero"), com documentação interativa em Swagger.

## Tecnologias

- Go
- [Echo](https://echo.labstack.com/) (framework web)
- [GORM](https://gorm.io/) + SQLite (banco de dados)
- zerolog (logs)
- Swagger (swaggo/echo-swagger)

## Como rodar

```bash
git clone https://github.com/Tatascruz/api-students.git
cd api-students
go run main.go
```

- A API sobe em `http://localhost:8080`.
- O banco `student.db` é criado automaticamente na primeira execução.
- Documentação Swagger: `http://localhost:8080/swagger/index.html`

## Rotas

| Método | Rota | Descrição | Respostas |
|--------|------|-----------|-----------|
| GET | `/students` | Lista todos os alunos | 200, 500 |
| POST | `/students` | Cadastra um aluno | 201, 400, 409, 500 |
| GET | `/students/:id` | Busca um aluno pelo id | 200, 400, 404, 500 |
| PUT | `/students/:id` | Atualiza um aluno | 200, 400, 404, 500 |
| DELETE | `/students/:id` | Remove um aluno | 200, 400, 404, 500 |

## Modelo de aluno

| Campo | Tipo | Observação |
|-------|------|------------|
| `name` | string | obrigatório |
| `cpf` | string | obrigatório e único (texto, para manter zeros à esquerda) |
| `email` | string | obrigatório e único |
| `age` | int | obrigatório, maior que zero |
| `active` | bool | obrigatório |

## Regras da API

- CPF ou e-mail já cadastrados retornam **409 Conflict**.
- Dados inválidos no cadastro retornam **400 Bad Request**, assim como um id que não é número.
- Na atualização (`PUT`), só os campos enviados são alterados; os demais continuam como estavam.
- A exclusão é lógica: o registro some das consultas, mas continua no banco com a data de exclusão.

## Exemplo

Cadastrar um aluno:

```bash
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{"name":"Ana","email":"ana@email.com","cpf":"01234567890","age":20,"active":true}'
```

Resposta (`201 Created`):

```json
{
  "id": 1,
  "createdAt": "2026-11-23T08:00:00-03:00",
  "updatedAt": "2026-11-23T08:00:00-03:00",
  "name": "Ana",
  "cpf": "01234567890",
  "email": "ana@email.com",
  "age": 20,
  "active": true
}
```

## Estrutura do projeto

```
api-students/
├── api/        # rotas, handlers e validação das requisições
├── db/         # acesso ao banco de dados (GORM)
├── docs/       # documentação Swagger gerada
├── schemas/    # modelos de dados e formato das respostas
└── main.go     # ponto de entrada
```

## Atualizar a documentação Swagger

Depois de alterar os comentários dos handlers:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g api/api.go
```

## Próximos passos

- Testes automáticos
- Dockerfile

## Autora

Feito por [Tatascruz](https://github.com/Tatascruz) durante meus estudos de backend com Go.