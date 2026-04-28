package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"
)

func TestApp_UpdateCreateCardSecret_CtrlB(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.Focus()

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateCardSecret_CtrlC(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		token:  "Bearer test",
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.Focus()

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, a.state, newApp.state)
	assert.NotNil(t, cmd) // tea.Quit
}

func TestApp_UpdateCreateCardSecret_TitleRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.cardNumber.SetValue("1234")
	a.createCardSecretModel.cardHolder.SetValue("John Doe")
	a.createCardSecretModel.cardExp.SetValue("12/25")

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createCardSecretModel.err.Error(), "title is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateCardSecret_CardExpRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.SetValue("test-title")
	a.createCardSecretModel.cardNumber.SetValue("1234")
	a.createCardSecretModel.cardHolder.SetValue("John Doe")

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createCardSecretModel.err.Error(), "card exp is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateCardSecret_CardNumberRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.SetValue("test-title")
	a.createCardSecretModel.cardHolder.SetValue("John Doe")
	a.createCardSecretModel.cardExp.SetValue("12/25")

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createCardSecretModel.err.Error(), "card number is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateCardSecret_CardHolderRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretsState,
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.SetValue("test-title")
	a.createCardSecretModel.cardNumber.SetValue("1234")
	a.createCardSecretModel.cardExp.SetValue("12/25")

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createCardSecretModel.err.Error(), "card holder is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateCardSecret_ValidInput(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretCardCreateState,
		token:  "Bearer test",
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
		secretsModel: secretsPage{},
	}

	a.createCardSecretModel.title.SetValue("test-title")
	a.createCardSecretModel.metadata.SetValue("test-meta")
	a.createCardSecretModel.cardNumber.SetValue("1234")
	a.createCardSecretModel.cardHolder.SetValue("John Doe")
	a.createCardSecretModel.cardExp.SetValue("12/25")

	newModel, cmd := a.updateCreateCardSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Empty(t, newApp.createCardSecretModel.title.Value())
	assert.Empty(t, newApp.createCardSecretModel.metadata.Value())
	assert.Empty(t, newApp.createCardSecretModel.cardNumber.Value())
	assert.Empty(t, newApp.createCardSecretModel.cardHolder.Value())
	assert.Empty(t, newApp.createCardSecretModel.cardExp.Value())
	assert.Nil(t, newApp.createCardSecretModel.err)
	assert.NotNil(t, cmd)
}

func TestApp_CreateCardSecretView(t *testing.T) {
	a := app{
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			cardNumber: textinput.New(),
			cardHolder: textinput.New(),
			cardExp:    textinput.New(),
		},
	}

	a.createCardSecretModel.title.SetValue("test-card")
	a.createCardSecretModel.metadata.SetValue("test-meta")
	a.createCardSecretModel.cardNumber.SetValue("1234")
	a.createCardSecretModel.cardHolder.SetValue("John Doe")
	a.createCardSecretModel.cardExp.SetValue("12/25")

	s := a.createCardSecretView()

	assert.Contains(t, s, "New secret (type card)")
	assert.Contains(t, s, "Title*> test-card")
	assert.Contains(t, s, "Metadata> test-meta")
	assert.Contains(t, s, "Card number*> 1234")
	assert.Contains(t, s, "Card holder*> John Doe")
	assert.Contains(t, s, "Card exp*> 12/25")
	assert.Contains(t, s, "[Ctrl+B] - back in secrets")
}
