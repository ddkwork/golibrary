package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/ddkwork/golibrary/src/fynelib/fyneTheme"
	"github.com/ddkwork/golibrary/src/fynelib/myTable"
	"net/http"
	"time"
)

func main() {
	a := app.NewWithID("com.rows.app")
	a.SetIcon(nil)
	fyneTheme.Dark()
	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(1040, 780))
	w.SetMaster()
	w.CenterOnScreen()

	p := new(packets)

	selectionHandler := myTable.WithSelectionHandler(func(id int, selected bool) {
		//println(id)
		//println(selected)
	})
	doubleClickHandler := myTable.WithDoubleClickHandler(func(id int) {
		popUpMenu := widget.NewPopUpMenu(
			fyne.NewMenu("pop",
				fyne.NewMenuItem("copy", func() {
					//println("copy")
				}),
				fyne.NewMenuItem("cut", func() {
					//println("cut")
				}),
				fyne.NewMenuItem("no", func() {
					//println("no")
				}),
			),
			nil,
		)
		popUpMenu.Show()
	})
	list, err := myTable.NewTable(p, selectionHandler, doubleClickHandler)
	if err != nil {
		panic(err.Error())
	}
	go func() {
		//myTable.ColumnWidth = map[int]float32{
		//	2: 300,
		//} //todo
		ticker := time.NewTicker(1 * time.Second)
		defer func() { ticker.Stop() }()
		for i := 0; i < 10; i++ {
			now := time.Now()
			p.Append(&PacketInfo{
				PacketIndex: i + 1,
				Method:      http.MethodGet,
				Host:        "www.baidu.com",
				Path:        "/login",
				ConnectType: "json",
				Size:        246,
				PadTime:     now.Sub(time.Now()),
				StartTime:   now,
				Status:      http.StatusText(http.StatusOK),
				StatusCode:  http.StatusOK,
				Note:        "good",
				Req:         DecodedInfo{},
				Resp:        DecodedInfo{},
			})
			p.Append(&PacketInfo{
				PacketIndex: i + 2,
				Method:      http.MethodPost,
				Host:        "www.baidu.com",
				Path:        "/login",
				ConnectType: "json",
				Size:        246,
				PadTime:     now.Sub(time.Now()),
				StartTime:   now,
				Status:      http.StatusText(http.StatusOK),
				StatusCode:  http.StatusOK,
				Note:        "good",
				Req:         DecodedInfo{},
				Resp:        DecodedInfo{},
			})
			list.SetData(p)
			p.Filter(http.MethodPost, 1)
			<-ticker.C
			list.Refresh()
		}
	}()
	w.SetContent(list)
	w.ShowAndRun()
}
