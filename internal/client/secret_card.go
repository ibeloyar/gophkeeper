package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func (a app) updateSecretCard(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyCtrlB:
			a.state = secretsState
			a.selectedSecret = nil
			return a, nil
		}
	}

	return a, nil
}

func (a app) secretCardView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += fmt.Sprintf("Secret: %s\n\n", a.selectedSecret.Title)

	if a.selectedSecret.SecretType == gophkeeperv1.SecretType_LOGIN_PASSWORD {
		s += fmt.Sprintf("Login: %s\n", a.selectedSecret.Login)
		s += fmt.Sprintf("Passowrd: %s\n", a.selectedSecret.Password)
	}

	if a.selectedSecret.SecretType == gophkeeperv1.SecretType_TEXT {
		s += fmt.Sprintf("TextData: %s\n", a.selectedSecret.TextData)
	}

	if a.selectedSecret.SecretType == gophkeeperv1.SecretType_BINARY {
		s += fmt.Sprintf("BinaryData: file with size %d bytes\n", len(a.selectedSecret.BinaryData))
	}

	if a.selectedSecret.SecretType == gophkeeperv1.SecretType_CARD {
		s += fmt.Sprintf("Card holder: %s\n", a.selectedSecret.CardHolder)
		s += fmt.Sprintf("Card number: %s\n", a.selectedSecret.CardNumber)
		s += fmt.Sprintf("Card exp: %s\n", a.selectedSecret.CardExp)
	}

	s += "\n[Ctrl+B] – back in secrets, [Ctrl+C] – quit\n"

	return s
}
