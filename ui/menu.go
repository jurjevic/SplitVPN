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
	inetInfo *MenuInfo
	vpnInfo  *MenuInfo
	requestSplit bool
}

type MenuInfo struct {
	title *systray.MenuItem
	info1 *systray.MenuItem
	info2 *systray.MenuItem
	info3 *systray.MenuItem
	info4 *systray.MenuItem
	info5 *systray.MenuItem
	infoSub *systray.MenuItem
}

const (
	notConnected = "Not connected"
)

func newMenuInfo(title string) *MenuInfo {
	menuTitle := systray.AddMenuItem(title, "")
	menuTitle.Disable()
	menuInfo1 := systray.AddMenuItem("", "")
	menuInfo1.Disable()
	menuInfo1.Hide()
	menuInfo2 := systray.AddMenuItem("", "")
	menuInfo2.Disable()
	menuInfo2.Hide()
	menuInfo3 := systray.AddMenuItem("", "")
	menuInfo3.Disable()
	menuInfo3.Hide()
	menuInfo4 := systray.AddMenuItem("", "")
	menuInfo4.Disable()
	menuInfo4.Hide()
	menuInfo5 := systray.AddMenuItem("", "")
	menuInfo5.Disable()
	menuInfo5.Hide()
	return &MenuInfo{
		title: menuTitle,
		info1: menuInfo1,
		info2: menuInfo2,
		info3: menuInfo3,
		info4: menuInfo4,
		info5: menuInfo5,
	}
}

func (receiver *MenuInfo) Update(info1, info2, info3, info4, info5 string) {
	if info1 == "" {
		receiver.info1.Hide()
		receiver.info1.SetTitle("")
	} else {
		receiver.info1.Show()
		receiver.info1.SetTitle(info1)
	}
	if info2 == "" {
		receiver.info2.Hide()
		receiver.info2.SetTitle("")
	} else {
		receiver.info2.Show()
		receiver.info2.SetTitle(info2)
	}
	if info3 == "" {
		receiver.info3.Hide()
		receiver.info3.SetTitle("")
	} else {
		receiver.info3.Show()
		receiver.info3.SetTitle(info3)
	}
	if info4 == "" {
		receiver.info4.Hide()
		receiver.info4.SetTitle("")
	} else {
		receiver.info4.Show()
		receiver.info4.SetTitle(info4)
	}
	if info5 == "" {
		receiver.info5.Hide()
		receiver.info5.SetTitle("")
	} else {
		receiver.info5.Show()
		receiver.info5.SetTitle(info5)
	}
}

func (receiver *MenuInfo) UpdateNotConnected() {
	receiver.Update(notConnected, "", "", "", "")
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

	ping := ""
	if avg > 0 {
		ping = "Ping: " + strconv.Itoa(int(avg.Milliseconds())) + " ms"
	}

	receiver.Update(gw, ifname, route, host, ping)
}

func Setup() *Menu {

	m := &Menu{}

	infoVpn := newMenuInfo("🔐 VPN")
	infoVpn.UpdateNotConnected()

	systray.AddSeparator()

	infoInet := newMenuInfo("🌍 INET")
	infoInet.UpdateNotConnected()

	systray.AddSeparator()

	infoBox := newMenuInfo("🗂 Information")
	infoBox.Update(
		"Version: " + Version,
		"License: MIT",
		"Build with: " + runtime.Version(),
		"https://github.com/jurjevic/SplitVPN",
		"")

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

	r:=split.Response{}

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
