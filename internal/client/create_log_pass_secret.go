package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"google.golang.org/grpc/metadata"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

// createLogPassSecretModel manages state for login/password secret creation form.
// Handles four text inputs (title*, metadata, login*, password*) with validation.
type createLogPassSecretModel struct {
	title      textinput.Model
	metadata   textinput.Model
	login      textinput.Model
	password   textinput.Model
	secretType gophkeeperv1.SecretType
	err        error
}

// updateCreateLogPassSecret processes input events for login/password secret form.
// Validates required fields on Enter, sends gRPC CreateSecret with LOGIN_PASSWORD type.
// Supports field navigation with Tab/Shift+Tab/Up/Down, returns to secrets list on success.
func (a app) updateCreateLogPassSecret(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyCtrlB:
			a.state = secretsState
			a.selectedSecret = nil
			return a, nil

		case tea.KeyEnter:
			if a.createLogPassSecretModel.title.Value() == "" {
				a.createLogPassSecretModel.err = errors.New("title is required")
			}
			if a.createLogPassSecretModel.login.Value() == "" {
				a.createLogPassSecretModel.err = errors.New("login is required")
			}
			if a.createLogPassSecretModel.password.Value() == "" {
				a.createLogPassSecretModel.err = errors.New("password is required")
			}

			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

			_, err := a.client.Cmd.CreateSecret(ctx, &gophkeeperv1.CreateSecretRequest{
				Title:      a.createLogPassSecretModel.title.Value(),
				Metadata:   a.createLogPassSecretModel.metadata.Value(),
				SecretType: gophkeeperv1.SecretType_LOGIN_PASSWORD,
				Login:      a.createLogPassSecretModel.login.Value(),
				Password:   a.createLogPassSecretModel.password.Value(),
			})

			if err != nil {
				a.secretsModel.err = err
			}

			a.state = secretsState
			a.createLogPassSecretModel.title.SetValue("")
			a.createLogPassSecretModel.metadata.SetValue("")
			a.createLogPassSecretModel.login.SetValue("")
			a.createLogPassSecretModel.password.SetValue("")
			return a, a.pollSecrets()

		case tea.KeyUp:
			fallthrough
		case tea.KeyShiftTab:
			if a.createLogPassSecretModel.title.Focused() {
				a.createLogPassSecretModel.title.Blur()
				a.createLogPassSecretModel.password.Focus()
			} else if a.createLogPassSecretModel.metadata.Focused() {
				a.createLogPassSecretModel.metadata.Blur()
				a.createLogPassSecretModel.title.Focus()
			} else if a.createLogPassSecretModel.login.Focused() {
				a.createLogPassSecretModel.login.Blur()
				a.createLogPassSecretModel.metadata.Focus()
			} else if a.createLogPassSecretModel.password.Focused() {
				a.createLogPassSecretModel.password.Blur()
				a.createLogPassSecretModel.login.Focus()
			} else {
				a.createLogPassSecretModel.title.Focus()
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			if a.createLogPassSecretModel.title.Focused() {
				a.createLogPassSecretModel.title.Blur()
				a.createLogPassSecretModel.metadata.Focus()
			} else if a.createLogPassSecretModel.metadata.Focused() {
				a.createLogPassSecretModel.metadata.Blur()
				a.createLogPassSecretModel.login.Focus()
			} else if a.createLogPassSecretModel.login.Focused() {
				a.createLogPassSecretModel.login.Blur()
				a.createLogPassSecretModel.password.Focus()
			} else if a.createLogPassSecretModel.password.Focused() {
				a.createLogPassSecretModel.password.Blur()
				a.createLogPassSecretModel.title.Focus()
			} else {
				a.createLogPassSecretModel.title.Focus()
			}
		}
	}

	if a.createLogPassSecretModel.title.Focused() {
		a.createLogPassSecretModel.title, cmd = a.createLogPassSecretModel.title.Update(msg)
	}
	if a.createLogPassSecretModel.metadata.Focused() {
		a.createLogPassSecretModel.metadata, cmd = a.createLogPassSecretModel.metadata.Update(msg)
	}
	if a.createLogPassSecretModel.login.Focused() {
		a.createLogPassSecretModel.login, cmd = a.createLogPassSecretModel.login.Update(msg)
	}
	if a.createLogPassSecretModel.password.Focused() {
		a.createLogPassSecretModel.password, cmd = a.createLogPassSecretModel.password.Update(msg)
	}

	return a, cmd
}

// createLogPassSecretView renders form for login/password secret creation.
// Shows required fields (title*, login*, password*) with validation errors.
// Displays navigation help and build/server information.
func (a app) createLogPassSecretView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += "New secret (type login_password)\n* - required field\n\n"

	s += "Title*" + a.createLogPassSecretModel.title.View() + "\n"
	s += "Metadata" + a.createLogPassSecretModel.metadata.View() + "\n"

	s += "Login*" + a.createLogPassSecretModel.login.View() + "\n"
	s += "Password*" + a.createLogPassSecretModel.password.View() + "\n"

	if a.createLogPassSecretModel.err != nil {
		s += "\nError: " + a.createLogPassSecretModel.err.Error() + "\n"
	}

	s += "\n[Ctrl+B] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
