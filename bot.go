package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// Const variables of Prompts.
const ImagePrompt = "這是一張名片，你是一個名片秘書。請將以下資訊整理成 json 給我 Name, Title, Address, Email, Phone Number"

// replyText: Reply text message to LINE server.
func replyText(replyToken, text string) error {
	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				&messaging_api.TextMessage{
					Text: text,
				},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

// callbackHandler: Handle callback from LINE server.
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	card_prompt := os.Getenv("CARD_PROMPT")
	if card_prompt == "" {
		card_prompt = ImagePrompt
	}

	cb, err := webhook.ParseRequest(ChannelSecret, r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range cb.Events {
		log.Printf("Got event %v", event)
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			// Handle only on text message
			case webhook.TextMessageContent:
				log.Printf("Got text message, ID: %s, text: %s", message.Id, message.Text)

			// Handle only on Sticker message
			case webhook.StickerMessageContent:
				// log sticker id and package id.
				log.Printf("Got sticker message, packageID: %s, stickerID: %s", message.PackageId, message.StickerId)

			// Handle only image message
			case webhook.ImageMessageContent:
				log.Println("Got img msg ID:", message.Id)

				//Get image binary from LINE server based on message ID.
				data, err := GetImageBinary(blob, message.Id)
				if err != nil {
					log.Println("Got GetMessageContent err:", err)
					continue
				}

				ret, err := GeminiImage(data, card_prompt)
				if err != nil {
					ret = "無法辨識影片內容文字，請重新輸入:" + err.Error()
				}

				log.Println("Got GeminiImage ret:", ret)

				// Remove first and last line,	which are the backticks.
				jsonData := removeFirstAndLastLine(ret)
				log.Println("Got jsonData:", jsonData)

				// Parse json and insert NotionDB
				var person Person
				err = json.Unmarshal([]byte(jsonData), &person)
				if err != nil {
					log.Println("Error parsing JSON:", err)
				}

				nDB := &NotionDB{
					DatabaseID: os.Getenv("NOTION_DB_PAGEID"),
					Token:      os.Getenv("NOTION_INTEGRATION_TOKEN"),
				}

				err = nDB.AddPageToDatabase(person.Name, person.Title, person.Address, person.Email, person.PhoneNumber)
				if err != nil {
					log.Println("Error adding page to database:", err)
				}

				// Determine the push msg target.
				if _, err := bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.TextMessage{
								Text: ret,
							},
							&messaging_api.TextMessage{
								Text: jsonData,
							},
						},
					},
				); err != nil {
					log.Print(err)
				}

			// Handle only video message
			case webhook.VideoMessageContent:
				log.Println("Got video msg ID:", message.Id)

			default:
				log.Printf("Unknown message: %v", message)
			}
		case webhook.PostbackEvent:
			log.Printf("Got postback: %v", e.Postback.Data)
		case webhook.JoinEvent:
			log.Printf("Got join event")
		case webhook.FollowEvent:
			log.Printf("message: Got followed event")
		case webhook.BeaconEvent:
			log.Printf("Got beacon: " + e.Beacon.Hwid)
		}
	}
}

// ProcessImage: Process an image and reply with a text.
func processImage(target, m_id, prompt, errMsg string, blob *messaging_api.MessagingApiBlobAPI) {
	// Get image data
	data, err := GetImageBinary(blob, m_id)
	if err != nil {
		log.Printf("Got GetMessageContent err: %v", err)
		return
	}

	// Chat with Image
	ret, err := GeminiImage(data, prompt)
	if err != nil {
		log.Printf("Got %s err: %v", errMsg, err)
		return
	}

	// Determine the push msg target.
	if err := replyText(target, ret); err != nil {
		log.Print(err)
	}
}

// GetImageBinary: Get image binary from LINE server based on message ID.
func GetImageBinary(blob *messaging_api.MessagingApiBlobAPI, messageID string) ([]byte, error) {
	// Get image binary from LINE server based on message ID.
	content, err := blob.GetMessageContent(messageID)
	if err != nil {
		log.Println("Got GetMessageContent err:", err)
	}
	defer content.Body.Close()
	data, err := io.ReadAll(content.Body)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}

// extractJSONFromString takes a string that contains JSON with Markdown backticks and returns just the JSON string.
func extractJSONFromString(s string) string {
	// Trim the leading and trailing backticks.
	s = strings.Trim(s, "`")

	// Trim leading and trailing whitespace that may be present after removing backticks.
	s = strings.TrimSpace(s)

	return s
}

// removeFirstAndLastLine takes a string and removes the first and last lines.
func removeFirstAndLastLine(s string) string {
	// Split the string into lines.
	lines := strings.Split(s, "\n")

	// If there are less than 3 lines, return an empty string because removing the first and last would leave nothing.
	if len(lines) < 3 {
		return ""
	}

	// Join the lines back together, skipping the first and last lines.
	return strings.Join(lines[1:len(lines)-1], "\n")
}
