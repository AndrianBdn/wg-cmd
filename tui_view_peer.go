package main

import (
	"bytes"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mdp/qrterminal/v3"
	"io"
	"strings"
)

type qr struct {
	qrCode string
	size   int
}

type ViewPeer struct {
	PeerName  string
	PeerInfo  string
	Comment   string
	Config    string
	sSize     tea.WindowSizeMsg
	qrEnabled bool
	qrMode    bool
	qrCodes   []qr
}

var qrPersist bool

func NewViewPeer(sSize tea.WindowSizeMsg, srv *backend.Server, cl *backend.Client) ViewPeer {
	cfg, err := cl.GetPlainTextConfig(srv)
	qrEnabled := true
	qrMode := qrPersist
	if err != nil {
		cfg = "Error: " + err.Error()
		qrEnabled = false
		qrMode = false
	}
	qrs := make([]qr, 4)
	qrs[0] = qrGenerate(func(w io.Writer) { qrterminal.Generate(cfg, qrterminal.M, w) })
	qrs[1] = qrGenerate(func(w io.Writer) { qrterminal.Generate(cfg, qrterminal.L, w) })
	qrs[2] = qrGenerate(func(w io.Writer) { qrterminal.GenerateHalfBlock(cfg, qrterminal.L, w) })
	qrs[3] = qr{qrCode: "Terminal size is too small for QR code; enlarge windows / set smaller font", size: 1}

	ip6 := cl.GetIP6(srv)
	if ip6 != "" {
		ip6 = " IP6 " + ip6
	}

	return ViewPeer{
		sSize:     sSize,
		PeerName:  cl.GetName(),
		PeerInfo:  "IP4 " + cl.GetIP4(srv) + ip6,
		Config:    cfg,
		qrEnabled: qrEnabled,
		qrMode:    qrMode,
		qrCodes:   qrs,
	}
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
			qrPersist = m.qrMode
			return m, nil
		}
	}
	return m, cmd
}

func (m ViewPeer) View() string {
	header := theme.Current.ViewerTopBar.Copy().Width(m.sSize.Width)
	body := theme.Current.ViewerMain.Copy().Width(m.sSize.Width).Height(m.sSize.Height - 2)

	f9 := helpKey{key: "F9", help: "QR", hidden: true}
	if m.qrEnabled {
		if m.qrMode {
			f9.help = "Text"
		}
		f9.hidden = false
	}
	f10 := helpKey{key: "F10", help: "Close"}

	config := m.Config
	if m.qrMode {
		for _, q := range m.qrCodes {
			if q.size < m.sSize.Height-1 {
				config = q.qrCode
				break
			}
		}
	}

	return lipgloss.JoinVertical(0,
		header.Render("Peer \""+m.PeerName+"\" • "+m.PeerInfo),
		body.Render(config),
		RenderHelpLine(m.sSize.Width, f9, f10),
	)
}

func qrGenerate(cb func(io.Writer)) qr {
	buf := bytes.NewBuffer(nil)
	cb(buf)
	q := strings.TrimSpace(buf.String())
	return qr{qrCode: q, size: lipgloss.Height(q)}
}
