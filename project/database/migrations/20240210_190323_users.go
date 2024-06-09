package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Users_20240210_190323 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Users_20240210_190323{}
	m.Created = "20240210_190323"

	migration.Register("Users_20240210_190323", m)
}

// Run the migrations
func (m *Users_20240210_190323) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE IF NOT EXISTS users (id INT AUTO_INCREMENT PRIMARY KEY, email VARCHAR(255) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL)")

}

// Reverse the migrations
func (m *Users_20240210_190323) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS users")
}
