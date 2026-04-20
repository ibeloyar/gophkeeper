package main

import (
	"github.com/charmbracelet/bubbles/textinput"

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
		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 'b' {
				a.state = secretsState
				a.selectedSecret = nil
				return a, nil
			}

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
	s := "New secret (type binary)\n* - required field\n\n"

	s += "Title*" + a.createBinarySecretModel.title.View() + "\n"
	s += "Metadata" + a.createBinarySecretModel.metadata.View() + "\n"
	s += "File path*" + a.createBinarySecretModel.filePath.View() + "\n"

	if a.createBinarySecretModel.err != nil {
		s += "\nError: " + a.createBinarySecretModel.err.Error() + "\n"
	}
	s += "\n[Tab(Shift+Tab) or ↑↓] - move, [b] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
