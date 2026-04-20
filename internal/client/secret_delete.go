package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
	"google.golang.org/grpc/metadata"
)

func (a app) updateSecretDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 'y' {
				ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

				_, err := a.client.Cmd.DeleteSecret(ctx, &gophkeeperv1.DeleteSecretRequest{
					Title: a.selectedSecret.Title,
				})
				if err != nil {
					a.secretsModel.err = fmt.Errorf("failed to delete secret %s: %w", a.selectedSecret.Title, err)
				}

				a.state = secretsState
				a.selectedSecret = nil
				return a, nil
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 'n' {
				a.state = secretsState
				return a, nil
			}
		}
	}

	return a, nil
}

func (a app) secretDeleteView() string {
	s := "Delete secret\n\n"

	s += fmt.Sprintf("Are you sure you want to delete the secret %s?\n\n", a.selectedSecret.Title)

	s += "[y] - YES delete secret, [n] - NO back in secrets,  [Ctrl+C] – quit\n"

	return s
}
