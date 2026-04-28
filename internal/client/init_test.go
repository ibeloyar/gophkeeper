package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"
	"github.com/stretchr/testify/assert"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func TestInitialModel_BasicFields(t *testing.T) {
	buildVersion = "v1.0.0"
	buildDate = "2026-03-04"

	client := grpcclient.New(":8080")
	addr := "localhost:8080"

	a := initialModel(client, addr)

	assert.Equal(t, "v1.0.0", a.buildVersion)
	assert.Equal(t, "2026-03-04", a.buildDate)
	assert.Equal(t, "localhost:8080", a.serverAddr)

	assert.Equal(t, client, a.client)
	assert.Equal(t, "", a.token)
	assert.Equal(t, 0, a.secretsCursor)
	assert.NotNil(t, a.secrets)
	assert.Len(t, a.secrets, 0)
	assert.False(t, a.secretsInitRequestDone)
}

func TestInitialModel_LoginForm(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.loginModel.loginInput)
	assert.Equal(t, "Enter login...", a.loginModel.loginInput.Placeholder)

	assert.NotNil(t, a.loginModel.passwordInput)
	assert.Equal(t, "Enter password...", a.loginModel.passwordInput.Placeholder)
	assert.Equal(t, textinput.EchoPassword, a.loginModel.passwordInput.EchoMode)
}

func TestInitialModel_RegisterForm(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.registerModel.loginInput)
	assert.Equal(t, "Enter login...", a.registerModel.loginInput.Placeholder)

	assert.NotNil(t, a.registerModel.passwordInput)
	assert.Equal(t, "Enter password...", a.registerModel.passwordInput.Placeholder)
	assert.Equal(t, textinput.EchoPassword, a.registerModel.passwordInput.EchoMode)
}

func TestInitialModel_CreateLogPassSecretModel(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.createLogPassSecretModel.title)
	assert.Equal(t, "Enter title...", a.createLogPassSecretModel.title.Placeholder)

	assert.NotNil(t, a.createLogPassSecretModel.metadata)
	assert.Equal(t, "Enter metadata...", a.createLogPassSecretModel.metadata.Placeholder)

	assert.NotNil(t, a.createLogPassSecretModel.login)
	assert.Equal(t, "Enter login...", a.createLogPassSecretModel.login.Placeholder)

	assert.NotNil(t, a.createLogPassSecretModel.password)
	assert.Equal(t, "Enter password...", a.createLogPassSecretModel.password.Placeholder)

	assert.Equal(t, gophkeeperv1.SecretType_LOGIN_PASSWORD, a.createLogPassSecretModel.secretType)
}

func TestInitialModel_CreateTextSecretModel(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.createTextSecretModel.title)
	assert.Equal(t, "Enter title...", a.createTextSecretModel.title.Placeholder)

	assert.NotNil(t, a.createTextSecretModel.metadata)
	assert.Equal(t, "Enter metadata...", a.createTextSecretModel.metadata.Placeholder)

	assert.NotNil(t, a.createTextSecretModel.textData)
	assert.Equal(t, "Enter text...", a.createTextSecretModel.textData.Placeholder)

	assert.Equal(t, gophkeeperv1.SecretType_TEXT, a.createTextSecretModel.secretType)
}

func TestInitialModel_CreateBinarySecretModel(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.createBinarySecretModel.title)
	assert.Equal(t, "Enter title...", a.createBinarySecretModel.title.Placeholder)

	assert.NotNil(t, a.createBinarySecretModel.metadata)
	assert.Equal(t, "Enter metadata...", a.createBinarySecretModel.metadata.Placeholder)

	assert.NotNil(t, a.createBinarySecretModel.filePath)
	assert.Equal(t, "Enter filename...", a.createBinarySecretModel.filePath.Placeholder)

	assert.Equal(t, gophkeeperv1.SecretType_BINARY, a.createBinarySecretModel.secretType)
}

func TestInitialModel_CreateCardSecretModel(t *testing.T) {
	a := initialModel(grpcclient.New(":8080"), "localhost:8080")

	assert.NotNil(t, a.createCardSecretModel.title)
	assert.Equal(t, "Enter title...", a.createCardSecretModel.title.Placeholder)

	assert.NotNil(t, a.createCardSecretModel.metadata)
	assert.Equal(t, "Enter metadata...", a.createCardSecretModel.metadata.Placeholder)

	assert.NotNil(t, a.createCardSecretModel.cardExp)
	assert.Equal(t, "Enter card exp (Format: MM/YY)...", a.createCardSecretModel.cardExp.Placeholder)

	assert.NotNil(t, a.createCardSecretModel.cardNumber)
	assert.Equal(t, "Enter card number (16 digits)...", a.createCardSecretModel.cardNumber.Placeholder)

	assert.NotNil(t, a.createCardSecretModel.cardHolder)
	assert.Equal(t, "Enter cardholder...", a.createCardSecretModel.cardHolder.Placeholder)

	assert.Equal(t, gophkeeperv1.SecretType_CARD, a.createCardSecretModel.secretType)
}
