package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LANReport struct {
	IP string
}

type LANPinger struct {
	MWin        fyne.Window
	IPEntry     *widget.Entry
	Subnet      string
	Reports     []LANReport
	ReportsList *widget.List
	PTNBtn      *widget.Button
	Progress    *widget.ProgressBar
	Select      *widget.Select
}

func NewWindow() *LANPinger {
	fApp := app.New()
	mainWindow := fApp.NewWindow("LANPinger utility")

	ipin := widget.NewEntry()
	ipin.SetPlaceHolder("Enter IP address")
	ipin.Text = "192.168.42.0"
	this := LANPinger{
		MWin:    mainWindow,
		IPEntry: ipin,
		Subnet:  "/32",
		Reports: []LANReport{
			{IP: "NO IPs"},
		},
	}

	combo := widget.NewSelect([]string{"/32", "/24", "/16", "/8"}, func(value string) {
		this.Subnet = value
	})
	combo.SetSelected("/24")
	this.Select = combo

	inpcon := container.New(layout.NewGridLayout(2), ipin, combo)

	startbtn := widget.NewButton("PING THE NET", func() {
		go func() {
			this.runScanner()
		}()
	})
	this.PTNBtn = startbtn

	progress := widget.NewProgressBar()
	this.Progress = progress

	list := widget.NewList(
		func() int {
			return len(this.Reports)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(this.Reports[i].IP)
		})
	this.ReportsList = list

	toplt := container.NewVBox(inpcon, progress, startbtn)
	mainWindow.SetContent(container.New(
		layout.NewBorderLayout(toplt, nil, nil, nil), toplt, list))

	return &this
}

func (this *LANPinger) Run() {
	this.MWin.ShowAndRun()
}

/*

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

 *-/
myApp := app.New()
	myWindow := myApp.NewWindow("LAN Pinger")


	/*
	header := widget.NewVBox(widget.NewLabel("Hello Fyne!"), pb)
	footer := widget.NewVBox(cb, b)
	c := fyne.NewContainerWithLayout(layout.NewBorderLayout(header, footer, nil, nil), header, footer, scrollable)
*-/
list := widget.NewList(
func() int {
	return len(data)
},
func() fyne.CanvasObject {
	return widget.NewLabel("template")
},
func(i widget.ListItemID, o fyne.CanvasObject) {
	o.(*widget.Label).SetText(data[i])
})

go func() {
	for {
		time.Sleep(time.Second * 5)
		data = append(data, time.Now().Format(time.RFC3339))
		list.Refresh()
	}
}()

combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
	log.Println("Select set to", value)
})

myWindow.SetContent(container.NewMax(container.NewVBox(combo, list)))
myWindow.ShowAndRun()
*/
