package split

import (
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type State int

const vpnMask = "/8"

const (
	NoConnected State = iota
	InternetConnected
	VpnConnected
	Connected
)

type split struct {
	inet      *Zone
	vpn       *Zone
	State     State
	router    *router
	automatic bool
}

type routeConfig struct {
	isDefault     bool
	gateway       string
	route         string
	interfaceName string
}

type Bool int

func (b Bool) bool() bool {
	switch b {
	case True:
		return true
	default:
		return false
	}
}

const (
	Nil Bool = iota - 1
	False
	True
)

type Response struct {
	SplitNow      bool
	Diagnose      bool
	AutomaticMode Bool
}

func NewSplit() split {
	return split{
		router:    newRouter(),
		State:     NoConnected,
		inet:      &Zone{host: "google.com"},
		vpn:       &Zone{},
		automatic: true,
	}
}

func (s *split) Start(update func(state State, inet *Zone, vpn *Zone) Response, stateChanged func(state State, isp Isp)) {
	go run(s.inet)
	go run(s.vpn)
	go s.observe(update, stateChanged)
}

func (s *split) observe(update func(state State, inet *Zone, vpn *Zone) Response, stateChanged func(state State, isp Isp)) {
	laststate := NoConnected
	for true {
		if s.automatic {
			s.router.updateInterfaces()
			ifnets := s.router.getInterfacesWithGateway()
			for _, ifnet := range ifnets {
				s.vpn.host = ifnet.GatewayAddress
			}
		}
		if laststate != s.State {
			laststate = s.State
			isp := fetchISP()
			stateChanged(s.State, isp)
		}
		response := update(laststate, s.inet, s.vpn)
		if response.SplitNow {
			s.mapZone(s.resplit)
			response.SplitNow = false
		}
		if response.Diagnose {
			s.startDiagnose()
			response.Diagnose = false
		}
		if response.AutomaticMode != Nil {
			s.automatic = response.AutomaticMode.bool()
			response.AutomaticMode = Nil
		}
		time.Sleep(time.Second)
		if s.inet.active && s.vpn.active {
			s.State = Connected
			if s.inet.gateway == "" || s.vpn.gateway == "" || s.inet.host == "" || s.vpn.host == "" {
				s.mapZone(s.getRouteConfig)
			}
		} else if s.inet.active && !s.vpn.active {
			s.State = InternetConnected
			s.mapZone(s.getRouteConfig)
		} else if !s.inet.active && s.vpn.active {
			s.State = VpnConnected
			s.mapZone(s.resplit)
		} else if !s.inet.active && !s.vpn.active {
			s.State = NoConnected
			s.inet.gateway = ""
			s.vpn.gateway = ""
			s.vpn.host = ""
		}
	}
}

func (s *split) mapZone(call func() (vpnRouteConfig []routeConfig, inetRouteConfig []routeConfig, ok bool)) {
	vpnGateways, inetGateways, ok := call()
	if ok {
		s.vpn.interfaceName = ""
		s.vpn.gateway = ""
		s.vpn.route = ""
		s.vpn.isDefault = false
		for _, gw := range vpnGateways {
			if gw.gateway == s.vpn.Host() {
				s.vpn.interfaceName += gw.interfaceName + " "
				s.vpn.gateway += gw.gateway + " "
				s.vpn.route += gw.route + " "
				s.vpn.isDefault = gw.isDefault
			}
		}
		s.inet.interfaceName = ""
		s.inet.gateway = ""
		s.inet.route = ""
		s.inet.isDefault = false
		for _, gw := range inetGateways {
			s.inet.interfaceName = assignNonDuplicate(s.inet.interfaceName, gw.interfaceName)
			s.inet.gateway = assignNonDuplicate(s.inet.gateway, gw.gateway)
			s.inet.route = assignNonDuplicate(s.inet.route, gw.route)
			s.inet.isDefault = gw.isDefault
		}
	}
}

func assignNonDuplicate(in string, val string) string {
	if strings.Contains(in, val) {
		return in
	}
	if len(in) > 0 {
		return in + " | " + val
	} else {
		return val
	}
}

func (s *split) resplit() (vpnRouteConfig []routeConfig, inetRouteConfig []routeConfig, ok bool) {
	vpnRouteConfig, inetRouteConfig, ok = s.getRouteConfig()
	log.Printf("Perform split...")
	if ok {
		for _, gw := range vpnRouteConfig {
			cmd := "route -nv add -net " + gw.route + " -interface " + gw.interfaceName
			out, err := exec.Command("/bin/sh", "-c", cmd).Output()
			Debugln(string(out))
			if err != nil {
				log.Printf("Failed to add route. %s", err.Error())
			} else {
				log.Printf("VPN route added -- " + cmd)
			}
		}
		for _, r := range inetRouteConfig {
			cmd := "route change default " + r.gateway
			out, err := exec.Command("/bin/sh", "-c", cmd).Output()
			Debugln(string(out))
			if err != nil {
				log.Printf("Failed to set default route. %s", err.Error())
			} else {
				log.Printf("Default gateway set -- " + cmd)
			}
		}
	}
	log.Printf("Done!")
	return
}

func (s *split) getRouteConfig() (vpnRouteConfig []routeConfig, inetRouteConfig []routeConfig, ok bool) {
	s.router.updateInterfaces()
	ifnets := s.router.getInterfacesWithGateway()
	s.router.updateRoutes()
	routes := s.router.getDefaultRoutes()
	for _, r := range routes {
		found := false
		for _, ifnet := range ifnets {
			if r.gateway == ifnet.GatewayAddress {
				found = true
			}
		}
		if r.isDefault && !r.ipv6 && !found {
			inetRouteConfig = append(inetRouteConfig, routeConfig{
				isDefault:     r.isDefault,
				gateway:       r.gateway,
				route:         "default",
				interfaceName: r.netif,
			})
			ok = true
		}

	}
	for _, ifnet := range ifnets {
		vpnRouteConfig = append(vpnRouteConfig, routeConfig{
			isDefault:     false,
			gateway:       ifnet.GatewayAddress,
			route:         maskIp(ifnet.GatewayAddress+vpnMask) + vpnMask,
			interfaceName: ifnet.Name,
		})
	}

	return
}

// Extracts IP mask from CIDR address.
func maskIp(cidr string) string {
	ip, ipNet, _ := net.ParseCIDR(cidr)
	c := ip.Mask(ipNet.Mask)
	return c.String()
}

func run(status *Zone) { // todo: move to zone
	for true {
		if status.host != "" {
			err := status.Ping()
			if err != nil {
				log.Printf("Ping %s error: %s", status.host, err.Error())
				status.update(0)
			} else {
				if status.isDefault {
					go updateHttpRequest(status)
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func updateHttpRequest(status *Zone) {
	if err := requestNow(status.host); err != nil {
		log.Printf("Request http://%s error: %s", status.host, err.Error())
		status.httpRequest = false
	} else {
		status.httpRequest = true
	}
}

func requestNow(host string) error {
	resp, err := http.Get("http://" + host)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
