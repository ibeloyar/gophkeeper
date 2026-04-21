package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"

	tea "github.com/charmbracelet/bubbletea"
)

func TestApp_UpdateCreateLogPassSecret_CtrlC(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		token:  "Bearer test",
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.title.Focus()

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, a.state, newApp.state)
	assert.NotNil(t, cmd) // tea.Quit
}

func TestApp_UpdateCreateLogPassSecret_CtrlB(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.title.Focus()

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateLogPassSecret_TitleRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.login.SetValue("testuser")
	a.createLogPassSecretModel.password.SetValue("testpass")

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createLogPassSecretModel.err.Error(), "title is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateLogPassSecret_LoginRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.title.SetValue("test-title")
	a.createLogPassSecretModel.password.SetValue("testpass")

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createLogPassSecretModel.err.Error(), "login is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateLogPassSecret_PasswordRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.title.SetValue("test-title")
	a.createLogPassSecretModel.login.SetValue("testuser")

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createLogPassSecretModel.err.Error(), "password is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateLogPassSecret_ValidInput(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretLogPassCreateState,
		token:  "Bearer test",
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
		secretsModel: secretsPage{},
	}

	a.createLogPassSecretModel.title.SetValue("test-login-pass")
	a.createLogPassSecretModel.metadata.SetValue("test-meta")
	a.createLogPassSecretModel.login.SetValue("testuser")
	a.createLogPassSecretModel.password.SetValue("testpass")

	newModel, cmd := a.updateCreateLogPassSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Empty(t, newApp.createLogPassSecretModel.title.Value())
	assert.Empty(t, newApp.createLogPassSecretModel.metadata.Value())
	assert.Empty(t, newApp.createLogPassSecretModel.login.Value())
	assert.Empty(t, newApp.createLogPassSecretModel.password.Value())
	assert.Nil(t, newApp.createLogPassSecretModel.err)
	assert.NotNil(t, cmd)
}

func TestApp_CreateLogPassSecretView(t *testing.T) {
	a := app{
		createLogPassSecretModel: createLogPassSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			login:    textinput.New(),
			password: textinput.New(),
		},
	}

	a.createLogPassSecretModel.title.SetValue("test-logpass")
	a.createLogPassSecretModel.metadata.SetValue("test-meta")
	a.createLogPassSecretModel.login.SetValue("testuser")
	a.createLogPassSecretModel.password.SetValue("testpass")

	s := a.createLogPassSecretView()

	assert.Contains(t, s, "New secret (type login_password)")
	assert.Contains(t, s, "Title*> test-logpass")
	assert.Contains(t, s, "Metadata> test-meta")
	assert.Contains(t, s, "Login*> testuser")
	assert.Contains(t, s, "Password*> testpass")
	assert.Contains(t, s, "[Ctrl+B] - back in secrets")
}
