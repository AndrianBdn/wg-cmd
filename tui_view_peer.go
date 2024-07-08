package main

import (
	"bytes"
	"io"
	"strings"

	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"github.com/andrianbdn/wg-cmd/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mdp/qrterminal/v3"
)

type qr struct {
	qrCode string
	size   int
}

type ViewPeer struct {
	title     string
	config    string
	sSize     tea.WindowSizeMsg
	qrEnabled bool
	qrMode    bool
	qrCodes   []qr
	app       *app.App
}

func NewViewPeer(sSize tea.WindowSizeMsg, app *app.App, cl *backend.Client) ViewPeer {
	qrEnabled := true
	qrMode := app.Settings.ViewerQRMode
	title := ""
	var cfg string
	var err error

	if cl != nil {
		cfg, err = cl.GetPlainTextConfig(app.State.Server)
		if err != nil {
			cfg = "Error: " + err.Error()
			qrEnabled = false
			qrMode = false
		}
		title = "Peer \"" + cl.GetName() + "\" • IP4 " + cl.GetIP4(app.State.Server)
		ip6 := cl.GetIP6(app.State.Server)
		if ip6 != "" {
			title += " • IP6 " + ip6
		}
	} else {
		cfg = app.State.Server.GetInterfaceString()
		qrMode = false
		qrEnabled = false
		title = "Server interface " + app.State.Server.Interface
	}

	qrs := make([]qr, 4)

	if qrEnabled {
		qrs[0] = qrGenerate(func(w io.Writer) { qrterminal.Generate(cfg, qrterminal.M, w) })
		qrs[1] = qrGenerate(func(w io.Writer) { qrterminal.Generate(cfg, qrterminal.L, w) })
		qrs[2] = qrGenerate(func(w io.Writer) { qrterminal.GenerateHalfBlock(cfg, qrterminal.L, w) })
		qrs[3] = qr{qrCode: "Terminal size is too small for QR code; enlarge windows / set smaller font", size: 1}
	}

	return ViewPeer{
		sSize:     sSize,
		title:     title,
		config:    cfg,
		qrEnabled: qrEnabled,
		qrMode:    qrMode,
		qrCodes:   qrs,
		app:       app,
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

			m.app.Settings.ViewerQRMode = m.qrMode
			_ = m.app.SaveSettings()
			return m, nil
		}
	}
	return m, cmd
}

func (m ViewPeer) View() string {
	header := theme.Current.ViewerTopBar.Width(m.sSize.Width)
	body := theme.Current.ViewerMain.Width(m.sSize.Width).Height(m.sSize.Height - 2)

	f9 := helpKey{key: "F9", help: "QR", hidden: true}
	if m.qrEnabled {
		if m.qrMode {
			f9.help = "Text"
		}
		f9.hidden = false
	}
	f10 := helpKey{key: "F10", help: "Close"}

	config := m.config
	if m.qrMode {
		for _, q := range m.qrCodes {
			if q.size < m.sSize.Height-1 {
				config = q.qrCode
				break
			}
		}
	}

	return lipgloss.JoinVertical(0,
		header.Render(m.title),
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
