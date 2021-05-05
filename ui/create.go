package ui

import (
	"github.com/getlantern/systray"
	"github.com/jurjevic/SplitVPN/split"
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

func createDebug(m *Menu) {
	systray.AddSeparator()
	mDiagnose := systray.AddMenuItemCheckbox("Diagnose", "Print diagnose information to the console", false)
	go func() {
		for true {
			<-mDiagnose.ClickedCh
			m.requestDiagnose = true
		}
	}()
	mDebugMode := systray.AddMenuItemCheckbox("Debug Mode", "Produce verbose debug messages to the console", false)
	go func() {
		for true {
			<-mDebugMode.ClickedCh
			if mDebugMode.Checked() {
				mDebugMode.Uncheck()
				split.DebugFlag = false
			} else {
				mDebugMode.Check()
				split.DebugFlag = true
			}
		}
	}()
	mAutomaticMode := systray.AddMenuItemCheckbox("Automatic Mode", "Detected network changes and perform split automatically", true)
	go func() {
		for true {
			<-mAutomaticMode.ClickedCh
			if mAutomaticMode.Checked() {
				mAutomaticMode.Uncheck()
				m.requestAutomaticMode = split.False
			} else {
				mAutomaticMode.Check()
				m.requestAutomaticMode = split.True
			}
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
		[]string{"Version: " + split.Version,
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
