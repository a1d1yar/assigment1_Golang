package mailer

import (
	"testing"
)

func TestMailer_Send(t *testing.T) {
	m := New("smtp.mailtrap.io", 2525, "6a5a341ed3b09b", "a600a72a5771db", "sender@example.com")

	err := m.Send("recipient@example.com", "user_welcome.tmpl", map[string]interface{}{"key": "value"})

	if err != nil {
		t.Errorf("Error sending email: %v", err)
	}
}

func TestMailer_Send_InvalidRecipient(t *testing.T) {
	m := New("smtp.mailtrap.io", 2525, "6a5a341ed3b09b", "a600a72a5771db", "sender@example.com")

	err := m.Send("invalid-email", "user_welcome.tmpl", map[string]interface{}{"key": "value"})

	if err == nil {
		t.Error("Expected error for invalid recipient but got none")
	}
}

func TestMailer_Send_TemplateNotFound(t *testing.T) {
	m := New("smtp.mailtrap.io", 2525, "6a5a341ed3b09b", "a600a72a5771db", "sender@example.com")

	err := m.Send("recipient@example.com", "nonexistent_template.txt", map[string]interface{}{"key": "value"})

	if err == nil {
		t.Error("Expected error for non-existent template but got none")
	}
}

func TestMailer_Send_ValidEmail(t *testing.T) {
	m := New("smtp.mailtrap.io", 2525, "6a5a341ed3b09b", "a600a72a5771db", "sender@example.com")

	recipient := "recipient@example.com"
	templateFile := "user_welcome.tmpl"
	data := map[string]interface{}{
		"key": "value",
	}
	err := m.Send(recipient, templateFile, data)

	if err != nil {
		t.Errorf("Error sending email: %v", err)
	}
}
