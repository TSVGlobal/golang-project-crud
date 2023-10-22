package main

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT * FROM products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query DB: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func getProduct(db *sql.DB, id int) (Product, error) {
	query := "SELECT * FROM products WHERE id = ?"
	row := db.QueryRow(query, id)

	var p Product
	if err := row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Price); err != nil {
		return p, fmt.Errorf("could not scan row: %w", err)
	}
	return p, nil
}

func (p *Product) createProduct(db *sql.DB) (Product, error) {
	query := "INSERT INTO products(name, quantiry, price) VALUES(?, ?, ?)"
	result, err := db.Exec(query, p.Name, p.Quantity, p.Price)
	if err != nil {
		return *p, fmt.Errorf("could not exec query: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return *p, fmt.Errorf("could not get last insert id: %w", err)
	}
	p.ID = int(id)
	return *p, nil
}

func (p *Product) updateProduct(db *sql.DB) error {
	query := "UPDATE products SET name = ?, quantiry = ?, price = ? WHERE id = ?"

	result, err := db.Exec(query, p.Name, p.Quantity, p.Price, p.ID)
	if err != nil {
		return fmt.Errorf("could not exec query: %w", err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)

	}
	return nil

}

func (p *Product) deleteProduct(db *sql.DB, id int) error {
	query := "DELETE FROM products WHERE id = ?"
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not exec query: %w", err)
	}
	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %w", err)
	}
	return nil
}
