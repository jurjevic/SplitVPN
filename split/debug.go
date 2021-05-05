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
	PrintTitle("Diagnose")
	println("Version: ", Version)
	PrintTitle("Route configs")
	vpnRC, inetRC, ok := s.getRouteConfig()
	if !ok {
		println("...route not ok")
	}
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
}
