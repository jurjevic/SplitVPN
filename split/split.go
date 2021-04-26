package split

import (
	"github.com/go-ping/ping"
	"os/exec"
	"time"
)

type State int

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

func NewSplit() split {
	return split{
		router: newRouter(),
		State:  NoConnected,
		inet:   &Zone{host: "google.com"},
		vpn:    &Zone{},
	}
}

func (s *split) Start(update func (state State, inet *Zone, vpn *Zone)) {
	go run(s.inet)
	go run(s.vpn)
	go s.observe(update)
}

func (s *split) observe(update func (state State, inet *Zone, vpn *Zone)) {
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
		update(laststate, s.inet, s.vpn)
		time.Sleep(time.Second)
		if s.inet.active && s.vpn.active {
			s.State = Connected
		} else if s.inet.active && !s.vpn.active {
			s.State = InternetConnected
		} else if !s.inet.active && s.vpn.active {
			s.State = VpnConnected
			s.reconnect()
		}  else if !s.inet.active && !s.vpn.active {
			s.State = NoConnected
			s.reconnect()
		}
	}
}

func (s *split) reconnect() {
	vpnAddresses := make(map[string]IfNet)
	inetAddresses := make(map[string]Route)
	s.router.updateInterfaces()
	ifnets := s.router.getInterfacesWithGateway()
	s.router.updateRoutes()
	routes := s.router.getDefaultRoutes()
	for _, r := range routes {
		found := false
		if !r.ipv6 {
			for _, ifnet := range ifnets {
				if r.gateway == ifnet.GatewayAddress {
					vpnAddresses[ifnet.GatewayAddress] = ifnet
					found = true
				}
			}
			if !found {
				inetAddresses[r.gateway] = r
			}
		}
	}
	if len(vpnAddresses) > 0 && len(inetAddresses) > 0 {
		for _, ifnet := range vpnAddresses {
			_, err := exec.Command("/bin/sh", "-c", "route -nv add -net " + ifnet.GatewayAddress +"/8 -interface "+ifnet.Name).Output()
			if err != nil {
				println(err.Error()) // todo:
			}
		}
		for _, r := range inetAddresses {
			_, err := exec.Command("/bin/sh", "-c", "route change default " + r.gateway).Output()
			if err != nil {
				println(err.Error()) // todo:
			}
		}
	}
}

func run(status *Zone) {
	for true {
		dur, _ := pingNow(status.host)
		status.update(dur)
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
