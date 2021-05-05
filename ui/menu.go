package ui

import (
	"github.com/getlantern/systray"
	"github.com/jurjevic/SplitVPN/icon"
	"github.com/jurjevic/SplitVPN/split"
	"strconv"
	"time"
)

type Menu struct {
	inetInfo             *MenuInfo
	vpnInfo              *MenuInfo
	ispInfo              *MenuInfo
	extIpInfo            *MenuInfo
	requestSplit         bool
	requestDiagnose      bool
	requestAutomaticMode split.Bool
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

func newSubMenuInfo(title string, size int) *MenuInfo {
	menuTitle := systray.AddMenuItem(title, "")
	var menuItems []*systray.MenuItem
	for i := 0; i < size; i++ {
		menuInfo := menuTitle.AddSubMenuItem("", "")
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
	avg, rateValue, _ := s.Average()

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

	http := ""
	if s.IsDefault() {
		if s.HttpRequest() {
			http = "HTTP: Ok"
		} else {
			http = "HTTP: Error"
		}
	}

	since := ""
	if s.ActiveSince() != (time.Time{}) {
		since = "Since: " + s.ActiveSince().Format("2006-01-02 15:04:05")
	}

	ping := ""
	if avg > 0 {
		ping = "Ping: " + strconv.Itoa(int(avg.Milliseconds())) + " ms"
	}

	rate := ""
	if rateValue > 0 {
		rate = "Rate: " + textScala(rateValue, 100)
	}

	receiver.Update([]string{gw, ifname, route, host, http, ping, rate, since})
}

func textScala(value, max int) string {
	rate := 10 * value / max
	switch rate {
	case 10:
		return "■■■■■■■■■■"
	case 9:
		return "■■■■■■■■■□"
	case 8:
		return "■■■■■■■■□□"
	case 7:
		return "■■■■■■■□□□"
	case 6:
		return "■■■■■■□□□□"
	case 5:
		return "■■■■■□□□□□"
	case 4:
		return "■■■■□□□□□□"
	case 3:
		return "■■■□□□□□□□"
	case 2:
		return "■■□□□□□□□□"
	case 1:
		return "■□□□□□□□□□"
	default:
		return "□□□□□□□□□□"
	}
}

func Setup() *Menu {

	m := &Menu{}

	infoVpn, infoInet := createNetMenu()

	createBrowserExternalIp(m)
	createSplitNow(m)
	createDebug(m)
	createInfoBox()
	createExitMenu()

	m.inetInfo = infoInet
	m.vpnInfo = infoVpn

	return m
}

func (m *Menu) StateChanged(state split.State, isp split.Isp) {
	switch state {
	case split.NoConnected:
		m.updateConnectedState(split.Isp{})
	case split.VpnConnected:
		m.updateConnectedState(isp)
	case split.InternetConnected:
		m.updateConnectedState(isp)
	case split.Connected:
		m.updateConnectedState(isp)
	}
}

func (m *Menu) updateConnectedState(isp split.Isp) {
	if isp.Query == "" {
		m.extIpInfo.title.Hide()
		m.ispInfo.title.Hide()
	} else {
		m.extIpInfo.title.SetTitle("📡 " + isp.Query)
		m.extIpInfo.title.Show()
		m.extIpInfo.Update([]string{
			"Status: " + isp.Status,
			"Mobile: " + strconv.FormatBool(isp.Mobile),
			"Proxy: " + strconv.FormatBool(isp.Proxy),
			"Hosting: " + strconv.FormatBool(isp.Hosting),
		})
		m.ispInfo.title.SetTitle("🔭 " + isp.Isp)
		m.ispInfo.title.Show()
		m.ispInfo.Update([]string{
			isp.As,
			isp.Asname,
			isp.City,
			isp.RegionName,
			isp.Country,
			isp.Continent,
		})
	}
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
	_, _, ok := vpn.Average()
	if !ok {
		m.vpnInfo.UpdateNotConnected()
	} else {
		m.vpnInfo.UpdateConnected(vpn)
	}

	_, _, ok = inet.Average()
	if !ok {
		m.inetInfo.UpdateNotConnected()
	} else {
		m.inetInfo.UpdateConnected(inet)
	}
	systray.SetIcon(ico)

	r := split.Response{
		SplitNow:      false,
		Diagnose:      false,
		AutomaticMode: split.Nil,
	}

	if m.requestSplit {
		r.SplitNow = true
		m.requestSplit = false
	}
	if m.requestDiagnose {
		r.Diagnose = true
		m.requestDiagnose = false
	}
	if m.requestAutomaticMode != split.Nil {
		r.AutomaticMode = m.requestAutomaticMode
		m.requestAutomaticMode = split.Nil
	}

	return r
}

func getStatus(duration time.Duration, rate int, ok bool) int {
	if !ok {
		return 0
	} else if duration < 20*time.Millisecond && rate > 50 {
		return 3
	} else if duration < 200*time.Millisecond {
		return 2
	} else {
		return 1
	}
}
