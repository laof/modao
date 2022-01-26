package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"github.com/nadoo/glider/cmd/file"
	"github.com/nadoo/glider/cmd/sys"
)

const (
	start   = "Start"
	stop    = "Stop"
	exit    = "Exit"
	version = "v1.0.2"
)

var MinW *walk.MainWindow
var StartButton *walk.PushButton
var StopButton *walk.PushButton
var ExitButton *walk.PushButton
var AssignedTime *walk.Label
var AssignedNodes *walk.LineEdit
var ALog *walk.Label

type versionJson struct {
	Modao string `json:"modao"`
}

func Run() {

	modao, e := walk.NewIconFromImage(ModaoIcon())

	if e != nil {
		fmt.Println("..")
	}

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
		Title:    "modao",
		Size:     Size{380, 240},
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
								AssignedTime.SetText("更新中，请等待...")
								file.UpdateNodes()
								AssignedTime.SetText("更新完成，请重启")
							},
						},
					},
				},
				Menu{
					Text: "help",
					Items: []MenuItem{
						Action{
							Text: "about",
							OnTriggered: func() {
								Popup(MinW, version)
							},
						},
						Action{
							Text: "update",
							OnTriggered: func() {
								res, err := http.Get("https://laof.github.io/assets/version.json")

								if err != nil {
									Popup(MinW, "network bad")
									return
								}
								defer res.Body.Close()
								txt, rerr := ioutil.ReadAll(res.Body)

								if rerr != nil {
									Popup(MinW, "获取版本信息失败")
									return
								}

								data := &versionJson{}
								json.Unmarshal(txt, data)

								if data.Modao == version {
									Popup(MinW, "亲，"+data.Modao+" 已是最新版本，无需更新")
								} else {
									sys.DownloadZip()
								}

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
					Label{AssignTo: &AssignedTime, Text: time},
					LineEdit{AssignTo: &AssignedNodes, Text: nodes, MinSize: Size{120, 80}, ReadOnly: true},
				},
			},
			Composite{
				Layout: Grid{Columns: 3, Alignment: AlignHNearVFar},
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
