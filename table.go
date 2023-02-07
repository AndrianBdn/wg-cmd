package main

import (
	"github.com/76creates/stickers"
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"sort"
	"strconv"
)

func newTable(app *app.App, sSize tea.WindowSizeMsg) *stickers.TableSingleType[string] {

	rows := make([][]string, 0, len(app.State.Clients)+1)
	rows = append(rows, []string{"0001", "Server (" + app.State.Server.Interface + ")",
		app.State.Server.Address4, app.State.Server.Address6})

	keys := make([]int, len(app.State.Clients))

	i := 0
	for k := range app.State.Clients {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	for _, k := range keys {
		cl := app.State.Clients[k]
		if cl == nil {
			continue
		}
		rows = append(rows, []string{cl.GetIPNumberString(), cl.GetName(), cl.GetIP4(app.State.Server), cl.GetIP6(app.State.Server)})
	}

	headers := []string{
		"ID",
		"Peer Name",
		"IPv4",
		"IPv6",
	}

	t := stickers.NewTableSingleType[string](0, 0, headers)
	ratio := []int{1, 10, 5, 7}
	minSize := []int{5, 10, 11, 16}
	t.SetRatio(ratio).SetMinWidth(minSize)
	t.AddRows(rows)
	t.SetHeight(sSize.Height - 1)
	t.SetWidth(sSize.Width)

	return t
}

func peerRow(app *app.App, table *stickers.TableSingleType[string]) *backend.Client {
	cellID := table.GetCursorValue()
	peerID, err := strconv.Atoi(cellID)
	if err != nil {
		log.Println("can't convert", cellID, "to int", err)
		return nil
	}
	return app.State.Clients[peerID]
}
