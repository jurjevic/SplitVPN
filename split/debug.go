package split

import (
	"fmt"
	"log"
)

var DebugFlag bool = false

var (
	Version string
)


func Debug(v ...interface{}) {
	if DebugFlag {
		log.Output(2, fmt.Sprint(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if DebugFlag {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debugln(v ...interface{}) {
	if DebugFlag {
		log.Output(2, fmt.Sprintln(v...))
	}
}

func PrintTitle(title string) {
	println()
	length := len(title) + 4
	printChar("░", length)
	println("░ " + title + " ░")
	printChar("░", length)
	println()
}

func printChar(char string, length int) {
	for i := 0; i < length; i++ {
		print(char)
	}
	println()
}

func (s split) startDiagnose() {
	tempAutomatic := s.automatic
	s.automatic = false
	tempDebugFlag := DebugFlag
	DebugFlag = true

	PrintTitle("Diagnose")
	println("Version: ", Version)
	println("Debug Mode: ", tempDebugFlag)
	println("Automatic Mode: ", tempAutomatic)

	PrintTitle("Network configuration")
	vpnRC, inetRC, ok := s.getRouteConfig()
	if !ok {
		println("...route not ok")
	}
	PrintTitle("Processing results")
	for _, rc := range vpnRC {
		println("VNET gateway:", rc.gateway)
		println(".......route:", rc.route)
		println(".....default:", rc.isDefault)
		println("...interface:", rc.interfaceName)
	}
	for _, rc := range inetRC {
		println("INET gateway:", rc.gateway)
		println(".......route:", rc.route)
		println(".....default:", rc.isDefault)
		println("...interface:", rc.interfaceName)
	}

	PrintTitle("VPN Zone")
	s.vpn.diagnose()
	PrintTitle("INET Zone")
	s.inet.diagnose()

	PrintTitle("Internet Service Provider - System proxy")
	fetchISP()

	PrintTitle("Internet Service Provider - No proxy")
	fetchNoProxyISP()


	println()
	printChar("▬", 80)
	println()
	DebugFlag = tempDebugFlag
	s.automatic = tempAutomatic
}
