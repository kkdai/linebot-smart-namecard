package main

import "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"

// SendFlexMsg: Send flex message to LINE server.
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
						&messaging_api.FlexBox{
							BorderColor: "#6EC4C4",
							BorderWidth: "1px",
							Flex:        0,
							Height:      "120px",
							Layout:      messaging_api.FlexBoxLAYOUT_VERTICAL,
							Contents: []messaging_api.FlexComponentInterface{
								&messaging_api.FlexFiller{},
							},
						},
						&messaging_api.FlexBox{
							Flex:   3,
							Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
							Contents: []messaging_api.FlexComponentInterface{
								&messaging_api.FlexBox{
									Flex:   1,
									Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
									Contents: []messaging_api.FlexComponentInterface{
										&messaging_api.FlexFiller{},
									},
								},
								&messaging_api.FlexText{
									Color:  "#6EC4C4",
									Size:   "sm",
									Text:   "Company",
									Weight: "bold",
								},
								&messaging_api.FlexText{
									Color:  "#81C997",
									Margin: "xxl",
									Size:   "xxs",
									Text:   "Title",
								},
								&messaging_api.FlexText{
									Color:  "#81C997",
									Size:   "xl",
									Text:   "Name",
									Weight: "bold",
								},
								&messaging_api.FlexText{
									Text: "address",
								},
								&messaging_api.FlexText{
									Text: "email",
								},
								&messaging_api.FlexText{
									Text: "phone",
								},
							},
						},
					},
				},
			},
		},
	}

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
