package roundrobin

import (
	"errors"
	"net"
	"regexp"
	"sync"
	"sync/atomic"
)

var (
	ErrNoBackends = errors.New("no backends specified")
	ErrNotValid   = errors.New("one or more addresses are not valid")
	ErrDuplicates = errors.New("found duplicates")
)

type Backend struct {
	Address string
	IsAlive bool
	Mu      sync.RWMutex
}

type Service struct {
	Backends []*Backend
	next     uint32
}

func findDuplicates(addresses []string) []string {
	found := make(map[string]bool)
	duplicates := []string{}

	for _, address := range addresses {
		if _, value := found[address]; value {
			duplicates = append(duplicates, address)

		} else {
			found[address] = true

		}
	}
	return duplicates
}

func validateAddresses(addresses []string) error {
	for _, address := range addresses {
		ip, port, err := net.SplitHostPort(address)
		if err != nil {
			return err
		}

		ipValid := net.ParseIP(ip)

		re := regexp.MustCompile(`^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`)
		portValid := re.FindAllString(port, -1)

		if len(ipValid) == 0 || len(portValid) == 0 {
			return ErrNotValid
		}
	}
	return nil
}

func validate(addresses []string) error {
	if len(addresses) == 0 {
		return ErrNoBackends
	}

	err := validateAddresses(addresses)
	if err != nil {
		return err
	}

	duplicates := findDuplicates(addresses)
	if len(duplicates) > 0 {
		return ErrDuplicates
	}

	return nil
}

// NewService returns a service
func NewService(addresses []string) (*Service, error) {
	err := validate(addresses)
	if err != nil {
		return nil, err
	}

	var backends []*Backend
	for _, address := range addresses {
		backend := Backend{
			Address: address,
			IsAlive: true,
			Mu:      sync.RWMutex{},
		}

		backends = append(backends, &backend)
	}

	return &Service{
		Backends: backends,
	}, nil
}

// Next returns the next url from a service
func (s *Service) Next() *Backend {
	n := atomic.AddUint32(&s.next, 1)
	return s.Backends[(int(n)-1)%len(s.Backends)]
}
