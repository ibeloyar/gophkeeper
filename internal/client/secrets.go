package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

// secretsPage manages state for secrets list page in TUI.
// Contains only error field for displaying gRPC operation failures.
type secretsPage struct {
	err error
}

// pollSecrets returns a ticker command that triggers loadSecretsMsg every 3 seconds.
// Enables automatic refresh of secrets list without manual user action.
func (a app) pollSecrets() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return loadSecretsMsg{}
	})
}

// updateSecrets handles keyboard navigation and auto-refresh for secrets list.
// Supports creating new secrets (l,t,b,c), deleting (d), viewing (Enter), navigation (↑↓).
// Auto-refreshes every 3 seconds via loadSecretsMsg. Manages cursor position.
func (a app) updateSecrets(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 'l' {
				a.state = secretLogPassCreateState
				a.selectedSecret = nil
				return a, nil
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 't' {
				a.state = secretTextCreateState
				a.selectedSecret = nil
				return a, nil
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 'b' {
				a.state = secretBinaryCreateState
				a.selectedSecret = nil
				return a, nil
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 'c' {
				a.state = secretCardCreateState
				a.selectedSecret = nil
				return a, nil
			}
			if len(msg.Runes) > 0 && msg.Runes[0] == 'd' {
				a.selectedSecret = a.secrets[a.secretsCursor]
				a.state = secretDeleteState
				return a, nil
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit

		case tea.KeyEnter:
			a.selectedSecret = a.secrets[a.secretsCursor]
			a.state = secretCardState

		case tea.KeyDown:
			a.secretsModel.err = nil

			if a.secretsCursor < len(a.secrets)-1 {
				a.secretsCursor++
			}

		case tea.KeyUp:
			a.secretsModel.err = nil

			if a.secretsCursor > 0 {
				a.secretsCursor--
			}
		}

		return a, nil

	case loadSecretsMsg:
		ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", a.token)

		res, err := a.client.Cmd.GetSecrets(ctx, &gophkeeperv1.GetSecretsRequest{})
		if err != nil {
			a.secretsModel.err = err
		} else {
			a.secretsModel.err = nil
			a.secrets = res.Secrets
			a.secretsInitRequestDone = true
		}

		return a, a.pollSecrets()
	}

	return a, nil
}

// secretsView renders paginated secrets list with cursor navigation.
// Shows loading state, empty state, errors, and comprehensive keyboard shortcuts.
// Displays build info, server address, and secret titles with cursor indicator.
func (a app) secretsView() string {
	s := fmt.Sprintf("Gophkeeper-cli %s (%s) Server: %s \n==================================================\n",
		a.buildVersion, a.buildDate, a.serverAddr,
	)
	s += "Secrets list:\n\n"

	for i, choice := range a.secrets {
		cursor := " "
		if a.secretsCursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice.Title)
	}

	if len(a.secrets) < 1 && a.secretsInitRequestDone {
		s += "Secrets not found\n"
	}
	if len(a.secrets) < 1 && !a.secretsInitRequestDone {
		s += "Loading...\n"
	}

	if a.secretsModel.err != nil {
		s += "\nError: " + a.secretsModel.err.Error() + "\n"
	}

	s += "\nAdd secret: [l] - login/password, [t] - text, [b] - binary, [c] - card\n"
	s += "Control:    [Tab(Shift+Tab) or ↑↓] - move, [d] - delete, [Enter] - show, [Ctrl+C] – quit \n"

	return s
}
