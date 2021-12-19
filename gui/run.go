package gui

import (
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/nadoo/glider/cmd/file"
	"github.com/nadoo/glider/cmd/sys"
)

const (
	start = "Start"
	stop  = "Stop"
	exit  = "Exit"
)

var MinW *walk.MainWindow
var StartButton *walk.PushButton
var StopButton *walk.PushButton
var ExitButton *walk.PushButton
var AssignedNodes *walk.Label
var AssignedTime *walk.Label
var ALog *walk.Label

func Run() {

	modao, e := walk.NewIconFromImage(ModaoIcon())

	var time, nodes string
	if file.ConfiInfo != "" {
		s := strings.Split(file.ConfiInfo, ",")
		time = s[0]
		nodes = s[1]
	} else {
		nodes = "配置信息读取失败，请更新节点"
	}

	if e != nil {
		return
	}
	mw := MainWindow{
		Title:    "Modao",
		Size:     Size{400, 260},
		Font:     Font{PointSize: 11},
		Layout:   VBox{},
		Icon:     modao,
		AssignTo: &MinW,

		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonTextOnly,
			Items: []MenuItem{
				Menu{
					Text: "node",
					Items: []MenuItem{

						Action{
							Text: "update",
							OnTriggered: func() {
								AssignedNodes.SetText("")
								AssignedTime.SetText("节点更新中，请等待...")
								file.UpdateNodes()
								AssignedTime.SetText("更新完成，请重启")
							},
						},
					},
				},
				Menu{
					Text: "version",
					Items: []MenuItem{
						Action{
							Text: "about",
							OnTriggered: func() {
								Popup2(MinW, "v1.0.1")
							},
						},
					},
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: VBox{},
				Children: []Widget{
					Label{AssignTo: &AssignedNodes, Text: nodes},
					Label{AssignTo: &AssignedTime, Text: time},
				},
			},
			Composite{
				Layout: Grid{Columns: 3, Margins: Margins{Top: 5}, Alignment: AlignHNearVFar, MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					PushButton{
						Text:     start,
						AssignTo: &StartButton,
						OnClicked: func() {
							StartButton.SetEnabled(false)
							StopButton.SetEnabled(true)
							sys.SetProxy(1)
						},
					},
					PushButton{
						Text:     stop,
						Enabled:  false,
						AssignTo: &StopButton,
						OnClicked: func() {
							StartButton.SetEnabled(true)
							StopButton.SetEnabled(false)
							sys.SetProxy(0)
						},
					},
					PushButton{
						Text:     exit,
						AssignTo: &ExitButton,
						OnClicked: func() {
							sys.SetProxy(-1)
						},
					},
				},
			},
		},
		OnSizeChanged: func() {
		},
	}

	mw.Create()

	ww, wh := int(win.GetSystemMetrics(0)), int(win.GetSystemMetrics(1))

	width, height := MinW.HeightPixels(), MinW.WidthPixels()

	MinW.SetXPixels((ww - width) / 2)
	MinW.SetYPixels((wh - height) / 2)

	MinW.Run()

}

func update() {

}
