package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func getDB(c string) *sql.DB {
	if c == "in" {
		db, err := sql.Open("mysql",
			"root:root@tcp(127.0.0.1:3306)/demo_magento2_20171222")
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatal("DB connection error!")
		}
		return db
	}
	return nil
}

func showBundleSkus(country string, skus []string) {
	db := getDB(country)

	var sku string
	s := strings.Join(skus, ",")
	q := `SELECT DISTINCT sku FROM
		(SELECT value, row_id FROM catalog_product_entity_varchar
		   WHERE attribute_id = 170 AND FIND_IN_SET(value, ?)) s
		LEFT JOIN catalog_product_entity c
		ON s.row_id = c.row_id
	      ORDER BY sku`
	stmt, _ := db.Prepare(q)
	rows, err := stmt.Query(s)
	defer stmt.Close()

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println()
	m := "Related bundle products:"
	fmt.Println(m)
	fmt.Println(strings.Repeat("-", utf8.RuneCountInString(m)))
	fmt.Println()

	c := 0
	for rows.Next() {
		err = rows.Scan(&sku)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(sku)
		c++
	}
	fmt.Println()
	fmt.Println("Total records:", c)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify sku list, separated by comma.\n")
		fmt.Println("usage: hpquery sku1[,sku2,...]")
		os.Exit(0)
	}

	showBundleSkus("in", os.Args[1:])
}
