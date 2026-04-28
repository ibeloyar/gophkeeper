package main

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
	"github.com/stretchr/testify/assert"
)

func TestApp_UpdateSecrets_KeyRunes_l(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		token:                  "Bearer test",
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	runMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}

	newModel, cmd := a.updateSecrets(runMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretLogPassCreateState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyRunes_t(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	runMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")}

	newModel, cmd := a.updateSecrets(runMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretTextCreateState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyRunes_b(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	runMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")}

	newModel, cmd := a.updateSecrets(runMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretBinaryCreateState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyRunes_c(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	runMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")}

	newModel, cmd := a.updateSecrets(runMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretCardCreateState, newApp.state)
	assert.Nil(t, newApp.selectedSecret)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyRunes_d(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          1,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	runMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")}

	newModel, cmd := a.updateSecrets(runMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretDeleteState, newApp.state)
	assert.Equal(t, "secret-2", newApp.selectedSecret.Title)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyCtrlC(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateSecrets_KeyEsc(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateSecrets_KeyEnter(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretCardState, newApp.state)
	assert.Equal(t, "secret-1", newApp.selectedSecret.Title)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyDown(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{err: errors.New("test-error")},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyDown}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, 1, newApp.secretsCursor)
	assert.Nil(t, newApp.secretsModel.err)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyDown_WithoutError(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          1,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyDown}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, 1, newApp.secretsCursor)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyUp(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          1,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{err: errors.New("test-error")},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyUp}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, 0, newApp.secretsCursor)
	assert.Nil(t, newApp.secretsModel.err)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_KeyUp_WithoutError(t *testing.T) {
	a := app{
		state:  secretsState,
		client: nil,
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          1,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyUp}

	newModel, cmd := a.updateSecrets(keyMsg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, 0, newApp.secretsCursor)
	assert.Nil(t, cmd)
}

func TestApp_UpdateSecrets_LoadSecretsMsg_Success(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 grpcclient.New(":8080"),
		token:                  "Bearer test",
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{err: errors.New("test-error")},
	}

	msg := loadSecretsMsg{}

	newModel, cmd := a.updateSecrets(msg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.False(t, newApp.secretsInitRequestDone)
	assert.NotNil(t, newApp.secretsModel.err)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateSecrets_LoadSecretsMsg_Failure(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 grpcclient.New(":8080"),
		token:                  "Bearer test",
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsCursor:          0,
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{err: nil},
	}

	msg := loadSecretsMsg{}

	newModel, cmd := a.updateSecrets(msg)

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Len(t, newApp.secrets, 0)
	assert.False(t, newApp.secretsInitRequestDone)
	assert.Error(t, newApp.secretsModel.err)
	assert.Contains(t, newApp.secretsModel.err.Error(), "rpc error")
	assert.NotNil(t, cmd)
}

func TestApp_SecretsView_EmptySecrets_Init(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		buildVersion:           "v1.0.0",
		buildDate:              "2026-03-04",
		serverAddr:             "localhost:8080",
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsInitRequestDone: false,
		secretsModel:           secretsPage{},
	}

	s := a.secretsView()

	assert.Contains(t, s, "Gophkeeper-cli v1.0.0 (2026-03-04)")
	assert.Contains(t, s, "Server: localhost:8080")
	assert.Contains(t, s, "Secrets list:")
	assert.Contains(t, s, "Loading...")
	assert.NotContains(t, s, "Secrets not found")
}

func TestApp_SecretsView_EmptySecrets_RequestDone(t *testing.T) {
	a := app{
		state:                  secretsState,
		client:                 nil,
		buildVersion:           "v1.0.0",
		buildDate:              "2026-03-04",
		serverAddr:             "localhost:8080",
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsInitRequestDone: true,
		secretsModel:           secretsPage{err: nil},
	}

	s := a.secretsView()

	assert.Contains(t, s, "Secrets list:")
	assert.Contains(t, s, "Secrets not found")
	assert.NotContains(t, s, "Loading")
}

func TestApp_SecretsView_SingleSecret(t *testing.T) {
	a := app{
		state:        secretsState,
		client:       nil,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
		},
		secretsCursor:          0,
		secretsInitRequestDone: true,
		secretsModel:           secretsPage{err: nil},
	}

	s := a.secretsView()

	assert.Contains(t, s, "Secrets list:")
	assert.Contains(t, s, "> secret-1")
}

func TestApp_SecretsView_MultipleSecrets(t *testing.T) {
	a := app{
		state:        secretsState,
		client:       nil,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
			{Title: "secret-2"},
		},
		secretsCursor:          1,
		secretsInitRequestDone: true,
		secretsModel:           secretsPage{err: nil},
	}

	s := a.secretsView()

	assert.Contains(t, s, "Secrets list:")
	assert.Contains(t, s, " secret-1")
	assert.Contains(t, s, "> secret-2")
}

func TestApp_SecretsView_Error(t *testing.T) {
	a := app{
		state:        secretsState,
		client:       nil,
		buildVersion: "v1.0.0",
		buildDate:    "2026-03-04",
		serverAddr:   "localhost:8080",
		secrets: []*gophkeeperv1.Secret{
			{Title: "secret-1"},
		},
		secretsCursor:          0,
		secretsInitRequestDone: true,
		secretsModel:           secretsPage{err: errors.New("test-secrets-error")},
	}

	s := a.secretsView()

	assert.Contains(t, s, "Error: test-secrets-error")
}
