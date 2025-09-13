package commands

import (
	"context"
	"tax-helper/internal/infrastructure/bot/commands/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		text    string
		err     error
		expText string
	}{
		{
			name:    "valid 2 arguments",
			text:    " 01.01.2025 1000",
			err:     nil,
			expText: "Successfully registered",
		},
		{
			name:    "valid 3 arguments",
			text:    " 01.01.2025 1000 01.01.2025",
			err:     nil,
			expText: "Successfully registered",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewMockentrepreneurCreator(ctrl)
			service.EXPECT().
				CreateEntrepreneur(
					gomock.Any(),
					int64(12345),
					gomock.Any(),
					gomock.Any(),
					1000.0,
				).
				Return(tt.err)

			handler := NewRegisterHandler(service, nil)
			msg := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 12345},
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
					},
				},
				Text: tt.text,
			}

			result := handler.Handle(context.Background(), msg)
			if result.ChatID != msg.Chat.ID || result.Text != tt.expText {
				t.Fatalf("unexpected result: %+v", result)
			}
		})
	}
}
