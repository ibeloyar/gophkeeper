package main

import (
	"os"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func TestApp_UpdateCreateBinarySecret_CtrlC(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretBinaryCreateState,
		token:  "Bearer test",
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
	}

	a.createBinarySecretModel.title.Focus()

	newModel, cmd := a.updateCreateBinarySecret(tea.KeyMsg{Type: tea.KeyCtrlC})

	assert.NotNil(t, cmd) // tea.Quit

	newApp, ok := newModel.(app)
	assert.True(t, ok)
	assert.Equal(t, a.state, newApp.state)
}

func TestApp_UpdateCreateBinarySecret_CtrlB(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretBinaryCreateState,
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
	}

	a.createBinarySecretModel.title.Focus()

	newModel, cmd := a.updateCreateBinarySecret(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateBinarySecret_TitleRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretBinaryCreateState,
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
	}

	a.createBinarySecretModel.filePath.SetValue("/fake/path")

	newModel, cmd := a.updateCreateBinarySecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretBinaryCreateState, newApp.state)
	assert.Contains(t, newApp.createBinarySecretModel.err.Error(), "failed to open file")
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateBinarySecret_NoSuchFile(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretBinaryCreateState,
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
	}

	a.createBinarySecretModel.title.SetValue("test-title")

	newModel, cmd := a.updateCreateBinarySecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretBinaryCreateState, newApp.state)
	assert.Contains(t, newApp.createBinarySecretModel.err.Error(), "no such file or directory")
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateBinarySecret_ValidInput(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.Write([]byte("test data"))
	assert.NoError(t, err)

	a := app{
		client: grpcclient.New(":8080"),
		state:  secretBinaryCreateState,
		token:  "Bearer test",
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
		secretsModel: secretsPage{},
		secrets:      []*gophkeeperv1.Secret{},
	}

	a.createBinarySecretModel.title.SetValue("test-title")
	a.createBinarySecretModel.metadata.SetValue("test-meta")
	a.createBinarySecretModel.filePath.SetValue(tmpFile.Name())

	newModel, cmd := a.updateCreateBinarySecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Empty(t, newApp.createBinarySecretModel.title.Value())
	assert.Empty(t, newApp.createBinarySecretModel.metadata.Value())
	assert.Empty(t, newApp.createBinarySecretModel.filePath.Value())
	assert.Nil(t, newApp.createBinarySecretModel.err)
	assert.NotNil(t, cmd)
}

func TestApp_CreateBinarySecretView(t *testing.T) {
	a := app{
		createBinarySecretModel: createBinarySecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			filePath: textinput.New(),
		},
	}

	a.createBinarySecretModel.title.SetValue("test-title")
	a.createBinarySecretModel.filePath.SetValue("/fake/path")

	s := a.createBinarySecretView()

	assert.Contains(t, s, "New secret (type binary)")
	assert.Contains(t, s, "Title*> test-title")
	assert.Contains(t, s, "File path*> /fake/path")
	assert.Contains(t, s, "[Tab(Shift+Tab) or ↑↓] - move")
}
