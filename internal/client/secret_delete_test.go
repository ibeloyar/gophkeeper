package main

import (
	"context"
	"testing"

	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type GRPCClient struct {
	client gophkeeperv1.GophkeeperClient
}

func (g *GRPCClient) DeleteSecret(ctx context.Context, req *gophkeeperv1.DeleteSecretRequest) (*gophkeeperv1.DeleteSecretResponse, error) {
	return g.client.DeleteSecret(ctx, req)
}

// Тестирование updateSecretDelete

func TestApp_UpdateSecretDelete_CtrlC(t *testing.T) {
	client := grpcclient.New(":8080")

	m := app{
		state:          secretDeleteState,
		token:          "Bearer test",
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
		client:         client,
	}

	newModel, cmd := m.updateSecretDelete(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, m.state, newApp.state)
	assert.Equal(t, m.selectedSecret.Title, newApp.selectedSecret.Title)
	assert.Equal(t, m.token, newApp.token)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateSecretDelete_ConfirmYes(t *testing.T) {
	client := grpcclient.New(":8080")

	m := app{
		state:          secretCardState,
		token:          "Bearer test",
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
		client:         client,
	}

	newModel, cmd := m.updateSecretDelete(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// After y: the state transitions to secretsState, selectedSecret = nil
	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecretDelete_ConfirmNo(t *testing.T) {
	m := app{
		state:          secretDeleteState,
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	newModel, cmd := m.updateSecretDelete(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// After n: state goes to secretsState, selectedSecret = nil, no commands
	assert.Equal(t, secretsState, newApp.state)
	assert.NotNil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecretDelete_OtherKey(t *testing.T) {
	m := app{
		state:          secretCardState,
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	newModel, cmd := m.updateSecretDelete(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	// Не y / n / Ctrl+C -> состояние не меняется, команд нет
	assert.Equal(t, m.state, newApp.state)
	assert.Equal(t, m.selectedSecret.Title, newApp.selectedSecret.Title)
	assert.Nil(t, cmd)
}

// Тестирование secretDeleteView

func TestApp_SecretDeleteView(t *testing.T) {
	m := app{
		selectedSecret: &gophkeeperv1.Secret{Title: "test-secret"},
	}

	output := m.secretDeleteView()

	assert.Contains(t, output, "Delete secret")
	assert.Contains(t, output, "Are you sure you want to delete the secret test-secret?")
	assert.Contains(t, output, "[y] - YES delete secret, [n] - NO back in secrets,  [Ctrl+C] – quit")
}
