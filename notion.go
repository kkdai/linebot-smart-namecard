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
}

type NotionDBEntry struct {
	Name        string
	Title       string
	Address     string
	Email       string
	PhoneNumber string
	Tags        []string
	ImgURL      string
}

// QueryDatabase 根據提供的屬性和值查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabase(UId, property, value string) ([]NotionDBEntry, error) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// Add UId to the filter conditions
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
					Equals: UId,
				},
			},
		},
	}

	// 調用 Notion API 來查詢資料庫
	result, err := client.Database.Query(context.Background(), notionapi.DatabaseID(n.DatabaseID), filter)
	if err != nil {
		return nil, err
	}

	var entries []NotionDBEntry

	for _, page := range result.Results {
		entry := n.createEntryFromPage(&page)
		entries = append(entries, entry)
	}
	return entries, nil
}

// createEntryFromPage creates a NotionDBEntry from a page.
func (n *NotionDB) createEntryFromPage(page *notionapi.Page) NotionDBEntry {
	entry := NotionDBEntry{}

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

// QueryDatabaseByName 根據提供的名稱和UId查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseByName(name, UId string) ([]NotionDBEntry, error) {
	return n.QueryDatabase(UId, "Name", name)
}

// QueryDatabaseByEmail 根據提供的電子郵件地址和UId查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseByEmail(email, UId string) ([]NotionDBEntry, error) {
	return n.QueryDatabase(UId, "Email", email)
}

// AddPageToDatabase adds a new page with the provided field values to the specified Notion database.
func (n *NotionDB) AddPageToDatabase(Uid string, name string, title string, address string, email string, phoneNumber string) error {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立 Properties 物件來設置頁面屬性
	properties := notionapi.Properties{
		"UID": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					PlainText: name,
					Text:      &notionapi.Text{Content: Uid},
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

	log.Println("Page added successfully.")
	return nil
}

// RetrieveDatabaseContents 從指定的 Notion 資料庫檢索內容並返回結構體切片。
func (n *NotionDB) RetrieveDatabaseContents() ([]NotionDBEntry, error) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 讀取 Notion 資料庫中的頁面
	query := &notionapi.DatabaseQueryRequest{}
	result, err := client.Database.Query(context.Background(), notionapi.DatabaseID(n.DatabaseID), query)
	if err != nil {
		return nil, err
	}

	var entries []NotionDBEntry

	for _, page := range result.Results {
		entry := NotionDBEntry{}

		if prop, ok := page.Properties["Name"].(*notionapi.TitleProperty); ok && len(prop.Title) > 0 {
			entry.Name = prop.Title[0].PlainText
		}

		if prop, ok := page.Properties["Title"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
			entry.Title = prop.RichText[0].PlainText
		}

		if prop, ok := page.Properties["Address"].(*notionapi.RichTextProperty); ok && len(prop.RichText) > 0 {
			entry.Address = prop.RichText[0].PlainText
		}

		if prop, ok := page.Properties["Email"].(*notionapi.EmailProperty); ok {
			entry.Email = prop.Email
		}

		if prop, ok := page.Properties["Phone Number"].(*notionapi.PhoneNumberProperty); ok {
			entry.PhoneNumber = prop.PhoneNumber
		}

		if tagsProp, ok := page.Properties["Tags"].(*notionapi.MultiSelectProperty); ok {
			for _, tag := range tagsProp.MultiSelect {
				entry.Tags = append(entry.Tags, tag.Name)
			}
		}

		if imgProp, ok := page.Properties["Img"].(*notionapi.FilesProperty); ok {
			for _, file := range imgProp.Files {
				if file.Type == "external" {
					entry.ImgURL = file.External.URL
					break
				}
			}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
