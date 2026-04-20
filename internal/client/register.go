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

type registerForm struct {
	loginInput    textinput.Model
	passwordInput textinput.Model
	err           error
}

func (a app) updateRegister(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyCtrlB:
			a.state = loginState
			a.registerModel.loginInput.SetValue("")
			a.registerModel.passwordInput.SetValue("")
			a.registerModel.err = nil
			return a, nil

		case tea.KeyEnter:
			if a.registerModel.loginInput.Value() != "" && a.registerModel.passwordInput.Value() != "" {
				var header metadata.MD

				_, err := a.client.Cmd.Register(context.Background(), &gophkeeperv1.RegisterRequest{
					Login:    a.registerModel.loginInput.Value(),
					Password: a.registerModel.passwordInput.Value(),
				}, grpc.Header(&header))

				if err != nil {
					a.registerModel.err = err
				} else {
					headersList := header.Get("token")

					if len(headersList) == 0 {
						a.registerModel.err = fmt.Errorf("token not found")
					}

					a.state = secretsState
					a.registerModel.err = nil
					a.token = headersList[0]
					return a, a.pollSecrets()
				}
			} else {
				a.registerModel.err = fmt.Errorf("login and password are required")
			}

		case tea.KeyUp:
			fallthrough
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			if a.registerModel.loginInput.Focused() {
				a.registerModel.loginInput.Blur()
				a.registerModel.passwordInput.Focus()
			} else {
				a.registerModel.passwordInput.Blur()
				a.registerModel.loginInput.Focus()
			}
		}
	}

	if a.registerModel.loginInput.Focused() {
		a.registerModel.loginInput, cmd = a.registerModel.loginInput.Update(msg)
	} else {
		a.registerModel.passwordInput, cmd = a.registerModel.passwordInput.Update(msg)
	}

	return a, cmd
}

func (a app) registerView() string {
	s := "Registration\n* - required field\n\n"

	s += "Login* " + a.registerModel.loginInput.View() + "\n"
	s += "Password* " + a.registerModel.passwordInput.View() + "\n"

	if a.registerModel.err != nil {
		s += "\nError: " + a.registerModel.err.Error() + "\n"
	}

	s += "\n\nIf you already have an account, you can login\n"

	s += "\n> [Enter] – register, [Ctrl+B] – back to login, [Ctrl+C] – quit"

	return s
}
