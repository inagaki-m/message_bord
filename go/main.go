package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type user struct {
	Name string `db:"name"`
}

func mainTest() {
	// MySQLデータベースに接続する
	db, err := sqlx.Open("mysql", "root:myrootpassword@tcp(127.0.0.1:3306)/test_1")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatalln("エラー1: ", err)
	}

	// データベースにクエリを送信する
	query := `INSERT INTO users(
		name
	)
	VALUES(
		:name
	)`

	record := &user{
		Name: "ああああああ",
	}
	if _, err = tx.NamedExecContext(ctx, query, record); err != nil {
		log.Fatalln("エラー2: ", err)
	}

	tx.Commit()

	rows, err := db.Queryx("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	// 結果を処理する
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			panic(err.Error())
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
}
