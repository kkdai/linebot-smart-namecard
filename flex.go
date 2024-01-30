package main

import (
	"net/url"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

const LogoImageUrl = "https://raw.githubusercontent.com/kkdai/linebot-smart-namecard/main/img/logo.jpeg"

// getSingleCard: Send flex message to LINE server.
func getSingleCard(card Person) messaging_api.FlexBubble {

	return messaging_api.FlexBubble{
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
									Url:         LogoImageUrl,
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
							Text:   card.Company,
							Weight: "bold",
						},
						&messaging_api.FlexText{
							Color:  "#81C997",
							Margin: "xxl",
							Size:   "xxs",
							Text:   card.Title,
						},
						&messaging_api.FlexText{
							Color:  "#81C997",
							Size:   "xl",
							Text:   card.Name,
							Weight: "bold",
						},
						&messaging_api.FlexText{
							Text: card.Address,
						},
						&messaging_api.FlexText{
							Text: card.Email,
						},
						&messaging_api.FlexText{
							Text: card.Phone,
						},
					},
				},
			},
		},
	}
}

// SendFlexMsg: Send flex message to LINE server.
func SendFlexMsg(replyToken string, people []Person, msg string) error {
	var cards []messaging_api.FlexBubble
	for _, card := range people {
		cards = append(cards, getNewSingleCard(card))
	}

	contents := &messaging_api.FlexCarousel{
		Contents: cards,
	}

	if _, err := bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				&messaging_api.TextMessage{
					Text: msg,
				},
				&messaging_api.FlexMessage{
					Contents: contents,
					AltText:  "請到手機上查看名片資訊",
				},
			},
		},
	); err != nil {
		return err
	}
	return nil
}

// getNewSingleCard: Send flex message to LINE server.
func getNewSingleCard(card Person) messaging_api.FlexBubble {
	// Get URL encode for company name and address
	companyEncode := url.QueryEscape(card.Company)
	addressEncode := url.QueryEscape(card.Address)

	return messaging_api.FlexBubble{
		Size: messaging_api.FlexBubbleSIZE_GIGA,
		Body: &messaging_api.FlexBox{
			Layout:  messaging_api.FlexBoxLAYOUT_HORIZONTAL,
			Spacing: "md",
			Contents: []messaging_api.FlexComponentInterface{
				&messaging_api.FlexImage{
					AspectMode:  "cover",
					AspectRatio: "1:1",
					Flex:        1,
					Size:        "full",
					Url:         LogoImageUrl,
				},
				&messaging_api.FlexBox{
					Flex:   4,
					Layout: messaging_api.FlexBoxLAYOUT_VERTICAL,
					Contents: []messaging_api.FlexComponentInterface{
						&messaging_api.FlexText{
							Align:  "end",
							Size:   "xxl",
							Text:   card.Name,
							Weight: "bold",
						},
						&messaging_api.FlexText{
							Align: "end",
							Size:  "sm",
							Text:  card.Title,
						},
						&messaging_api.FlexText{
							Align:  "end",
							Margin: "xxl",
							Size:   "lg",
							Text:   card.Company,
							Weight: "bold",
							Action: &messaging_api.UriAction{
								Uri: "https://www.google.com/maps/search/?api=1&query=" + companyEncode + "&openExternalBrowser=1",
							},
						},
						&messaging_api.FlexText{
							Align: "end",
							Size:  "sm",
							Text:  card.Address,
							Action: &messaging_api.UriAction{
								Uri: "https://www.google.com/maps/search/?api=1&query=" + addressEncode + "&openExternalBrowser=1",
							},
						},
						&messaging_api.FlexText{
							Align:  "end",
							Margin: "xxl",
							Text:   card.Phone,
							Action: &messaging_api.UriAction{
								Uri: "tel:" + card.Phone,
							},
						},
						&messaging_api.FlexText{
							Align: "end",
							Text:  card.Email,
							Action: &messaging_api.UriAction{
								Uri: "mailto:" + card.Email,
							},
						},
						&messaging_api.FlexText{
							Align: "end",
							Text:  "更多資訊",
							Action: &messaging_api.UriAction{
								Uri: "https://github.com/kkdai/linebot-smart-namecard",
							},
						},
					},
				},
			},
		},
	}
}
