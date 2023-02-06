package main

import (
	"bytes"
	"github.com/andrianbdn/wg-dir-conf/backend"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mdp/qrterminal/v3"
)

type ViewPeer struct {
	PeerName  string
	Comment   string
	Config    string
	sSize     tea.WindowSizeMsg
	qrEnabled bool
	qrMode    bool
}

func NewViewPeer(sSize tea.WindowSizeMsg, srv *backend.Server, cl *backend.Client) ViewPeer {
	cfg, err := cl.GetPlainTextConfig(srv)
	qrEnabled := true
	if err != nil {
		cfg = "Error: " + err.Error()
		qrEnabled = false
	}
	return ViewPeer{sSize: sSize, PeerName: cl.GetName(), Config: cfg, qrEnabled: qrEnabled}
}

func (m ViewPeer) Init() tea.Cmd {
	return nil
}

func (m ViewPeer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.sSize = msg
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape, tea.KeyF10:
			return m, func() tea.Msg {
				return TuiDialogMsgResult{}
			}
		case tea.KeyF9:
			m.qrMode = !m.qrMode
			return m, nil
		}
	}
	return m, cmd
}

func (m ViewPeer) View() string {
	ccyan := lipgloss.Color("6")
	cblk := lipgloss.Color("0")
	cblue := lipgloss.Color("4")
	cgray := lipgloss.Color("7")
	cwhite := lipgloss.Color("15")

	header := lipgloss.NewStyle().Background(ccyan).Foreground(cblk).Width(m.sSize.Width)
	body := lipgloss.NewStyle().Background(cblue).Foreground(cgray).Width(m.sSize.Width).Height(m.sSize.Height - 2)

	whiteOnBlack := lipgloss.NewStyle().Background(cblk).Foreground(cwhite)
	blackOnCyan := lipgloss.NewStyle().Background(ccyan).Foreground(cblk)

	fbtn := func(btn, text string) string {
		b := blackOnCyan.Copy().Width(12)
		return whiteOnBlack.Render(btn) + b.Render(text)
	}

	helpLine := ""
	if m.qrEnabled {
		qrButtonTitle := "QR"
		if m.qrMode {
			qrButtonTitle = "Text"
		}
		helpLine = fbtn("F9", qrButtonTitle) + whiteOnBlack.Render("  ")
	}
	helpLine += fbtn("F10", "Close")

	bw := m.sSize.Width - lipgloss.Width(helpLine)
	helpLine = whiteOnBlack.Copy().Width(bw).Render(" ") + helpLine

	config := m.Config
	if m.qrMode {
		buf := bytes.NewBuffer(nil)
		qrterminal.GenerateHalfBlock(config, qrterminal.L, buf)
		config = buf.String()
	}

	return lipgloss.JoinVertical(0,
		header.Render("Peer \""+m.PeerName+"\""),
		body.Render(config),
		helpLine,
	)
}
