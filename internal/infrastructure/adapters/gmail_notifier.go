package adapters

type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (e *EmailNotifier) SendMessage(userID int64, message string) error {
	// отправка email
	return nil
}
