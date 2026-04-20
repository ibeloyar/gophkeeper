package main

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/textinput"
	"google.golang.org/grpc/metadata"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type createCardSecretModel struct {
	title      textinput.Model
	metadata   textinput.Model
	secretType gophkeeperv1.SecretType
	cardExp    textinput.Model
	cardHolder textinput.Model
	cardNumber textinput.Model
	err        error
}

func (a app) updateCreateCardSecret(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlB:
			a.state = secretsState
			a.selectedSecret = nil
			return a, nil

		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyEnter:
			if a.createCardSecretModel.title.Value() == "" {
				a.createCardSecretModel.err = errors.New("title is required")
			}
			if a.createCardSecretModel.cardExp.Value() == "" {
				a.createCardSecretModel.err = errors.New("card exp is required")
			}
			if a.createCardSecretModel.cardNumber.Value() == "" {
				a.createCardSecretModel.err = errors.New("card number is required")
			}
			if a.createCardSecretModel.cardHolder.Value() == "" {
				a.createCardSecretModel.err = errors.New("card holder is required")
			}

			ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

			_, err := a.client.Cmd.CreateSecret(ctx, &gophkeeperv1.CreateSecretRequest{
				Title:      a.createCardSecretModel.title.Value(),
				Metadata:   a.createCardSecretModel.metadata.Value(),
				SecretType: gophkeeperv1.SecretType_CARD,
				CardExp:    a.createCardSecretModel.cardExp.Value(),
				CardHolder: a.createCardSecretModel.cardHolder.Value(),
				CardNumber: a.createCardSecretModel.cardNumber.Value(),
			})

			if err != nil {
				a.secretsModel.err = err
			}

			a.state = secretsState
			a.createCardSecretModel.title.SetValue("")
			a.createCardSecretModel.metadata.SetValue("")
			a.createCardSecretModel.cardExp.SetValue("")
			a.createCardSecretModel.cardHolder.SetValue("")
			a.createCardSecretModel.cardNumber.SetValue("")
			return a, a.pollSecrets()

		case tea.KeyUp:
			fallthrough
		case tea.KeyShiftTab:
			if a.createCardSecretModel.title.Focused() {
				a.createCardSecretModel.title.Blur()
				a.createCardSecretModel.cardExp.Focus()
			} else if a.createCardSecretModel.metadata.Focused() {
				a.createCardSecretModel.metadata.Blur()
				a.createCardSecretModel.title.Focus()
			} else if a.createCardSecretModel.cardNumber.Focused() {
				a.createCardSecretModel.cardNumber.Blur()
				a.createCardSecretModel.metadata.Focus()
			} else if a.createCardSecretModel.cardHolder.Focused() {
				a.createCardSecretModel.cardHolder.Blur()
				a.createCardSecretModel.cardNumber.Focus()
			} else if a.createCardSecretModel.cardExp.Focused() {
				a.createCardSecretModel.cardExp.Blur()
				a.createCardSecretModel.cardHolder.Focus()
			} else {
				a.createCardSecretModel.title.Focus()
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyTab:
			if a.createCardSecretModel.title.Focused() {
				a.createCardSecretModel.title.Blur()
				a.createCardSecretModel.metadata.Focus()
			} else if a.createCardSecretModel.metadata.Focused() {
				a.createCardSecretModel.metadata.Blur()
				a.createCardSecretModel.cardNumber.Focus()
			} else if a.createCardSecretModel.cardNumber.Focused() {
				a.createCardSecretModel.cardNumber.Blur()
				a.createCardSecretModel.cardHolder.Focus()
			} else if a.createCardSecretModel.cardHolder.Focused() {
				a.createCardSecretModel.cardHolder.Blur()
				a.createCardSecretModel.cardExp.Focus()
			} else if a.createCardSecretModel.cardExp.Focused() {
				a.createCardSecretModel.cardExp.Blur()
				a.createCardSecretModel.title.Focus()
			} else {
				a.createCardSecretModel.title.Focus()
			}
		}
	}

	if a.createCardSecretModel.title.Focused() {
		a.createCardSecretModel.title, cmd = a.createCardSecretModel.title.Update(msg)
	}
	if a.createCardSecretModel.metadata.Focused() {
		a.createCardSecretModel.metadata, cmd = a.createCardSecretModel.metadata.Update(msg)
	}
	if a.createCardSecretModel.cardExp.Focused() {
		a.createCardSecretModel.cardExp, cmd = a.createCardSecretModel.cardExp.Update(msg)
	}
	if a.createCardSecretModel.cardNumber.Focused() {
		a.createCardSecretModel.cardNumber, cmd = a.createCardSecretModel.cardNumber.Update(msg)
	}
	if a.createCardSecretModel.cardHolder.Focused() {
		a.createCardSecretModel.cardHolder, cmd = a.createCardSecretModel.cardHolder.Update(msg)
	}

	return a, cmd
}

func (a app) createCardSecretView() string {
	s := "New secret (type card)\n* - required field\n\n"

	s += "Title*" + a.createCardSecretModel.title.View() + "\n"
	s += "Metadata" + a.createCardSecretModel.metadata.View() + "\n"
	s += "Card number*" + a.createCardSecretModel.cardNumber.View() + "\n"
	s += "Card holder*" + a.createCardSecretModel.cardHolder.View() + "\n"
	s += "Card exp*" + a.createCardSecretModel.cardExp.View() + "\n"

	if a.createCardSecretModel.err != nil {
		s += "\nError: " + a.createCardSecretModel.err.Error() + "\n"
	}

	s += "\n[Ctrl+B] - back in secrets, [Enter] - create, [Ctrl+C] – quit \n"

	return s
}
