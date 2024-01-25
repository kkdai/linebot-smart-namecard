package main

import (
	"fmt"
	"os"
	"testing"
)

func TestQueryNotionDB(t *testing.T) {
	token := os.Getenv("NOTION_INTEGRATION_TOKEN")
	pageid := os.Getenv("NOTION_DB_PAGEID")

	// If not set token and pageid , skip this test
	if token == "" || pageid == "" {
		t.Skip("NOTION_INTEGRATION_TOKEN or NOTION_DB_PAGEID not set")
	}

	db := &NotionDB{
		DatabaseID: pageid,
		Token:      token,
		UID:        "uid",
	}

	entries, err := db.QueryDatabaseByName("name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", entries)

	entries, err = db.QueryDatabaseByEmail("email@email.com")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", entries)
}

func TestAddNotionDB(t *testing.T) {
	token := os.Getenv("NOTION_INTEGRATION_TOKEN")
	pageid := os.Getenv("NOTION_DB_PAGEID")

	// If not set token and pageid , skip this test
	if token == "" || pageid == "" {
		t.Skip("NOTION_INTEGRATION_TOKEN or NOTION_DB_PAGEID not set")
	}

	db := &NotionDB{
		DatabaseID: pageid,
		Token:      token,
		UID:        "uid",
	}

	err := db.AddPageToDatabase("name", "title", "address", "emai@email.com", "phone")
	if err != nil {
		t.Fatal(err)
	}
}

func TestQueryContainNotionDB(t *testing.T) {
	token := os.Getenv("NOTION_INTEGRATION_TOKEN")
	pageid := os.Getenv("NOTION_DB_PAGEID")

	// If not set token and pageid , skip this test
	if token == "" || pageid == "" {
		t.Skip("NOTION_INTEGRATION_TOKEN or NOTION_DB_PAGEID not set")
	}

	db := &NotionDB{
		DatabaseID: pageid,
		Token:      token,
		UID:        "uid",
	}

	entries, err := db.QueryDatabaseContainsByName("name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", entries)

	entries, err = db.QueryDatabaseContainsByEmail("email")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", entries)

	//test contains all columns (name, title, email)
	entries, err = db.QueryDatabaseContains("keyword")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", entries)
}
