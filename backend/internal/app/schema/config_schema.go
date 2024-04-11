package schema

// SchemaSQL contém os comandos SQL para criar o banco de dados e as tabelas.
const SchemaSQL = `
CREATE DATABASE IF NOT EXISTS myDatabase;

USE myDatabase;

CREATE TABLE IF NOT EXISTS config_table (
    id INT AUTO_INCREMENT PRIMARY KEY,
    config_data TEXT NOT NULL
);
`
