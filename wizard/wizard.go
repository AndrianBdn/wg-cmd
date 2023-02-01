package wizard

import (
	"github.com/andrianbdn/wg-dir-conf/app"
	"github.com/andrianbdn/wg-dir-conf/backend"
	"github.com/andrianbdn/wg-dir-conf/sysinfo"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
)

const optIdNAT = "optIdNAT"
const optIdDNS = "optIdDNS"

type Done struct {
	InterfaceName string
}

type RootModel struct {
	app           *app.App
	stepInterface interfaceScreen
	currentModel  tea.Model
	blueprint     backend.ServerBlueprint
	sSize         tea.WindowSizeMsg
}

func NewRootModel(app *app.App) RootModel {
	m := RootModel{}
	m.app = app
	m.currentModel = newWelcomeDirTestScreen(app, tea.WindowSizeMsg{})
	return m
}

func (m RootModel) Init() tea.Cmd {
	if m.currentModel != nil {
		return m.currentModel.Init()
	}
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
	}

	if _, ok := msg.(welcomeDirTestResult); ok {
		m.currentModel = newInterfaceScreen(m.app, m.sSize)
		return m, m.currentModel.Init()
	}

	if msg, ok := msg.(interfaceScreenResult); ok {
		m.blueprint.InterfaceName = string(msg)
		//step2 := newEndpointStep(m.app, m.sSize)

		portStep := newPortScreen(m.sSize)
		m.currentModel = portStep
		return m, portStep.Init()
	}

	if msg, ok := msg.(portStepResult); ok {
		m.blueprint.Port = uint16(msg)

		endStep := newEndpointStep(m.sSize)
		m.currentModel = endStep
		return m, endStep.Init()
	}

	if msg, ok := msg.(endpointStepResult); ok {
		m.blueprint.Endpoint = string(msg)
		netStep := newNetScreen(m.sSize)
		m.currentModel = netStep
		return m, netStep.Init()
	}

	if msg, ok := msg.(netStepResult); ok {
		m.blueprint.Net4 = msg.net4
		m.blueprint.Net6 = msg.net6
		return m.presentNatDialog()
	}

	if msg, ok := msg.(optScreenResult); ok {
		if msg.id == optIdNAT {

			if msg.result.id == "nat46" {
				m.blueprint.Nat4 = true
				m.blueprint.Nat4 = true
			} else if msg.result.id == "nat4" {
				m.blueprint.Nat4 = true
				m.blueprint.Nat6 = false
			} else {
				m.blueprint.Nat4 = false
				m.blueprint.Nat6 = false
			}

			return m.presentDNSDialog()
		}

		if msg.id == optIdDNS {
			m.blueprint.DNS = msg.result.id
			doneStep := newDoneScreen(m.app, m.sSize, m.blueprint)
			m.currentModel = doneStep
			return m, doneStep.Init()
		}

	}

	if _, ok := msg.(doneScreenResult); ok {
		m.currentModel = nil
		return m, func() tea.Msg {
			return Done{InterfaceName: m.blueprint.InterfaceName}
		}
	}

	if m.currentModel != nil {
		w, c := m.currentModel.Update(msg)
		m.currentModel = w
		return m, c
	}
	return m, nil
}

func (m RootModel) View() string {
	if m.currentModel != nil {
		return m.currentModel.View()
	}
	return ""
}

func (m RootModel) presentNatDialog() (tea.Model, tea.Cmd) {
	opts := []opt{
		{"nat46", "Use NAT to provide the Internet access to the VPN, both via IP4 and IP6"},
		{"nat4", "Use NAT to provide the Internet access to the VPN, only via IP4"},
		{"no", "Skip NAT, only setup basic networking"},
	}

	if sysinfo.HasIP6() == false || m.blueprint.Net6 == "" {
		opts = opts[1:]
	}

	optStep := newOptionScreen(m.app, m.sSize, optIdNAT, opts)
	optStep.setPrompt("Choose if the Setup should configure Network Address Translation (NAT) for the VPN network")

	m.currentModel = optStep
	return m, optStep.Init()
}

func (m RootModel) presentDNSDialog() (tea.Model, tea.Cmd) {
	opts := []opt{
		{"1.1.1.1", "Use Cloudflare DNS https://1.1.1.1"},
		{"8.8.8.8", "Use Google DNS https://developers.google.com/speed/public-dns"},
		{"9.9.9.9", "Use Quad9 DNS https://www.quad9.net/"},
		{"208.67.222.222", "Use OpenDNS https://use.opendns.com/"},
	}

	// shuffle options, so no service will get default treatment (on avarage)
	rand.Shuffle(len(opts), func(i, j int) { opts[i], opts[j] = opts[j], opts[i] })

	opts = append(opts, opt{"", "Do not setup any DNS for clients (may cause leakage)"})

	optStep := newOptionScreen(m.app, m.sSize, optIdDNS, opts)
	optStep.setPrompt("Choose DNS service for client configuration files. Last option is not recommended for providing " +
		"the Internet access via NAT, due to possible leakage. You can later change DNS server to any other address.")

	m.currentModel = optStep
	return m, optStep.Init()
}
