package ui

import (
	"github.com/getlantern/systray"
	"log"
	"os/exec"
	"runtime"
)


func createExitMenu() {
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit SplitVPN")
	go func() {
		<-mQuit.ClickedCh
		log.Println("Exiting SplitVPN")
		systray.Quit()
	}()
}

func createNetMenu() (*MenuInfo, *MenuInfo) {
	infoVpn := newMenuInfo("🔐 VPN", 7)
	infoVpn.UpdateNotConnected()
	systray.AddSeparator()
	infoInet := newMenuInfo("🌍 INET", 7)
	infoInet.UpdateNotConnected()
	return infoVpn, infoInet
}

func createSplitNow(m *Menu) {
	systray.AddSeparator()
	mSplit := systray.AddMenuItem("Split now", "Execute a split manually")
	go func() {
		for true {
			<-mSplit.ClickedCh
			m.requestSplit = true
		}
	}()
}

func createBrowserExternalIp(m *Menu) {
	systray.AddSeparator()
	extIpBox := newSubMenuInfo("📡 Unknown IP", 4)
	extIpBox.title.Hide()
	m.extIpInfo = extIpBox
	go func() {
		for true {
			<-extIpBox.title.ClickedCh
			cmd := `open "https://ip-api.com/"`
			_, err := exec.Command("/bin/sh", "-c", cmd).Output()
			if err != nil {
				log.Printf("Failed to open info page. Got error. %s", err.Error())
			}
		}
	}()
	infoBox := newSubMenuInfo("🔭 Unknown ISP", 6)
	infoBox.title.Hide()
	m.ispInfo = infoBox
	go func() {
		for true {
			<-infoBox.title.ClickedCh
			cmd := `open "https://ip-api.com/"`
			_, err := exec.Command("/bin/sh", "-c", cmd).Output()
			if err != nil {
				log.Printf("Failed to open info page. Got error. %s", err.Error())
			}
		}
	}()
}

func createInfoBox() {
	systray.AddSeparator()
	infoBox := newSubMenuInfo("🚀 Information", 4)
	infoBox.Update(
		[]string{"Version: " + Version,
			"License: MIT",
			"Build with: " + runtime.Version(),
			"https://github.com/jurjevic/SplitVPN"})
	go func() {
		for true {
			<-infoBox.title.ClickedCh
			cmd := `open "https://github.com/jurjevic/SplitVPN"`
			_, err := exec.Command("/bin/sh", "-c", cmd).Output()
			if err != nil {
				log.Printf("Failed to open info page. Got error. %s", err.Error())
			}
		}
	}()
}
