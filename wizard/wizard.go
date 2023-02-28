package wizard

import (
	"math/rand"

	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/sysinfo"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	optIdNAT = "optIdNAT"
	optIdDNS = "optIdDNS"
)

type Done struct {
	InterfaceName string
}

type RootModel struct {
	app          *app.App
	currentModel tea.Model
	blueprint    backend.ServerBlueprint
	sSize        tea.WindowSizeMsg
}

func NewRootModel(app *app.App) RootModel {
	ok, noDepsScreen := newNoDepsScreen(app, tea.WindowSizeMsg{})

	m := RootModel{}
	m.app = app
	if !ok {
		m.currentModel = noDepsScreen
		return m
	}
	m.currentModel = newWelcomeScreen(app, tea.WindowSizeMsg{})
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

	if do, m, c := m.WizardStateUpdate(msg); do {
		return m, c
	}

	if m.currentModel != nil {
		w, c := m.currentModel.Update(msg)
		m.currentModel = w
		return m, c
	}
	return m, nil
}

func (m RootModel) WizardStateUpdate(msg tea.Msg) (bool, tea.Model, tea.Cmd) {
	if _, ok := msg.(welcomeScreenResult); ok {
		m.currentModel = newInterfaceScreen(m.app, m.sSize)
		return true, m, m.currentModel.Init()
	}

	if msg, ok := msg.(interfaceScreenResult); ok {
		m.blueprint.InterfaceName = string(msg)
		// step2 := newEndpointStep(m.app, m.sSize)

		portStep := newPortScreen(m.sSize)
		m.currentModel = portStep
		return true, m, portStep.Init()
	}

	if msg, ok := msg.(portStepResult); ok {
		m.blueprint.Port = uint16(msg)

		endStep := newEndpointStep(m.sSize)
		m.currentModel = endStep
		return true, m, endStep.Init()
	}

	if msg, ok := msg.(endpointStepResult); ok {
		m.blueprint.Endpoint = string(msg)
		netStep := newNetScreen(m.sSize)
		m.currentModel = netStep
		return true, m, netStep.Init()
	}

	if msg, ok := msg.(netStepResult); ok {
		m.blueprint.Net4 = msg.net4
		m.blueprint.Net6 = msg.net6
		m, c := m.presentNatDialog()
		return true, m, c
	}

	if msg, ok := msg.(optScreenResult); ok {
		if msg.id == optIdNAT {

			if msg.result.id == "nat46" {
				m.blueprint.Nat4 = true
				m.blueprint.Nat6 = true
			} else if msg.result.id == "nat4" {
				m.blueprint.Nat4 = true
				m.blueprint.Nat6 = false
			} else {
				m.blueprint.Nat4 = false
				m.blueprint.Nat6 = false
			}

			m, c := m.presentDNSDialog()
			return true, m, c
		}

		if msg.id == optIdDNS {
			m.blueprint.DNS = dnsDict[msg.result.id]
			doneStep := newDoneScreen(m.app, m.sSize, m.blueprint)
			m.currentModel = doneStep
			return true, m, doneStep.Init()
		}

	}

	if _, ok := msg.(doneScreenResult); ok {
		m.currentModel = nil

		if sysinfo.IsRoot() && (needNatConfigurationChange(m.blueprint) || sysinfo.HasSystemd()) {
			// extra step to check NAT and setup systemd
			rootLinux := newLinuxMoreScreen(m.app, m.sSize, m.blueprint)
			m.currentModel = rootLinux
			return true, m, rootLinux.Init()
		}

		return true, m, func() tea.Msg {
			return Done{InterfaceName: m.blueprint.InterfaceName}
		}
	}

	if _, ok := msg.(linuxMoreDone); ok {
		m.currentModel = nil
		return true, m, func() tea.Msg {
			return Done{InterfaceName: m.blueprint.InterfaceName}
		}
	}

	return false, nil, nil
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

	if !sysinfo.HasIP6() || m.blueprint.Net6 == "" {
		opts = opts[1:]
	}

	optStep := newOptionScreen(m.app, m.sSize, optIdNAT, opts)
	optStep.setPrompt("Choose if the Setup should configure Network Address Translation (NAT) for the VPN network")

	m.currentModel = optStep
	return m, optStep.Init()
}

func (m RootModel) presentDNSDialog() (tea.Model, tea.Cmd) {
	opts := []opt{
		{"cloudflare", "Use Cloudflare DNS https://1.1.1.1"},
		{"google", "Use Google DNS https://developers.google.com/speed/public-dns"},
		{"quad9", "Use Quad9 DNS https://www.quad9.net/"},
		{"opendns", "Use OpenDNS https://use.opendns.com/"},
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

func needNatConfigurationChange(b backend.ServerBlueprint) bool {
	if b.Nat4 {
		if !sysinfo.NatEnabledIPv4() {
			return true
		}
	}
	if b.Nat6 {
		if !sysinfo.NatEnabledIPv6() {
			return true
		}
	}
	return false
}
