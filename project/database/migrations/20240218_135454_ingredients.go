package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type Ingredients_20240218_135454 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Ingredients_20240218_135454{}
	m.Created = "20240218_135454"

	migration.Register("Ingredients_20240218_135454", m)
}

// Run the migrations
func (m *Ingredients_20240218_135454) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE IF NOT EXISTS ingredients (id INT AUTO_INCREMENT PRIMARY KEY,name VARCHAR(255) NOT NULL,quantity INT NOT NULL,units ENUM('шт.', 'г.', 'мл.') NOT NULL,recipe_id INT,FOREIGN KEY (recipe_id) REFERENCES items(id) ON DELETE CASCADE)")
}

// Reverse the migrations
func (m *Ingredients_20240218_135454) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE IF EXISTS ingredients")
}
