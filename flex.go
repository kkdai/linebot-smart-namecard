package main

import "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

//  {
//  "type": "carousel",
// 	"contents": [
// 		{
// 			"type": "bubble",
// 			"size": "giga",
// 			"body": {
// 				"layout": "horizontal",
// 				"spacing": "lg",
// 				"type": "box",
// 				"contents": [
// 					{
// 						"layout": "vertical",
// 						"type": "box",
// 						"width": "100px",
// 						"contents": [
// 							{
// 								"flex": 1,
// 								"layout": "vertical",
// 								"type": "box",
// 								"contents": [
// 									{
// 										"type": "filler"
// 									}
// 								]
// 							},
// 							{
// 								"height": "100px",
// 								"layout": "vertical",
// 								"type": "box",
// 								"width": "100px",
// 								"contents": [
// 									{
// 										"align": "center",
// 										"aspectMode": "cover",
// 										"aspectRatio": "1:1",
// 										"gravity": "center",
// 										"type": "image",
// 										"url": "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/logo.jpeg"
// 									}
// 								]
// 							},
// 							{
// 								"flex": 1,
// 								"layout": "vertical",
// 								"type": "box",
// 								"contents": [
// 									{
// 										"type": "filler"
// 									}
// 								]
// 							}
// 						]
// 					},
// 					{
// 						"borderColor": "#6EC4C4",
// 						"borderWidth": "1px",
// 						"flex": 0,
// 						"height": "120px",
// 						"layout": "vertical",
// 						"type": "box",
// 						"contents": [
// 							{
// 								"type": "filler"
// 							}
// 						]
// 					},
// 					{
// 						"flex": 3,
// 						"layout": "vertical",
// 						"type": "box",
// 						"contents": [
// 							{
// 								"flex": 1,
// 								"layout": "vertical",
// 								"type": "box",
// 								"contents": [
// 									{
// 										"type": "filler"
// 									}
// 								]
// 							},
// 							{
// 								"color": "#6EC4C4",
// 								"size": "sm",
// 								"text": "Company",
// 								"type": "text",
// 								"weight": "bold"
// 							},
// 							{
// 								"color": "#81C997",
// 								"margin": "xxl",
// 								"size": "xxs",
// 								"type": "text",
// 								"text": "Title"
// 							},
// 							{
// 								"color": "#81C997",
// 								"size": "xl",
// 								"text": "Name",
// 								"type": "text",
// 								"weight": "bold"
// 							},
// 							{
// 								"type": "text",
// 								"text": "address"
// 							},
// 							{
// 								"type": "text",
// 								"text": "email"
// 							},
// 							{
// 								"type": "text",
// 								"text": "phone"
// 							}
// 						]
// 					}
// 				]
// 			}
// 		},

func SendFlexMsg(replyToken string) error {
	contents := &messaging_api.FlexCarousel{
		Contents: []messaging_api.FlexBubble{
			{
				Size: messaging_api.FlexBubbleSIZE_GIGA,
				Body: &messaging_api.FlexBox{
					Layout:  messaging_api.FlexBoxLAYOUT_HORIZONTAL,
					Spacing: "lg",
					Contents: []messaging_api.FlexComponentInterface{
						&messaging_api.FlexBox{
							Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
							Width:  "100px",
							Contents: []messaging_api.FlexComponentInterface{
								&messaging_api.FlexBox{
									Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
									Flex:   1,
									Contents: []messaging_api.FlexComponentInterface{
										&messaging_api.FlexFiller{},
									},
								},
								&messaging_api.FlexBox{
									Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
									Width:  "100px",
									Height: "100px",
									Contents: []messaging_api.FlexComponentInterface{
										&messaging_api.FlexImage{
											Align:       "center",
											AspectMode:  "cover",
											AspectRatio: "1:1",
											Gravity:     "center",
											Url:         "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/logo.jpeg",
										},
									},
								},
								&messaging_api.FlexBox{
									Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
									Flex:   1,
									Contents: []messaging_api.FlexComponentInterface{
										&messaging_api.FlexFiller{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// 									{
	// 										"align": "center",
	// 										"aspectMode": "cover",
	// 										"aspectRatio": "1:1",
	// 										"gravity": "center",
	// 										"type": "image",
	// 										"url": "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/logo.jpeg"
	// 									}
	// 								]
	// 							},
	// 							{
	// 								"flex": 1,
	// 								"layout": "vertical",
	// 								"type": "box",
	// 								"contents": [
	// 									{
	// 										"type": "filler"
	// 									}
	// 								]
	// 							}
	// 						]
	// 					},
	// 					{
	// 						"borderColor": "#6EC4C4",
	// 						"borderWidth": "1px",
	// 						"flex": 0,
	// 						"height": "120px",
	// 						"layout": "vertical",
	// 						"type": "box",
	// 						"contents": [
	// 							{
	// 								"type": "filler"
	// 							}
	// 						]
	// 					},
	// 					{
	// 						"flex": 3,
	// 						"layout": "vertical",
	// 						"type": "box",
	// 						"contents": [
	// 							{
	// 								"flex": 1,
	// 								"layout": "vertical",
	// 								"type": "box",
	// 								"contents": [
	// 									{
	// 										"type": "filler"
	// 									}
	// 								]
	// 							},
	// 							{
	// 								"color": "#6EC4C4",
	// 								"size": "sm",
	// 								"text": "Company",
	// 								"type": "text",
	// 								"weight": "bold"
	// 							},
	// 							{
	// 								"color": "#81C997",
	// 								"margin": "xxl",
	// 								"size": "xxs",
	// 								"type": "text",
	// 								"text": "Title"
	// 							},
	// 							{
	// 								"color": "#81C997",
	// 								"size": "xl",
	// 								"text": "Name",
	// 								"type": "text",
	// 								"weight": "bold"
	// 							},
	// 							{
	// 								"type": "text",
	// 								"text": "address"
	// 							},
	// 							{
	// 								"type": "text",
	// 								"text": "email"
	// 							},
	// 							{
	// 								"type": "text",
	// 								"text": "phone"
	// 							}
	// 						]
	// 					}
	// 				]
	// 			}
	// 		},

	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{&messaging_api.FlexMessage{
				Contents: contents,
				AltText:  "Flex message alt text",
			}},
		},
	); err != nil {
		return err
	}
	return nil
}
