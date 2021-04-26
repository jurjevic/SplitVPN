package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"github.com/jurjevic/SplitVPN/icon"
	"github.com/jurjevic/SplitVPN/split"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if getProcessOwner() == "root" {
		fmt.Println("You're sudo!")
	} else {
		fmt.Println("root privileges required! Please start with 'sudo'")
		os.Exit(1)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data_0_0)
	systray.SetTooltip("Automatic VPN and Internet Split")
	mQuit := systray.AddMenuItem("Quit", "Quit SplitVPN")
	go func() {
		<-mQuit.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	splitInstance := split.NewSplit()
	splitInstance.Start(func(state split.State, inet *split.Zone, vpn *split.Zone) {
		vpnStatus := getStatus(vpn.Average())
		inetStatus := getStatus(inet.Average())
		ico := icon.Data_0_0
		if vpnStatus == 3 && inetStatus == 3 {
			ico = icon.Data_3_3
		} else if vpnStatus == 3 && inetStatus == 2 {
			ico = icon.Data_3_2
		} else if vpnStatus == 3 && inetStatus == 1 {
			ico = icon.Data_3_1
		} else if vpnStatus == 3 && inetStatus == 0 {
			ico = icon.Data_3_0
		} else if vpnStatus == 2 && inetStatus == 3 {
			ico = icon.Data_2_3
		} else if vpnStatus == 2 && inetStatus == 2 {
			ico = icon.Data_2_2
		} else if vpnStatus == 2 && inetStatus == 1 {
			ico = icon.Data_2_1
		} else if vpnStatus == 2 && inetStatus == 0 {
			ico = icon.Data_2_0
		} else if vpnStatus == 1 && inetStatus == 3 {
			ico = icon.Data_1_3
		} else if vpnStatus == 1 && inetStatus == 2 {
			ico = icon.Data_1_2
		} else if vpnStatus == 1 && inetStatus == 1 {
			ico = icon.Data_1_1
		} else if vpnStatus == 1 && inetStatus == 0 {
			ico = icon.Data_1_0
		} else if vpnStatus == 0 && inetStatus == 3 {
			ico = icon.Data_0_3
		} else if vpnStatus == 0 && inetStatus == 2 {
			ico = icon.Data_0_2
		} else if vpnStatus == 0 && inetStatus == 1 {
			ico = icon.Data_0_1
		} else if vpnStatus == 0 && inetStatus == 0 {
			ico = icon.Data_0_0
		}
		systray.SetIcon(ico)
	})
}

func getStatus(duration time.Duration, ok bool) int {
	if !ok {
		return 0
	} else if duration < 20*time.Millisecond {
		return 3
	} else if duration < 200*time.Millisecond {
		return 2
	} else {
		return 1
	}
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
