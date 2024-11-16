package handler

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/oke11o/sb-habits-bot/internal/fsm/help"
)

func (s *Suite) TestHandler_Help() {
	s.T().Run("/start", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(help.InstructionText, msg.Text)
			},
		})
		err := h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/start"}, // Notice!!!
		})
		s.Require().NoError(err)
	})

	s.T().Run("/help", func(t *testing.T) {
		h := s.createHandler()
		h.SetSender(testSender{
			assert: func(c tgbotapi.Chattable) {
				msg, ok := c.(tgbotapi.MessageConfig)
				s.Require().True(ok)
				s.Require().Equal(help.InstructionText, msg.Text)
			},
		})
		err := h.HandleUpdate(context.Background(), tgbotapi.Update{
			UpdateID: 1,
			Message: &tgbotapi.Message{MessageID: 45,
				From: &tgbotapi.User{ID: 20, FirstName: "Tmp", LastName: "User", UserName: "tmp_user", LanguageCode: "en"},
				Date: 1712312739,
				Chat: &tgbotapi.Chat{ID: 20, Type: "private", UserName: "tmp_user", FirstName: "Tmp", LastName: "User"},
				Text: "/help"}, // Notice!!!
		})
		s.Require().NoError(err)
	})

	q := `select count(*) from user`
	var count int
	err := s.dbx.GetContext(context.Background(), &count, q)
	s.Require().NoError(err)
	s.Require().Equal(1, count)

	q = `select count(*) from income_request`
	err = s.dbx.GetContext(context.Background(), &count, q)
	s.Require().NoError(err)
	s.Require().Equal(2, count)
}
