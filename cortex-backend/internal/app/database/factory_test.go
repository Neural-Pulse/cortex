package database

import (
	"testing"
)

func TestNewDatabaseFactory(t *testing.T) {
	tests := []struct {
		name   string
		dbType string
		want   string
	}{
		{"MySQL", "mysql", "MySQL"},
		{"PostgreSQL", "postgres", "PostgreSQL"},
		{"Unsupported", "unsupported", "Tipo de banco de dados não suportado"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if r != tt.want {
						t.Errorf("NewDatabaseFactory() = %v, want %v", r, tt.want)
					}
				}
			}()

			//db := NewDatabaseFactory(tt.dbType)
			if tt.dbType == "unsupported" {
				return
			}

			// Aqui você pode adicionar verificações adicionais para garantir que db é do tipo esperado.
			// Por exemplo, você pode verificar se db implementa a interface Database.
		})
	}
}
