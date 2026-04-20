package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"

	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

func initialModel(client *grpcclient.GRPCClient) app {
	a := app{
		client:                 client,
		token:                  "",
		secretsCursor:          0,
		secrets:                make([]*gophkeeperv1.Secret, 0),
		secretsInitRequestDone: false,
		loginModel: loginForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
		registerModel: registerForm{
			loginInput:    textinput.New(),
			passwordInput: textinput.New(),
		},
		createLogPassSecretModel: createLogPassSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			secretType: gophkeeperv1.SecretType_LOGIN_PASSWORD,
			login:      textinput.New(),
			password:   textinput.New(),
		},
		createTextSecretModel: createTextSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			secretType: gophkeeperv1.SecretType_TEXT,
			textData:   textinput.New(),
		},
		createBinarySecretModel: createBinarySecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			secretType: gophkeeperv1.SecretType_BINARY,
			filePath:   textinput.New(),
		},
		createCardSecretModel: createCardSecretModel{
			title:      textinput.New(),
			metadata:   textinput.New(),
			secretType: gophkeeperv1.SecretType_CARD,
			cardExp:    textinput.New(),
			cardHolder: textinput.New(),
			cardNumber: textinput.New(),
		},
	}

	// loginForm
	a.loginModel.loginInput.Placeholder = "Enter login..."
	a.loginModel.passwordInput.Placeholder = "Enter password..."
	a.loginModel.passwordInput.EchoMode = textinput.EchoPassword

	// registerForm
	a.registerModel.loginInput.Placeholder = "Enter login..."
	a.registerModel.passwordInput.Placeholder = "Enter password..."
	a.registerModel.passwordInput.EchoMode = textinput.EchoPassword

	// createSecretForm
	a.createLogPassSecretModel.title.Placeholder = "Enter title..."
	a.createLogPassSecretModel.metadata.Placeholder = "Enter metadata..."
	a.createLogPassSecretModel.login.Placeholder = "Enter login..."
	a.createLogPassSecretModel.password.Placeholder = "Enter password..."

	// createTextSecretForm
	a.createTextSecretModel.title.Placeholder = "Enter title..."
	a.createTextSecretModel.metadata.Placeholder = "Enter metadata..."
	a.createTextSecretModel.textData.Placeholder = "Enter text..."

	// createBinarySecretModel
	a.createBinarySecretModel.title.Placeholder = "Enter title..."
	a.createBinarySecretModel.metadata.Placeholder = "Enter metadata..."
	a.createBinarySecretModel.filePath.Placeholder = "Enter filename..."

	// createCardSecretModel
	a.createCardSecretModel.title.Placeholder = "Enter title..."
	a.createCardSecretModel.metadata.Placeholder = "Enter metadata..."
	a.createCardSecretModel.cardExp.Placeholder = "Enter card exp (Format: MM/YY)..."
	a.createCardSecretModel.cardNumber.Placeholder = "Enter card number (16 digits)..."
	a.createCardSecretModel.cardHolder.Placeholder = "Enter cardholder..."

	return a
}
