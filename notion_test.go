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
	}

	entries, err := db.QueryDatabaseByName("name", "uid")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", entries)

	entries, err = db.QueryDatabaseByEmail("email@email.com", "uid")
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
	}

	err := db.AddPageToDatabase("uid", "name", "title", "address", "emai@email.com", "phone")
	if err != nil {
		t.Fatal(err)
	}
}
