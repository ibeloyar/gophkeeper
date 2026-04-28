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

// createTextSecretModel manages UI state for text secret creation form.
// Handles three text inputs (title*, metadata, textData*) with validation.
type createTextSecretModel struct {
	title      textinput.Model
	metadata   textinput.Model
	textData   textinput.Model
	secretType gophkeeperv1.SecretType
	err        error
}

// updateCreateTextSecret handles keyboard input for text secret creation.
// Validates title and textData on Enter, sends gRPC CreateSecret with TEXT type.
// Circular navigation between three fields using Tab/Shift+Tab/Up/Down.
// Clears form and refreshes secrets list after successful creation.
func (a app) updateCreateTextSecret(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if a.createTextSecretModel.title.Value() == "" {
				a.createTextSecretModel.err = errors.New("title is required")
			}
			if a.createTextSecretModel.textData.Value() == "" {
				a.createTextSecretModel.err = errors.New("text data is required")
			}

			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

			_, err := a.client.Cmd.CreateSecret(ctx, &gophkeeperv1.CreateSecretRequest{
				Title:      a.createTextSecretModel.title.Value(),
				Metadata:   a.createTextSecretModel.metadata.Value(),
				SecretType: gophkeeperv1.SecretType_TEXT,
				TextData:   a.createTextSecretModel.textData.Value(),
			})
			if err != nil {
				a.secretsModel.err = err
			}

			a.state = secretsState
			a.createTextSecretModel.title.SetValue("")
			a.createTextSecretModel.metadata.SetValue("")
			a.createTextSecretModel.textData.SetValue("")
			return a, a.pollSecrets()

		case tea.KeyUp:
			fallthrough
		case tea.KeyShiftTab:
			if a.createTextSecretModel.title.Focused() {
				a.createTextSecretModel.title.Blur()
				a.createTextSecretModel.textData.Focus()
			} else if a.createTextSecretModel.metadata.Focused() {
				a.createTextSecretModel.metadata.Blur()
				a.createTextSecretModel.title.Focus()
			} else if a.createTextSecretModel.textData.Focused() {
				a.createTextSecretModel.textData.Blur()
				a.createTextSecretModel.metadata.Focus()
			} else {
				a.createTextSecretModel.title.Focus()
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			if a.createTextSecretModel.title.Focused() {
				a.createTextSecretModel.title.Blur()
				a.createTextSecretModel.metadata.Focus()
			} else if a.createTextSecretModel.metadata.Focused() {
				a.createTextSecretModel.metadata.Blur()
				a.createTextSecretModel.textData.Focus()
			} else if a.createTextSecretModel.textData.Focused() {
				a.createTextSecretModel.textData.Blur()
				a.createTextSecretModel.title.Focus()
			} else {
				a.createTextSecretModel.title.Focus()
			}
		}
	}

	if a.createTextSecretModel.title.Focused() {
		a.createTextSecretModel.title, cmd = a.createTextSecretModel.title.Update(msg)
	}
	if a.createTextSecretModel.metadata.Focused() {
		a.createTextSecretModel.metadata, cmd = a.createTextSecretModel.metadata.Update(msg)
	}
	if a.createTextSecretModel.textData.Focused() {
		a.createTextSecretModel.textData, cmd = a.createTextSecretModel.textData.Update(msg)
	}

	return a, cmd
}

// createTextSecretView renders TUI form for text secret creation.
// Displays required fields (title*, text data*) with validation errors.
// Includes build info, server address, and keyboard shortcuts.
func (a app) createTextSecretView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += "New secret (type text)\n* - required field\n\n"

	s += "Title*" + a.createTextSecretModel.title.View() + "\n"
	s += "Metadata" + a.createTextSecretModel.metadata.View() + "\n"

	s += "Text data*" + a.createTextSecretModel.textData.View() + "\n"

	if a.createTextSecretModel.err != nil {
		s += "\nError: " + a.createTextSecretModel.err.Error() + "\n"
	}

	s += "\n[Ctrl+B] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
