package main

import (
	"log"
	"strconv"

	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
)

func stringRowsFromApp(app *app.App) [][]string {
	rows := make([][]string, 0, len(app.State.Clients)+1)
	rows = append(rows, []string{
		"0001", "Server (" + app.State.Server.Interface + ")",
		app.State.Server.Address4, app.State.Server.Address6,
	})

	err := app.State.IterateClients(func(cl *backend.Client) error {
		rows = append(rows, []string{
			cl.GetIPNumberString(),
			cl.GetName(),
			cl.GetIP4(app.State.Server),
			cl.GetIP6(app.State.Server),
		})
		return nil
	})
	if err != nil {
		panic(err)
	}

	return rows
}

func newAppDynamicTableList(app *app.App, table *DynamicTableList) DynamicTableList {
	t := NewMainTable(
		[]string{"#", "Name", "IPv4", "IPv6"},
		stringRowsFromApp(app),
		[]int{1, 3, 2, 3},
		[]int{4, 20, 16, 46},
	)
	if table != nil {
		t.CopyTableState(table)
	}
	return t
}

func clientFromRow(app *app.App, row []string) *backend.Client {
	peerID, err := strconv.Atoi(row[0])
	if err != nil {
		log.Println("can't convert", row[0], "to int", err)
		return nil
	}
	return app.State.Clients[peerID]
}
