package main

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"
)

func TestApp_UpdateLogin_CtrlC(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        loginState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, loginState, newApp.state)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateLogin_CtrlR(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.SetValue("testuser")
	a.loginModel.passwordInput.SetValue("testpass")
	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyCtrlR})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, registerState, newApp.state)
	assert.Empty(t, newApp.loginModel.loginInput.Value())
	assert.Empty(t, newApp.loginModel.passwordInput.Value())
	assert.Nil(t, newApp.loginModel.err)
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_RequiredLogin(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.passwordInput.SetValue("testpass")

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, loginState, newApp.state)
	assert.Contains(t, newApp.loginModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_RequiredPassword(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.SetValue("testuser")

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, loginState, newApp.state)
	assert.Contains(t, newApp.loginModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_RequiredBoth(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, loginState, newApp.state)
	assert.Contains(t, newApp.loginModel.err.Error(), "login and password are required")
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_Navigation_Up(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyUp})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.loginModel.loginInput.Focused())
	assert.True(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_Navigation_Down(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyDown})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.loginModel.loginInput.Focused())
	assert.True(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)

	newModel, cmd = newApp.updateLogin(tea.KeyMsg{Type: tea.KeyDown})
	newApp, ok = newModel.(app)
	assert.True(t, ok)

	assert.True(t, newApp.loginModel.loginInput.Focused())
	assert.False(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_Navigation_Tab(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyTab})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.loginModel.loginInput.Focused())
	assert.True(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)

	newModel, cmd = newApp.updateLogin(tea.KeyMsg{Type: tea.KeyTab})
	newApp, ok = newModel.(app)
	assert.True(t, ok)

	assert.True(t, newApp.loginModel.loginInput.Focused())
	assert.False(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_UpdateLogin_Navigation_ShiftTab(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  loginState,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.Focus()

	newModel, cmd := a.updateLogin(tea.KeyMsg{Type: tea.KeyShiftTab})
	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.loginModel.loginInput.Focused())
	assert.True(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)

	newModel, cmd = newApp.updateLogin(tea.KeyMsg{Type: tea.KeyShiftTab})
	newApp, ok = newModel.(app)
	assert.True(t, ok)

	assert.True(t, newApp.loginModel.loginInput.Focused())
	assert.False(t, newApp.loginModel.passwordInput.Focused())
	assert.Nil(t, cmd)
}

func TestApp_LoginView(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        loginState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	a.loginModel.loginInput.SetValue("testuser")
	a.loginModel.passwordInput.SetValue("testpass")

	s := a.loginView()

	assert.Contains(t, s, "Gophkeeper-cli v1.0.0 (2026-03-04)")
	assert.Contains(t, s, "Server: localhost:8080")
	assert.Contains(t, s, "Login")
	assert.Contains(t, s, "Login* > testuser")
	assert.Contains(t, s, "Password* > testpass")
	assert.Contains(t, s, "[Enter] – login, [Ctrl+R] – registration, [Ctrl+C] – quit")
}

func TestApp_LoginView_Error(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        loginState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
			err:           nil,
		},
	}

	a.loginModel.err = errors.New("test login error")

	s := a.loginView()

	assert.Contains(t, s, "Error: test login error")
}

func TestApp_LoginView_EmptyFields(t *testing.T) {
	a := app{
		client:       grpcclient.New(":8080"),
		state:        loginState,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
	}

	s := a.loginView()

	assert.Contains(t, s, "Login")
	assert.Contains(t, s, "Password")
	assert.NotContains(t, s, "Login*>")
	assert.NotContains(t, s, "Password*>")
}
