package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ibeloyar/gophkeeper/internal/client/config"
	"github.com/ibeloyar/gophkeeper/internal/client/grpcclient"

	tea "github.com/charmbracelet/bubbletea"
	gophkeeperv1 "github.com/ibeloyar/gophkeeper/proto/gophkeeper/v1"
)

type state int
type loadSecretsMsg struct{}

const (
	loginState state = iota
	registerState
	secretsState
	secretCardState
	secretLogPassCreateState
	secretTextCreateState
	secretBinaryCreateState
	secretCardCreateState
	secretDeleteState
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
)

type app struct {
	buildVersion string
	buildDate    string
	serverAddr   string

	client                 *grpcclient.GRPCClient
	state                  state
	token                  string
	secrets                []*gophkeeperv1.Secret
	selectedSecret         *gophkeeperv1.Secret
	secretsInitRequestDone bool
	secretsCursor          int

	loginModel    loginForm
	registerModel registerForm
	secretsModel  secretsPage

	createLogPassSecretModel createLogPassSecretModel
	createTextSecretModel    createTextSecretModel
	createBinarySecretModel  createBinarySecretModel
	createCardSecretModel    createCardSecretModel
}

func (a app) Init() tea.Cmd {
	if a.state == secretsState {
		return a.pollSecrets()
	}

	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch a.state {
	case loginState:
		return a.updateLogin(msg)

	case registerState:
		return a.updateRegister(msg)

	case secretsState:
		return a.updateSecrets(msg)

	case secretCardState:
		return a.updateSecretCard(msg)
	case secretLogPassCreateState:
		return a.updateCreateLogPassSecret(msg)
	case secretTextCreateState:
		return a.updateCreateTextSecret(msg)
	case secretBinaryCreateState:
		return a.updateCreateBinarySecret(msg)
	case secretCardCreateState:
		return a.updateCreateCardSecret(msg)
	case secretDeleteState:
		return a.updateSecretDelete(msg)
	}

	return a, nil
}

func (a app) View() string {
	switch a.state {
	case loginState:
		return a.loginView()
	case registerState:
		return a.registerView()
	case secretsState:
		return a.secretsView()
	case secretCardState:
		return a.secretCardView()
	case secretLogPassCreateState:
		return a.createLogPassSecretView()
	case secretTextCreateState:
		return a.createTextSecretView()
	case secretBinaryCreateState:
		return a.createBinarySecretView()
	case secretCardCreateState:
		return a.createCardSecretView()
	case secretDeleteState:
		return a.secretDeleteView()
	}

	return ""
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	grpcClient := grpcclient.New(cfg.Addr)

	p := tea.NewProgram(initialModel(grpcClient, cfg.Addr))
	if _, err := p.Run(); err != nil {
		grpcClient.Close()
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
