package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Items_20240210_190339 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Items_20240210_190339{}
	m.Created = "20240210_190339"

	migration.Register("Items_20240210_190339", m)
}

// Run the migrations
func (m *Items_20240210_190339) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE IF NOT EXISTS items (id INT AUTO_INCREMENT PRIMARY KEY, image VARCHAR(255), title VARCHAR(80), description TEXT, recipe TEXT, ingredients JSON, cook_time_in_minutes INT, author_id INT, FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *Items_20240210_190339) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS items")
}
