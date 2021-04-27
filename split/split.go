package split

import (
	"github.com/go-ping/ping"
	"log"
	"net"
	"os/exec"
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
	inet   *Zone
	vpn    *Zone
	State  State
	router *router
}

type routeConfig struct {
	isDefault     bool
	gateway       string
	route         string
	interfaceName string
}

type Response struct {
	SplitNow bool
}

func NewSplit() split {
	return split{
		router: newRouter(),
		State:  NoConnected,
		inet:   &Zone{host: "google.com"},
		vpn:    &Zone{},
	}
}

func (s *split) Start(update func(state State, inet *Zone, vpn *Zone) Response) {
	go run(s.inet)
	go run(s.vpn)
	go s.observe(update)
}

func (s *split) observe(update func(state State, inet *Zone, vpn *Zone) Response) {
	laststate := NoConnected
	for true {
		if s.vpn.host == "" {
			s.router.updateInterfaces()
			ifnets := s.router.getInterfacesWithGateway()
			for _, ifnet := range ifnets {
				s.vpn.host = ifnet.GatewayAddress
			}
		}
		if laststate != s.State {
			laststate = s.State
		}
		response := update(laststate, s.inet, s.vpn)
		if response.SplitNow {
			s.mapZone(s.resplit)
			response.SplitNow = false
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

func (s *split) mapZone(call func() (vpnRouteConfig []routeConfig, inetRouteConfig routeConfig, ok bool)) {
	vpnGateways, inetGateway, ok := call()
	if ok {
		for _, gw := range vpnGateways {
			if gw.gateway == s.vpn.Host() {
				s.vpn.interfaceName = gw.interfaceName
				s.vpn.gateway = gw.gateway
				s.vpn.route = gw.route
			}
		}
		s.inet.interfaceName = inetGateway.interfaceName
		s.inet.gateway = inetGateway.gateway
		s.inet.route = inetGateway.route
	}
}

func (s *split) resplit() (vpnRouteConfig []routeConfig, inetRouteConfig routeConfig, ok bool) {
	vpnRouteConfig, inetRouteConfig, ok = s.getRouteConfig()
	if ok {
		for _, gw := range vpnRouteConfig {
			cmd := "route -nv add -net " + gw.route + " -interface " + gw.interfaceName
			_, err := exec.Command("/bin/sh", "-c", cmd).Output()
			if err != nil {
				log.Printf("Failed to add route. %s", err.Error())
			} else {
				log.Printf("VPN route added -- " + cmd)
			}
		}
		cmd := "route change default " + inetRouteConfig.gateway
		_, err := exec.Command("/bin/sh", "-c", cmd).Output()
		if err != nil {
			log.Printf("Failed to set default route. %s", err.Error())
		} else {
			log.Printf("Default gateway set -- " + cmd)
		}
	}
	return
}

func (s *split) getRouteConfig() (vpnRouteConfig []routeConfig, inetRouteConfig routeConfig, ok bool) {
	//vpnAddresses := make(map[string]IfNet)
	inetAddresses := make(map[string]Route)
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
		if !r.ipv6 {
			if !found {
				inetAddresses[r.gateway] = r
			}
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

	for _, r := range inetAddresses {
		if r.isDefault {
			inetRouteConfig = routeConfig{
				isDefault:     r.isDefault,
				gateway:       r.gateway,
				route:         "default",
				interfaceName: r.netif,
			}
			ok = true
			return // take the first pick
		}
	}
	return
}

// Extracts IP mask from CIDR address.
func maskIp(cidr string) string {
	ip, ipNet, _ := net.ParseCIDR(cidr)
	c := ip.Mask(ipNet.Mask)
	return c.String()
}

func run(status *Zone) {
	for true {
		if status.host != "" {
			dur, err := pingNow(status.host)
			if err != nil {
				log.Printf("Ping %s error: %s", status.host, err.Error())
			}
			status.update(dur)
		}
		time.Sleep(time.Second)
	}
}

func pingNow(host string) (duration time.Duration, err error) {
	pinger, perr := ping.NewPinger(host)
	if perr != nil {
		err = perr
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 5
	perr = pinger.Run() // Blocks until finished.
	if perr != nil {
		err = perr
	}
	stats := pinger.Statistics()
	duration = stats.AvgRtt
	return
}
