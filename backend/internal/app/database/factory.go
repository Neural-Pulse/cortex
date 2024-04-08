package database

// NewDatabaseFactory retorna uma instância de Database com base no tipo fornecido.
func NewDatabaseFactory(dbType string) Database {
	switch dbType {
	case "mysql":
		return MySQL{}
	case "postgres":
		return PostgreSQL{}
	// Adicione outros tipos de banco de dados aqui conforme necessário
	default:
		panic("Tipo de banco de dados não suportado")
	}
}
