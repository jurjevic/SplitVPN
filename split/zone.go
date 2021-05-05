package split

import (
	"errors"
	"github.com/go-ping/ping"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MaxResults = 10
)

type Zone struct {
	host          string
	gateway       string
	isDefault     bool
	route         string
	interfaceName string
	httpRequest   bool
	active        bool
	activeSince   time.Time
	time          []time.Duration
	mutex         sync.Mutex
}

func (z *Zone) IsDefault() bool {
	return z.isDefault
}

func (z *Zone) HttpRequest() bool {
	return z.httpRequest
}

func (z *Zone) ActiveSince() time.Time {
	return z.activeSince
}

func (z *Zone) Route() string {
	return z.route
}

func (z *Zone) InterfaceName() string {
	return z.interfaceName
}

func (z *Zone) Gateway() string {
	return z.gateway
}

func (z *Zone) Host() string {
	return z.host
}

func (z *Zone) update(times time.Duration) {
	z.mutex.Lock()
	// add data
	if len(z.time) >= MaxResults {
		z.time = z.time[1:]
	}
	z.time = append(z.time, times)
	z.mutex.Unlock()
	// set active state
	_, _, ok := z.Average()
	if  ok  {
		z.active = true
		if z.activeSince == (time.Time{}) {
			z.activeSince = time.Now()
		}
	} else {
		z.active = false
		z.activeSince = time.Time{}
	}
}

func (z *Zone) Average() (time.Duration, int, bool) {
	z.mutex.Lock()
	var avg time.Duration = 0
	count := 0
	for i, _ := range z.time {
		d := z.time[len(z.time)-1-i]
		if d > 0 {
			avg += d
			count += 1
		}
	}
	z.mutex.Unlock()
	if count == 0 {
		return 0, 0, false
	}
	rate := (100 * count / len(z.time))
	return avg / time.Duration(count), rate, true
}

func (z *Zone) Ping2() error {
	if z.host == "" {
		return errors.New("no host defined to ping")
	}
	Debugf("Ping: %s\n", z.host)
	pinger, err := ping.NewPinger(z.host)
	if err != nil {
		return err
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 5
	perr := pinger.Run() // Blocks until finished.
	if perr != nil {
		return perr
	}
	stats := pinger.Statistics()
	if stats.AvgRtt == 0 {
		Debugf("...ping was not successful!")
	}
	z.update(stats.AvgRtt)
	return nil
}

func (z *Zone) Ping() error {
	if z.host == "" {
		return errors.New("no host defined to ping")
	}
	cmd := "ping -n -o -c 1 -t 5 " + z.host
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	Debugln(string(out))
	if err != nil {
		log.Printf("Failed to ping host. %s", z.host)
		return err
	}

	str := string(out)
	idxTime := strings.Index(str, "time=")
	idxMs := strings.Index(str, " ms")
	str = str[idxTime+5:idxMs]
	f, ferr := strconv.ParseFloat(str,64)
	if ferr != nil {
		log.Printf("Failed to get ping time for host. %s", z.host)
		return err
	}
	d := int(f)
	z.update(time.Duration(d)*time.Millisecond)
	return nil
}
