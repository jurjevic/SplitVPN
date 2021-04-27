package split

import (
	"sync"
	"time"
)

const (
	maxResults = 10
)

type Zone struct {
	host          string
	gateway       string
	route         string
	interfaceName string
	active        bool
	time          []time.Duration
	mutex         sync.Mutex
}

func (s *Zone) Route() string {
	return s.route
}

func (s *Zone) InterfaceName() string {
	return s.interfaceName
}

func (s *Zone) Gateway() string {
	return s.gateway
}

func (s *Zone) Host() string {
	return s.host
}

func (s *Zone) update(time time.Duration) {
	s.mutex.Lock()
	s.active = time > 0
	if len(s.time) >= maxResults {
		s.time = s.time[1:]
	}
	s.time = append(s.time, time)
	s.mutex.Unlock()
}

func (s *Zone) Average() (time.Duration, bool) {
	if s.active == false {
		return 0, false
	}
	s.mutex.Lock()
	var avg time.Duration = 0
	count := 0
	for i, _ := range s.time {
		d := s.time[len(s.time)-1-i]
		if d > 0 {
			avg += d
			count += 1
		}
	}
	s.mutex.Unlock()
	if count == 0 {
		return 0, false
	}
	return avg / time.Duration(count), true
}
