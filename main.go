package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/jurjevic/SplitVPN/icon"
	"github.com/jurjevic/SplitVPN/split"
	"github.com/jurjevic/SplitVPN/ui"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)


func main() {

	ui.Version = Version
	if getProcessOwner() != "root" {
		log.Fatal("'root' privileges required! Please start with 'sudo'")
		os.Exit(1)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data_0_0)
	systray.SetTooltip("Automatic VPN and Internet Split")

	menu := ui.Setup()

	splitInstance := split.NewSplit()
	splitInstance.Start(menu.Refresh)
}

func onExit() {
	// clean up here

}

func getProcessOwner() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.ReplaceAll(string(stdout), "\n", "")
}
