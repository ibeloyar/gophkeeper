package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"
)

func TestApp_UpdateCreateTextSecret_CtrlC(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		token:  "Bearer test",
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()

	newModel, cmd := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyCtrlC})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, a.state, newApp.state)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateTextSecret_CtrlB(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()

	newModel, cmd := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyCtrlB})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Nil(t, cmd)
}

func TestApp_UpdateCreateTextSecret_TitleRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.textData.SetValue("test text")

	newModel, cmd := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createTextSecretModel.err.Error(), "title is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateTextSecret_TextDataRequired(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.SetValue("test-title")

	newModel, cmd := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Contains(t, newApp.createTextSecretModel.err.Error(), "text data is required")
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateTextSecret_ValidInput(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		token:  "Bearer test",
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
		secretsModel: secretsPage{},
	}

	a.createTextSecretModel.title.SetValue("test-text-secret")
	a.createTextSecretModel.metadata.SetValue("test-meta")
	a.createTextSecretModel.textData.SetValue("test text data")

	newModel, cmd := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyEnter})

	newApp, ok := newModel.(app)
	assert.True(t, ok)

	assert.Equal(t, secretsState, newApp.state)
	assert.Empty(t, newApp.createTextSecretModel.title.Value())
	assert.Empty(t, newApp.createTextSecretModel.metadata.Value())
	assert.Empty(t, newApp.createTextSecretModel.textData.Value())
	assert.Nil(t, newApp.createTextSecretModel.err)
	assert.NotNil(t, cmd)
}

func TestApp_UpdateCreateTextSecret_WithError(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
			err:      nil,
		},
	}

	newModel, _ := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyEnter})
	newApp := newModel.(app)
	assert.Equal(t, "text data is required", newApp.createTextSecretModel.err.Error())

	a.createTextSecretModel.title.SetValue("test-title")
	a.createTextSecretModel.textData.SetValue("test text")

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyCtrlB})
	newApp = newModel.(app)
	assert.Equal(t, "text data is required", newApp.createTextSecretModel.err.Error())
	assert.Equal(t, secretsState, newApp.state)
	assert.Empty(t, newApp.createTextSecretModel.title.Value())
	assert.Empty(t, newApp.createTextSecretModel.textData.Value())
}

func TestApp_UpdateCreateTextSecret_Navigation_Down(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()
	newModel, _ := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyDown})
	newApp := newModel.(app)
	assert.True(t, newApp.createTextSecretModel.metadata.Focused())
	assert.False(t, newApp.createTextSecretModel.title.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyDown})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.textData.Focused())
	assert.False(t, newApp.createTextSecretModel.metadata.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyDown})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.title.Focused())
	assert.False(t, newApp.createTextSecretModel.textData.Focused())
}

func TestApp_UpdateCreateTextSecret_Navigation_Up(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()
	newModel, _ := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyUp})
	newApp := newModel.(app)
	assert.True(t, newApp.createTextSecretModel.textData.Focused())
	assert.False(t, newApp.createTextSecretModel.title.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyUp})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.metadata.Focused())
	assert.False(t, newApp.createTextSecretModel.textData.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyUp})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.title.Focused())
	assert.False(t, newApp.createTextSecretModel.metadata.Focused())
}

func TestApp_UpdateCreateTextSecret_Navigation_Tab(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()
	newModel, _ := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyTab})
	newApp := newModel.(app)
	assert.True(t, newApp.createTextSecretModel.metadata.Focused())
	assert.False(t, newApp.createTextSecretModel.title.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyTab})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.textData.Focused())
	assert.False(t, newApp.createTextSecretModel.metadata.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyTab})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.title.Focused())
	assert.False(t, newApp.createTextSecretModel.textData.Focused())
}

func TestApp_UpdateCreateTextSecret_Navigation_ShiftTab(t *testing.T) {
	a := app{
		client: grpcclient.New(":8080"),
		state:  secretTextCreateState,
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.Focus()
	newModel, _ := a.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyShiftTab})
	newApp := newModel.(app)
	assert.True(t, newApp.createTextSecretModel.textData.Focused())
	assert.False(t, newApp.createTextSecretModel.title.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyShiftTab})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.metadata.Focused())
	assert.False(t, newApp.createTextSecretModel.textData.Focused())

	newModel, _ = newApp.updateCreateTextSecret(tea.KeyMsg{Type: tea.KeyShiftTab})
	newApp = newModel.(app)
	assert.True(t, newApp.createTextSecretModel.title.Focused())
	assert.False(t, newApp.createTextSecretModel.metadata.Focused())
}

func TestApp_CreateTextSecretView(t *testing.T) {
	a := app{
		createTextSecretModel: createTextSecretModel{
			title:    textinput.New(),
			metadata: textinput.New(),
			textData: textinput.New(),
		},
	}

	a.createTextSecretModel.title.SetValue("test-text-secret")
	a.createTextSecretModel.metadata.SetValue("test-meta")
	a.createTextSecretModel.textData.SetValue("test text data")

	s := a.createTextSecretView()

	assert.Contains(t, s, "New secret (type text)")
	assert.Contains(t, s, "Title*> test-text-secret")
	assert.Contains(t, s, "Metadata> test-meta")
	assert.Contains(t, s, "Text data*> test text data")
	assert.Contains(t, s, "[Ctrl+B] - back in secrets")
}
