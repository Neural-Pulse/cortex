package main

import (
	"log"

	"github.com/neural-pulse/cortex-backend/internal/app/database"
)

func main() {
	// Obtenha o tipo de banco de dados da solicitação da API
	dbType := "mysql" // Exemplo: poderia ser obtido de uma solicitação HTTP

	// Crie uma instância do banco de dados com base no tipo
	db := database.NewDatabaseFactory(dbType)

	// Configurar o banco de dados
	dsn := "user:password@tcp(localhost:3306)/database"
	dbConnection, err := db.SetupDatabase(dsn)
	if err != nil {
		log.Fatal("Erro ao configurar o banco de dados:", err)
	}

	// Verificar a conexão do banco de dados
	if err := db.PingDatabase(dbConnection); err != nil {
		log.Fatal("Erro ao pingar o banco de dados:", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
}
