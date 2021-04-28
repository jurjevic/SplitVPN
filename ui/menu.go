package ui

import (
	"github.com/getlantern/systray"
	"github.com/jurjevic/SplitVPN/icon"
	"github.com/jurjevic/SplitVPN/split"
	"log"
	"runtime"
	"strconv"
	"time"
)

var (
	Version string
)

type Menu struct {
	inetInfo     *MenuInfo
	vpnInfo      *MenuInfo
	requestSplit bool
}

type MenuInfo struct {
	title *systray.MenuItem
	info  []*systray.MenuItem
}

const (
	notConnected = "Not connected"
)

func newMenuInfo(title string, size int) *MenuInfo {
	menuTitle := systray.AddMenuItem(title, "")
	menuTitle.Disable()
	var menuItems []*systray.MenuItem
	for i := 0; i < size; i++ {
		menuInfo := systray.AddMenuItem("", "")
		menuInfo.Disable()
		menuInfo.Hide()
		menuItems = append(menuItems, menuInfo)
	}

	return &MenuInfo{
		title: menuTitle,
		info:  menuItems,
	}
}

func (receiver *MenuInfo) Update(info []string) {
	for i, menuItem := range receiver.info {
		text := ""
		if i < len(info) {
			text = info[i]
		}
		if text == "" {
			menuItem.Hide()
			menuItem.SetTitle("")
		} else {
			menuItem.Show()
			menuItem.SetTitle(text)
		}
	}
}

func (receiver *MenuInfo) UpdateNotConnected() {
	receiver.Update([]string{notConnected})
}

func (receiver *MenuInfo) UpdateConnected(s *split.Zone) {
	avg, _ := s.Average()

	gw := ""
	if len(s.Gateway()) > 0 {
		gw = "Gateway: " + s.Gateway()
	}

	ifname := ""
	if len(s.InterfaceName()) > 0 {
		ifname = "Interface: " + s.InterfaceName()
	}

	host := ""
	if len(s.Host()) > 0 {
		host = "Host: " + s.Host()
	}

	route := ""
	if len(s.Route()) > 0 {
		route = "Route: " + s.Route()
	}

	since := ""
	if len(s.Route()) > 0 {
		since = "Since: " + s.ActiveSince().Format("2006-01-02 15:04:05")
	}

	ping := ""
	if avg > 0 {
		ping = "Ping: " + strconv.Itoa(int(avg.Milliseconds())) + " ms"
	}

	receiver.Update([]string{gw, ifname, route, host, ping, since})
}

func Setup() *Menu {

	m := &Menu{}

	infoVpn := newMenuInfo("🔐 VPN", 6)
	infoVpn.UpdateNotConnected()

	systray.AddSeparator()

	infoInet := newMenuInfo("🌍 INET", 6)
	infoInet.UpdateNotConnected()

	systray.AddSeparator()

	infoBox := newMenuInfo("🚀 Information", 4)
	infoBox.Update(
		[]string{"Version: " + Version,
			"License: MIT",
			"Build with: " + runtime.Version(),
			"https://github.com/jurjevic/SplitVPN"})

	systray.AddSeparator()
	mSplit := systray.AddMenuItem("Split now", "Execute a split manually")
	go func() {
		for true {
			<-mSplit.ClickedCh
			m.requestSplit = true
		}
	}()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit SplitVPN")

	go func() {
		<-mQuit.ClickedCh
		log.Println("Exiting SplitVPN")
		systray.Quit()
	}()

	m.inetInfo = infoInet
	m.vpnInfo = infoVpn

	return m
}

func (m *Menu) Refresh(state split.State, inet *split.Zone, vpn *split.Zone) split.Response {
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
	_, ok := vpn.Average()
	if !ok {
		m.vpnInfo.UpdateNotConnected()
	} else {
		m.vpnInfo.UpdateConnected(vpn)
	}

	_, ok = inet.Average()
	if !ok {
		m.inetInfo.UpdateNotConnected()
	} else {
		m.inetInfo.UpdateConnected(inet)
	}
	systray.SetIcon(ico)

	r := split.Response{}

	if m.requestSplit {
		r.SplitNow = true
		m.requestSplit = false
	}

	return r
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
