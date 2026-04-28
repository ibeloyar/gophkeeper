package main

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"
)

func TestApp_UpdateRegister_CtrlC(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        registerState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.Focus()

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, registerState, newApp.state)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateRegister_CtrlB(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.SetValue("testuser")
	a.registerModel.passwordInput.SetValue("testpass")
	a.registerModel.loginInput.Focus()

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, loginState, newApp.state)
	assert.Empty(t, newApp.registerModel.loginInput.Value())
	assert.Empty(t, newApp.registerModel.passwordInput.Value())
	assert.Nil(t, newApp.registerModel.err)
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_RequiredLogin(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.passwordInput.SetValue("testpass")

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, registerState, newApp.state)
	assert.Contains(t, newApp.registerModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_RequiredPassword(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.SetValue("testuser")

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, registerState, newApp.state)
	assert.Contains(t, newApp.registerModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_RequiredBoth(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, registerState, newApp.state)
	assert.Contains(t, newApp.registerModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_Navigation_Up(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.Focus()

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyUp})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.registerModel.loginInput.Focused())
	assert.True(t, newApp.registerModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_Navigation_Down(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.Focus()

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyDown})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.registerModel.loginInput.Focused())
	assert.True(t, newApp.registerModel.passwordInput.Focused())
	assert.Nil(t, cmd)

	newModel, cmd = newApp.updateRegister(tea.KeyMsg{Type: tea.KeyDown})
	newApp, ok = newModel.(app)
	assert.True(t, ok)

	assert.True(t, newApp.registerModel.loginInput.Focused())
	assert.False(t, newApp.registerModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_UpdateRegister_Navigation_Tab(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  registerState,
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.Focus()

	newModel, cmd := a.updateRegister(tea.KeyMsg{Type: tea.KeyTab})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.registerModel.loginInput.Focused())
	assert.True(t, newApp.registerModel.passwordInput.Focused())
	assert.Nil(t, cmd)

	newModel, cmd = newApp.updateRegister(tea.KeyMsg{Type: tea.KeyTab})
	newApp, ok = newModel.(app)
	assert.True(t, ok)

	assert.True(t, newApp.registerModel.loginInput.Focused())
	assert.False(t, newApp.registerModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_RegisterView(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        registerState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.loginInput.SetValue("testuser")
	a.registerModel.passwordInput.SetValue("testpass")

	s := a.registerView()

	assert.Contains(t, s, "Gophkeeper-cli v1.0.0 (2026-03-04)")
	assert.Contains(t, s, "Server: localhost:8080")
	assert.Contains(t, s, "Registration")
	assert.Contains(t, s, "Login* > testuser")
	assert.Contains(t, s, "Password* > testpass")
	assert.Contains(t, s, "[Enter] – register, [Ctrl+B] – back to login, [Ctrl+C] – quit")
}

func TestApp_RegisterView_Error(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        registerState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.registerModel.err = errors.New("test register error")

	s := a.registerView()

	assert.Contains(t, s, "Error: test register error")
}

func TestApp_RegisterView_EmptyFields(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        registerState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	s := a.registerView()

	assert.Contains(t, s, "Registration")
	assert.Contains(t, s, "Login")
	assert.Contains(t, s, "Password")
	assert.NotContains(t, s, "Login*>")
	assert.NotContains(t, s, "Password*>")
}
