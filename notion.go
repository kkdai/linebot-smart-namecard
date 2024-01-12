package main

import (
	"context"
	"log"

	"github.com/jomei/notionapi"
)

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

// QueryDatabaseByTitleAndName 根據提供的標題和名稱查詢 Notion 資料庫。
func (n *NotionDB) QueryDatabaseByTitleAndName(name string) ([]NotionDBEntry, error) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立查詢過濾條件
	filter := &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Name",
			// Replace Text with the correct field based on the notionapi package's documentation or source code
			RichText: &notionapi.TextFilterCondition{
				Equals: name,
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

// AddPageToDatabase adds a new page with the provided field values to the specified Notion database.
func (n *NotionDB) AddPageToDatabase(name string, title string, address string, email string, phoneNumber string, tags []string, imgURL string) {
	client := notionapi.NewClient(notionapi.Token(n.Token))

	// 建立 Properties 物件來設置頁面屬性
	properties := notionapi.Properties{
		"Name": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: name}},
			},
		},
		"Title": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: title}},
			},
		},
		"Address": notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: address}},
			},
		},
		// ... 其他屬性依此類推
	}

	// 如果 Email 不為空，添加 Email 屬性
	if email != "" {
		properties["Email"] = notionapi.EmailProperty{Email: email}
	}

	// 如果 Phone Number 不為空，添加 Phone Number 屬性
	if phoneNumber != "" {
		properties["Phone Number"] = notionapi.PhoneNumberProperty{PhoneNumber: phoneNumber}
	}

	// 如果 Tags 不為空，添加 Tags 屬性
	if len(tags) > 0 {
		multiSelect := make([]notionapi.Option, 0, len(tags))
		for _, tag := range tags {
			multiSelect = append(multiSelect, notionapi.Option{Name: tag})
		}
		properties["Tags"] = notionapi.MultiSelectProperty{MultiSelect: multiSelect}
	}

	// 如果 ImgURL 不為空，添加 Img 屬性
	if imgURL != "" {
		properties["Img"] = notionapi.FilesProperty{
			Files: []notionapi.File{
				{
					Name: "image.jpg",
					Type: "notionapi.FileTypeImage",
					File: &notionapi.FileObject{
						URL:        imgURL,
						ExpiryTime: nil, // or some time.Time value
					},
				},
			},
		}
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
		log.Fatalf("Error creating page: %v", err)
	}

	log.Println("Page added successfully.")
}

// func main() {
// 	// 替換成你的 Notion Integration Token
// 	token := "your_notion_integration_token"

// 	// 替換成你的 Notion 資料庫 ID
// 	databaseID := "your_notion_database_id"

// 	// 調用函數來添加一個新頁面到資料庫
// 	AddPageToDatabase(token, databaseID, "Example Name", "Example Title", "Example Address", "example@email.com", "1234567890", []string{"Tag1", "Tag2"}, "https://example.com/image.jpg")
// }

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

// func main() {
// 	// 替換成你的 Notion Integration Token
// 	token := "your_notion_integration_token"

// 	// 替換成你的 Notion 資料庫 ID
// 	databaseID := "your_notion_database_id"

// 	// 調用函數並接收返回的結構體切片
// 	entries, err := RetrieveDatabaseContents(token, databaseID)
// 	if err != nil {
// 		log.Fatalf("Error retrieving database contents: %v", err)
// 	}

// 	// 這裡你可以處理 entries 切片，例如列印或其他操作
// 	for _, entry := range entries {
// 		log.Printf("%+v\n", entry)
// 	}
// }
