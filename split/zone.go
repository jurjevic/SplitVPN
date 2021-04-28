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
	activeSince   time.Time
	time          []time.Duration
	mutex         sync.Mutex
}

func (s *Zone) ActiveSince() time.Time {
	return s.activeSince
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

func (s *Zone) update(times time.Duration) {
	s.mutex.Lock()
	if times > 0 {
		s.active = true
		if s.activeSince == (time.Time{}) {
			s.activeSince = time.Now()
		}
	} else {
		s.activeSince = time.Time{}
	}
	if len(s.time) >= maxResults {
		s.time = s.time[1:]
	}
	s.time = append(s.time, times)
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
