package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type createBinarySecretModel struct {
	title      textinput.Model
	metadata   textinput.Model
	filePath   textinput.Model
	secretType gophkeeperv1.SecretType
	err        error
}

func (a app) updateCreateBinarySecret(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if a.createBinarySecretModel.title.Value() == "" {
				a.createBinarySecretModel.err = errors.New("title is required")
			}
			if a.createBinarySecretModel.filePath.Value() == "" {
				a.createBinarySecretModel.err = errors.New("file path is required")
			}

			f, err := os.Open(a.createBinarySecretModel.filePath.Value())
			if err != nil {
				a.createBinarySecretModel.err = fmt.Errorf("failed to open file: %w", err)
				return a, cmd
			}
			defer f.Close()

			stat, err := f.Stat()
			if err != nil {
				a.createBinarySecretModel.err = fmt.Errorf("failed to stat file: %w", err)
				return a, cmd
			}

			data := make([]byte, stat.Size())
			_, err = f.Read(data)
			if err != nil {
				a.createBinarySecretModel.err = fmt.Errorf("failed to read file: %w", err)
				return a, cmd
			}

			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

			_, err = a.client.Cmd.CreateSecret(ctx, &gophkeeperv1.CreateSecretRequest{
				Title:      a.createBinarySecretModel.title.Value(),
				Metadata:   a.createBinarySecretModel.metadata.Value(),
				SecretType: gophkeeperv1.SecretType_BINARY,
				BinaryData: data,
			})
			if err != nil {
				a.secretsModel.err = err
			}

			a.state = secretsState
			a.createBinarySecretModel.title.SetValue("")
			a.createBinarySecretModel.metadata.SetValue("")
			a.createBinarySecretModel.filePath.SetValue("")
			return a, a.pollSecrets()

		case tea.KeyUp:
			fallthrough
		case tea.KeyShiftTab:
			if a.createBinarySecretModel.title.Focused() {
				a.createBinarySecretModel.title.Blur()
				a.createBinarySecretModel.filePath.Focus()
			} else if a.createBinarySecretModel.metadata.Focused() {
				a.createBinarySecretModel.metadata.Blur()
				a.createBinarySecretModel.title.Focus()
			} else if a.createBinarySecretModel.filePath.Focused() {
				a.createBinarySecretModel.filePath.Blur()
				a.createBinarySecretModel.metadata.Focus()
			} else {
				a.createBinarySecretModel.title.Focus()
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			if a.createBinarySecretModel.title.Focused() {
				a.createBinarySecretModel.title.Blur()
				a.createBinarySecretModel.metadata.Focus()
			} else if a.createBinarySecretModel.metadata.Focused() {
				a.createBinarySecretModel.metadata.Blur()
				a.createBinarySecretModel.filePath.Focus()
			} else if a.createBinarySecretModel.filePath.Focused() {
				a.createBinarySecretModel.filePath.Blur()
				a.createBinarySecretModel.title.Focus()
			} else {
				a.createBinarySecretModel.title.Focus()
			}
		}
	}

	if a.createBinarySecretModel.title.Focused() {
		a.createBinarySecretModel.title, cmd = a.createBinarySecretModel.title.Update(msg)
	}
	if a.createBinarySecretModel.metadata.Focused() {
		a.createBinarySecretModel.metadata, cmd = a.createBinarySecretModel.metadata.Update(msg)
	}
	if a.createBinarySecretModel.filePath.Focused() {
		a.createBinarySecretModel.filePath, cmd = a.createBinarySecretModel.filePath.Update(msg)
	}

	return a, cmd
}

func (a app) createBinarySecretView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += "New secret (type binary)\n* - required field\n\n"

	s += "Title*" + a.createBinarySecretModel.title.View() + "\n"
	s += "Metadata" + a.createBinarySecretModel.metadata.View() + "\n"
	s += "File path*" + a.createBinarySecretModel.filePath.View() + "\n"

	if a.createBinarySecretModel.err != nil {
		s += "\nError: " + a.createBinarySecretModel.err.Error() + "\n"
	}
	s += "\n[Tab(Shift+Tab) or ↑↓] - move, [Ctrl+B] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
