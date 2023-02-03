package main

import (
	"github.com/76creates/stickers"
	"github.com/andrianbdn/wg-dir-conf/app"
)

func newTable(app *app.App) *stickers.TableSingleType[string] {

	rows := make([][]string, 0, len(app.State.Clients)+1)
	rows = append(rows, []string{"0001", "Server (" + app.State.Server.Interface + ")",
		app.State.Server.Address4, app.State.Server.Address6})

	for _, cl := range app.State.Clients {
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

	return t
}
