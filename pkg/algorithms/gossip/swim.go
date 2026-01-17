package gossip

import (
	"math/rand"
	"sync"
	"time"
)

// State represents the state of a member.
type State int

const (
	Alive State = iota
	Suspect
	Dead
)

// Member represents a node in the cluster.
type Member struct {
	ID          string
	Address     string
	State       State
	Incarnation uint64
	LastUpdate  time.Time
}

// Config holds configuration for the Gossip protocol.
type Config struct {
	BindAddress    string
	ID             string
	ProtocolPeriod time.Duration
	PingTimeout    time.Duration
	SuspectTimeout time.Duration
	PingReqK       int // Number of members to ask to ping a suspect
}

// Protocol implements a basic SWIM-style gossip protocol logic.
// Transport is abstracted away; users must hook up networking.
type Protocol struct {
	config  Config
	members map[string]*Member
	mu      sync.RWMutex

	// Events
	events chan Event

	// Transport hook
	Transport Transport
}

type Event struct {
	Type   string // "Join", "Leave", "Fail", "Update"
	Member Member
}

// Transport abstracts network operations.
type Transport interface {
	Ping(target string) (bool, error)
	PingReq(target string, proxy string) (bool, error)
}

// New creates a new Gossip Protocol instance.
func New(config Config, transport Transport) *Protocol {
	if config.ProtocolPeriod == 0 {
		config.ProtocolPeriod = 1 * time.Second
	}
	if config.PingReqK == 0 {
		config.PingReqK = 3
	}

	return &Protocol{
		config:    config,
		members:   make(map[string]*Member),
		events:    make(chan Event, 100),
		Transport: transport,
	}
}

// Start starts the gossip loop.
func (p *Protocol) Start() {
	go p.loop()
}

// Join adds a member to the local list (seeds).
func (p *Protocol) Join(id, address string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.members[id]; !exists {
		p.members[id] = &Member{
			ID:          id,
			Address:     address,
			State:       Alive,
			Incarnation: 0,
			LastUpdate:  time.Now(),
		}
	}
}

// Members returns the list of known members.
func (p *Protocol) Members() []Member {
	p.mu.RLock()
	defer p.mu.RUnlock()

	list := make([]Member, 0, len(p.members))
	for _, m := range p.members {
		list = append(list, *m)
	}
	return list
}

func (p *Protocol) loop() {
	ticker := time.NewTicker(p.config.ProtocolPeriod)
	for range ticker.C {
		p.probe()
	}
}

// probe performs a gossip round (random ping).
func (p *Protocol) probe() {
	target := p.selectRandomMember()
	if target == nil {
		return
	}

	success, _ := p.Transport.Ping(target.Address)
	if !success {
		// Indirect Ping (PingReq)
		if !p.pingReq(target) {
			p.markSuspect(target)
		}
	} else {
		// If success and was suspect, refute?
		// Simpler model: Alive update
		p.markAlive(target)
	}
}

func (p *Protocol) selectRandomMember() *Member {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if len(p.members) == 0 {
		return nil
	}

	// inefficient for large lists, but functional for now
	keys := make([]string, 0, len(p.members))
	for k := range p.members {
		if k != p.config.ID {
			keys = append(keys, k)
		}
	}

	if len(keys) == 0 {
		return nil
	}

	id := keys[rand.Intn(len(keys))]
	return p.members[id]
}

func (p *Protocol) pingReq(target *Member) bool {
	p.mu.RLock()
	proxies := make([]*Member, 0)
	for _, m := range p.members {
		if m.ID != p.config.ID && m.ID != target.ID && m.State == Alive {
			proxies = append(proxies, m)
		}
	}
	p.mu.RUnlock()

	// Shuffle and pick K
	rand.Shuffle(len(proxies), func(i, j int) { proxies[i], proxies[j] = proxies[j], proxies[i] })
	k := p.config.PingReqK
	if len(proxies) < k {
		k = len(proxies)
	}

	for i := 0; i < k; i++ {
		proxy := proxies[i]
		ok, _ := p.Transport.PingReq(target.Address, proxy.Address)
		if ok {
			return true
		}
	}
	return false
}

func (p *Protocol) markSuspect(m *Member) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if m.State == Alive {
		m.State = Suspect
		m.LastUpdate = time.Now()
		// Emit event, schedule deletion/dead transition
		// In full SWIM, we multicast this.
	} else if m.State == Suspect {
		// Time out suspect -> Dead
		if time.Since(m.LastUpdate) > p.config.SuspectTimeout {
			m.State = Dead
			delete(p.members, m.ID)
		}
	}
}

func (p *Protocol) markAlive(m *Member) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if m.State != Alive {
		m.State = Alive
		m.Incarnation++
		m.LastUpdate = time.Now()
	}
}

// Events returns the event channel.
func (p *Protocol) Events() <-chan Event {
	return p.events
}
