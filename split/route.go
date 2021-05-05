package split

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

type IfNet struct {
	Name              string
	Status            string
	Multicast         bool
	Broadcast         bool
	Promisc           bool
	PointToPoint      bool
	Running           bool
	Simplex           bool
	Inet              string
	Inet6             string
	HasDefaultGateway bool
	GatewayAddress    string
}

type Route struct {
	dest      string
	isDefault bool
	gateway   string
	flags     string
	netif     string // todo: interfaceName
	ipv6      bool
}

type router struct {
	routes     []Route
	interfaces []IfNet
}

func newRouter() *router {
	var route *router = &router{}
	route.updateRoutes()
	route.updateInterfaces()
	return route
}

func (r *router) getDefaultRoutes() []Route {
	routes := []Route{}
	for _, route := range r.routes {
		if route.isDefault {
			routes = append(routes, route)
		}
	}
	return routes
}

func (r *router) updateRoutes() {
	routes := []Route{}
	out, err := exec.Command("/bin/sh", "-c", "netstat -rn").Output()
	if err != nil {
		log.Printf("Failed to fetch routes. %s", err.Error())
	}
	routing := string(out)
	Debugln(routing)
	rows := strings.Split(routing, "\n")
	inetMode := false
	inet6Mode := false
	for _, row := range rows {
		if row == "Internet:" {
			inetMode = true
			inet6Mode = false
		} else if row == "Internet6:" {
			inetMode = false
			inet6Mode = true
		} else if inetMode || inet6Mode {
			token := strings.Fields(row)
			if len(token) > 3 {
				routes = append(routes, Route{
					dest:      token[0],
					isDefault: token[0] == "default",
					gateway:   token[1],
					flags:     token[2],
					netif:     token[3],
					ipv6:      inet6Mode,
				})
			}
		}
	}
	r.routes = routes
}

func (r *router) getInterfacesWithGateway() []IfNet {
	result := []IfNet{}
	for _, interf := range r.interfaces {
		if interf.HasDefaultGateway {
			result = append(result, interf)
		}
	}
	return result
}

func (r *router) updateInterfaces() {

	interfaces, err := net.Interfaces()

	if err != nil {
		log.Printf("Error fetching interfaces: %s\n", err.Error())
		return
	}

	var interfaceNetworks []IfNet

	for _, i := range interfaces {
		var ifnet IfNet
		ifnet.Name = i.Name
		if strings.Contains(i.Flags.String(), "up") {
			ifnet.Status = "UP"
		} else {
			ifnet.Status = "DOWN"
		}
		if strings.Contains(i.Flags.String(), "multicast") {
			ifnet.Multicast = true
		} else {
			ifnet.Multicast = false
		}
		if strings.Contains(i.Flags.String(), "broadcast") {
			ifnet.Broadcast = true
		} else {
			ifnet.Broadcast = false
		}

		out, err := exec.Command("/bin/sh", "-c", "ifconfig -r "+i.Name).Output()
		Debugln(string(out))
		if err != nil {
			log.Printf("Failed to get interfaces. %s", err.Error())
		}
		ifdetails := string(out)

		rows := strings.Split(ifdetails, "\n")
		for i, row := range rows {
			if i == 0 {
				if strings.Contains(ifdetails, "PROMISC") {
					ifnet.Promisc = true
				} else {
					ifnet.Promisc = false
				}
				if strings.Contains(ifdetails, "POINTOPOINT") {
					ifnet.PointToPoint = true
				} else {
					ifnet.PointToPoint = false
				}
				if strings.Contains(ifdetails, "SIMPLEX") {
					ifnet.Simplex = true
				} else {
					ifnet.Simplex = false
				}
				if strings.Contains(ifdetails, "RUNNING") {
					ifnet.PointToPoint = true
				} else {
					ifnet.PointToPoint = false
				}
			} else {
				rowCut := strings.Trim(row, "\t")
				token := strings.Fields(rowCut)
				if len(token) > 0 {
					switch token[0] {
					case "inet":
						if len(token) > 1 {
							ifnet.Inet = token[1]
						}
						if len(token) > 3 {
							if token[2] == "-->" {
								ifnet.HasDefaultGateway = true
								ifnet.GatewayAddress = token[3]
							}
						}
					case "inet6":
						if len(token) > 1 {
							ifnet.Inet6 = token[1]
						}
					}
				}
			}
		}
		interfaceNetworks = append(interfaceNetworks, ifnet)
	}
	r.interfaces = interfaceNetworks
}
