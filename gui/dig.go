package gui

import "github.com/lxn/walk"

func Popup(mw *walk.MainWindow, str string) {
	// 按键结果以int形式返回
	// cmd :=walk.MsgBox(
	walk.MsgBox(
		mw,
		"提示",
		str,
		walk.MsgBoxOK,
		// walk.MsgBoxCancelTryContinue,
		// walk.MsgBoxYesNoCancel,
	)

}
