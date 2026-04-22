package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

// loginForm manages state for user login form in TUI.
// Contains two text inputs (login*, password*) with error display.
type loginForm struct {
	loginInput    textinput.Model
	passwordInput textinput.Model
	err           error
}

// updateLogin handles keyboard events for login form.
// Performs gRPC Login RPC on Enter (with both fields filled), extracts token from headers.
// Supports Tab navigation between fields, Ctrl+R to switch to registration.
func (a app) updateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlR:
			a.state = registerState
			a.loginModel.loginInput.SetValue("")
			a.loginModel.passwordInput.SetValue("")
			a.loginModel.err = nil

		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyEnter:
			if a.loginModel.loginInput.Value() != "" && a.loginModel.passwordInput.Value() != "" {
				var header metadata.MD

				_, err := a.client.Cmd.Login(context.Background(), &gophkeeperv1.LoginRequest{
					Login:    a.loginModel.loginInput.Value(),
					Password: a.loginModel.passwordInput.Value(),
				}, grpc.Header(&header))

				if err != nil {
					a.loginModel.err = err
				} else {
					headersList := header.Get("token")

					if len(headersList) == 0 {
						a.loginModel.err = fmt.Errorf("token not found")
					}

					a.state = secretsState
					a.loginModel.err = nil
					a.token = headersList[0]
					return a, a.pollSecrets()
				}
			} else {
				a.loginModel.err = fmt.Errorf("login and password are required")
			}

		case tea.KeyUp:
			fallthrough
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			fallthrough
		case tea.KeyShiftTab:
			if a.loginModel.loginInput.Focused() {
				a.loginModel.loginInput.Blur()
				a.loginModel.passwordInput.Focus()
			} else {
				a.loginModel.passwordInput.Blur()
				a.loginModel.loginInput.Focus()
			}
		}
	}

	if a.loginModel.loginInput.Focused() {
		a.loginModel.loginInput, cmd = a.loginModel.loginInput.Update(msg)
	} else {
		a.loginModel.passwordInput, cmd = a.loginModel.passwordInput.Update(msg)
	}

	return a, cmd
}

// loginView renders login form UI with input fields, validation errors, and navigation.
// Displays build info, server address, and registration prompt.
func (a app) loginView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += "Login\n* - required field\n\n"

	s += "Login* " + a.loginModel.loginInput.View() + "\n"
	s += "Password* " + a.loginModel.passwordInput.View() + "\n"

	if a.loginModel.err != nil {
		s += "\nError: " + a.loginModel.err.Error() + "\n"
	}

	s += "\n\nIf you don't have an account, go to registration\n"

	s += "\n> [Enter] – login, [Ctrl+R] – registration, [Ctrl+C] – quit"

	return s
}
