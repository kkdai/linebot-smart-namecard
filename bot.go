package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// Const variables of Prompts.
const ImagePrompt = "你是一個美食烹飪專家，根據這張圖片給予相關的食物敘述，越詳細越好。"
const CalcPrompt = "根據這張圖片，幫我計算食物的卡路里。"
const CookPrompt = "根據這張圖片，幫我找到相關的食譜。"

// Image statics link.
const CalcImg = "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/calc.jpg"
const CookImg = "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/cooking.png"

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

func handleCameraQuickReply(replyToken string) error {
	msg := &messaging_api.TextMessage{
		Text: "請上傳一張美食照片，開始相關功能吧！",
		QuickReply: &messaging_api.QuickReply{
			Items: []messaging_api.QuickReplyItem{
				{
					ImageUrl: "",
					Action: &messaging_api.CameraAction{
						Label: "Camera",
					},
				},
			},
		},
	}
	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages:   []messaging_api.MessageInterface{msg},
		},
	); err != nil {
		return err
	}
	return nil
}

// callbackHandler: Handle callback from LINE server.
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	cb, err := webhook.ParseRequest(os.Getenv("ChannelSecret"), r)
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
				// Resceive text message, reply with QuickReply buttons.
				err := handleCameraQuickReply(e.ReplyToken)
				if err != nil {
					log.Print(err)
				}
			// Handle only on Sticker message
			case webhook.StickerMessageContent:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerId, message.PackageId, kw, message.Text)
				if err := replyText(e.ReplyToken, outStickerResult); err != nil {
					log.Print(err)
				}

			// Handle only image message
			case webhook.ImageMessageContent:
				log.Println("Got img msg ID:", message.Id)

				//Get image binary from LINE server based on message ID.
				data, err := GetImageBinary(blob, message.Id)
				if err != nil {
					log.Println("Got GetMessageContent err:", err)
					continue
				}

				ret, err := GeminiImage(data, ImagePrompt)
				if err != nil {
					ret = "無法辨識影片內容文字，請重新輸入:" + err.Error()
				}

				// Prepare QuickReply buttons.
				qReply := &messaging_api.QuickReply{
					Items: []messaging_api.QuickReplyItem{
						{
							ImageUrl: CalcImg,
							Action: &messaging_api.PostbackAction{
								Label:       "calc",
								Data:        "action=calc&m_id=" + message.Id,
								DisplayText: "",
								Text:        "計算卡路里",
							},
						}, {
							ImageUrl: CookImg,
							Action: &messaging_api.PostbackAction{
								Label:       "cook",
								Data:        "action=cook&m_id=" + message.Id,
								DisplayText: "",
								Text:        "建議食譜",
							},
						},
					},
				}

				// Determine the push msg target.
				if _, err := bot.ReplyMessage(
					&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							&messaging_api.TextMessage{
								Text:       ret,
								QuickReply: qReply,
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
			// Using urls value to parse event.Postback.Data strings.
			ret, err := url.ParseQuery(e.Postback.Data)
			if err != nil {
				log.Print("action parse err:", err, " dat=", e.Postback.Data)
				continue
			}

			log.Println("Action:", ret["action"])
			log.Println("Calc calories m_id:", ret["m_id"])

			// 取得用戶 ID
			var target string
			switch source := e.Source.(type) {
			case *webhook.UserSource:
				target = source.UserId
			case *webhook.GroupSource:
				target = source.UserId
			case *webhook.RoomSource:
				target = source.UserId
			}

			// Handle only on Postback message
			if ret["action"][0] == "calc" {
				// Determine the push msg target.
				go processImage(target, ret["m_id"][0], CalcPrompt, "GeminiImage", blob) // for calcCalories
			} else if ret["action"][0] == "cook" {
				// Determine the push msg target.
				go processImage(target, ret["m_id"][0], CookImg, "GeminiImage", blob) // for searchCooking
			}
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
