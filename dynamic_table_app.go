package main

import (
	"fmt"
	"github.com/andrianbdn/wg-cmd/app"
	"github.com/andrianbdn/wg-cmd/backend"
	"log"
	"sort"
	"strconv"
)

func stringRowsFromApp(app *app.App) [][]string {
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

	kkk := 1
	for j := 0; j < 1; j++ {
		for _, k := range keys {
			cl := app.State.Clients[k]
			if cl == nil {
				continue
			}
			// cl.GetIPNumberString()
			kkk++
			fake := fmt.Sprintf("%04d", kkk)
			fake = cl.GetIPNumberString()

			rows = append(rows, []string{fake, cl.GetName(), cl.GetIP4(app.State.Server), cl.GetIP6(app.State.Server)})
		}
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
