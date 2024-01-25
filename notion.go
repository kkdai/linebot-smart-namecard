package main

import (
	"context"
	"log"

	"github.com/jomei/notionapi"
)

// Person 定義了 JSON 資料的結構體
type Person struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

// DatabaseEntry 定義了 Notion 資料庫條目的結構體。
type NotionDB struct {
	DatabaseID string
	Token      string
	UID        string
}

// QueryDatabase 根據提供的屬性和值查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabase(property, value string) ([]Person, error) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立查詢過濾條件
	filter := &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: property,
				RichText: &notionapi.TextFilterCondition{
					Equals: value,
				},
			},
			notionapi.PropertyFilter{
				Property: "UID",
				RichText: &notionapi.TextFilterCondition{
					Equals: n.UID,
				},
			},
		},
	}

	// 調用 Notion API 來查詢資料庫
	result, err := client.Database.Query(context.Background(), notionapi.DatabaseID(n.DatabaseID), filter)
	if err != nil {
		return nil, err
	}

	var entries []Person

	for _, page := range result.Results {
		entry := n.createEntryFromPage(&page)
		entries = append(entries, entry)
	}
	return entries, nil
}

// QueryContainsDatabase 根據提供的屬性和值查詢 Notion 資料庫。
func (n *NotionDB) QueryContainsDatabase(property, value string) ([]Person, error) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立查詢過濾條件
	filter := &notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{
			notionapi.PropertyFilter{
				Property: property,
				RichText: &notionapi.TextFilterCondition{
					Contains: value,
				},
			},
			notionapi.PropertyFilter{
				Property: "UID",
				RichText: &notionapi.TextFilterCondition{
					Equals: n.UID,
				},
			},
		},
	}

	// 調用 Notion API 來查詢資料庫
	result, err := client.Database.Query(context.Background(), notionapi.DatabaseID(n.DatabaseID), filter)
	if err != nil {
		return nil, err
	}

	var entries []Person

	for _, page := range result.Results {
		entry := n.createEntryFromPage(&page)
		entries = append(entries, entry)
	}
	return entries, nil
}

// QueryDatabaseByEmail 根據提供的電子郵件地址查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseByEmail(email string) ([]Person, error) {
	return n.QueryDatabase("Email", email)
}

// AddPageToDatabase adds a new page with the provided field values to the specified Notion database.
func (n *NotionDB) AddPageToDatabase(name string, title string, address string, email string, phoneNumber string) error {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立 Properties 物件來設置頁面屬性
	properties := notionapi.Properties{
		"UID": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: n.UID},
				},
			},
		},
		"Name": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: name},
				},
			},
		},
		"Title": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: title},
				},
			},
		},
		"Address": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: address},
				},
			},
		},
		"Email": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: email},
				},
			},
		},
		"Phone Number": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: phoneNumber},
				},
			},
		},
	}

	// 創建一個新頁面的請求
	pageRequest := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(n.DatabaseID),
		},
		Properties: properties,
	}

	// 調用 Notion API 來創建新頁面
	_, err := client.Page.Create(context.Background(), pageRequest)
	if err != nil {
		log.Println("Error creating page:", err)
		return err
	}

	log.Println("Page added successfully:", n.UID, name, title, address, email, phoneNumber)
	return nil
}

// createEntryFromPage creates a Person from a page.
func (n *NotionDB) createEntryFromPage(page *notionapi.Page) Person {
	entry := Person{}

	entry.Name = n.getPropertyValue(page, "Name")
	entry.Title = n.getPropertyValue(page, "Title")
	entry.Address = n.getPropertyValue(page, "Address")
	entry.Email = n.getPropertyValue(page, "Email")
	entry.PhoneNumber = n.getPropertyValue(page, "Phone Number")

	return entry
}

// getPropertyValue gets the plain text value of a property from a page.
func (n *NotionDB) getPropertyValue(page *notionapi.Page, property string) string {
	if prop, ok := page.Properties[property].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
		return prop.RichText[0].PlainText
	}

	return ""
}

// QueryDatabaseByName 根據提供的名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseByName(name string) ([]Person, error) {
	return n.QueryDatabase("Name", name)
}

// QueryDatabaseContainsByName 根據提供的名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseContainsByName(name string) ([]Person, error) {
	return n.QueryContainsDatabase("Name", name)
}

// QueryDatabaseContainsByEmail 根據提供的名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseContainsByEmail(name string) ([]Person, error) {
	return n.QueryContainsDatabase("Email", name)
}

// QueryDatabaseContainsByName 根據提供的名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseContainsByTitle(name string) ([]Person, error) {
	return n.QueryContainsDatabase("Title", name)
}

// QueryDatabaseContains 根據提供的名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseContains(query string) ([]Person, error) {
	// 初始化一個空的結果集
	var combinedResult []Person

	// 進行名稱查詢
	log.Println("QueryDatabaseContainsByName", query, n.UID)
	nameResult, err := n.QueryDatabaseContainsByName(query)
	log.Println("QueryDatabaseContainsByName", nameResult, err)
	if err != nil {
		return nil, err
	}
	combinedResult = append(combinedResult, nameResult...)

	// 進行電子郵件查詢
	log.Println("QueryDatabaseContainsByEmail", query, n.UID)
	emailResult, err := n.QueryDatabaseContainsByEmail(query)
	log.Println("QueryDatabaseContainsByEmail", emailResult, err)
	if err != nil {
		return nil, err
	}
	combinedResult = append(combinedResult, emailResult...)

	// 進行標題查詢
	log.Println("QueryDatabaseContainsByTitle", query, n.UID)
	titleResult, err := n.QueryDatabaseContainsByTitle(query)
	log.Println("QueryDatabaseContainsByTitle", titleResult, err)
	if err != nil {
		return nil, err
	}
	combinedResult = append(combinedResult, titleResult...)

	// 返回結合的結果
	return combinedResult, nil
}
