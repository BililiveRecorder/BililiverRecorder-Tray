package systemTray

import (
	_ "embed"
	"os"

	"github.com/skycoin/systray"
)

//go:embed icon.ico
var icon []byte
var eventCh chan os.Signal

func Setup(ch chan os.Signal) {
	eventCh = ch
	systray.Run(onReady, onExit)
}
func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("BililiveRecorder")
	systray.SetTooltip("BililiveRecorder")
	mQuit := systray.AddMenuItem("Quit", "Quit BililiveRecorder")
	mQuit.SetIcon(icon)
	mQuit.Enable()
	systray.AddSeparator()
	systray.AddMenuItem("About", "About BililiveRecorder")
	go func() {
		<-mQuit.ClickedCh
		eventCh <- os.Interrupt
	}()
}
func onExit() {}

func Quit() {
	systray.Quit()
}
