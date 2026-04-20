package main

import (
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type createTextSecretModel struct {
	title      textinput.Model
	metadata   textinput.Model
	textData   textinput.Model
	secretType gophkeeperv1.SecretType
	err        error
}

func (a app) updateCreateTextSecret(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (a app) createTextSecretView() string {
	s := "New secret (type text)\n* - required field\n\n"

	s += "Title*" + a.createTextSecretModel.title.View() + "\n"
	s += "Metadata" + a.createTextSecretModel.metadata.View() + "\n"

	s += "Text data*" + a.createTextSecretModel.textData.View() + "\n"

	if a.createTextSecretModel.err != nil {
		s += "\nError: " + a.createTextSecretModel.err.Error() + "\n"
	}

	s += "\n[b] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
