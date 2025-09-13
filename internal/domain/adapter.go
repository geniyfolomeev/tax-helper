package domain

type Notifier interface {
	SendMessage(userID int64, message string) error
}
