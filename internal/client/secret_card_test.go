package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func TestApp_UpdateSecretCard_CtrlC(t *testing.T) {
	m := app{
		state:          secretCardState,
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	newModel, cmd := m.updateSecretCard(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// Ctrl+C should not change state
	assert.Equal(t, m.state, newApp.state)
	assert.Equal(t, m.selectedSecret, newApp.selectedSecret)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateSecretCard_CtrlB(t *testing.T) {
	m := app{
		state:          99,
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	newModel, cmd := m.updateSecretCard(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// Ctrl+B:
	// 1. state should become secretsState
	// 2. selectedSecret is nil
	// 3. no commands
	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecretCard_OtherKey(t *testing.T) {
	m := app{
		state:          99,
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	newModel, cmd := m.updateSecretCard(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// Не Ctrl+C / Ctrl+B -> the state and selectedSecret do not change, there are no commands
	assert.Equal(t, m.state, newApp.state)
	assert.Equal(t, m.selectedSecret.Title, newApp.selectedSecret.Title)
	assert.Nil(t, cmd)
}

func TestApp_SecretCardView_Text(t *testing.T) {
	m := app{
		selectedSecret: &gophkeeperv1.Secret{
			Title:      "text-secret",
			SecretType: gophkeeperv1.SecretType_TEXT,
			TextData:   "some text data",
		},
	}

	output := m.secretCardView()

	assert.Contains(t, output, "Secret: text-secret")
	assert.Contains(t, output, "TextData: some text data")
	assert.Contains(t, output, "[Ctrl+B] – back in secrets, [Ctrl+C] – quit")
}

func TestApp_SecretCardView_Binary(t *testing.T) {
	m := app{
		selectedSecret: &gophkeeperv1.Secret{
			Title:      "binary-secret",
			SecretType: gophkeeperv1.SecretType_BINARY,
			BinaryData: []byte{1, 2, 3, 4},
		},
	}

	output := m.secretCardView()

	assert.Contains(t, output, "Secret: binary-secret")
	assert.Contains(t, output, "BinaryData: file with size 4 bytes")
	assert.Contains(t, output, "[Ctrl+B] – back in secrets, [Ctrl+C] – quit")
}

func TestApp_SecretCardView_Card(t *testing.T) {
	m := app{
		selectedSecret: &gophkeeperv1.Secret{
			Title:      "card-secret",
			SecretType: gophkeeperv1.SecretType_CARD,
			CardHolder: "John Doe",
			CardNumber: "1234567812345678",
			CardExp:    "12/25",
		},
	}

	output := m.secretCardView()

	assert.Contains(t, output, "Secret: card-secret")
	assert.Contains(t, output, "Card holder: John Doe")
	assert.Contains(t, output, "Card number: 1234567812345678")
	assert.Contains(t, output, "Card exp: 12/25")
	assert.Contains(t, output, "[Ctrl+B] – back in secrets, [Ctrl+C] – quit")
}
