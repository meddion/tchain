package core

import (
	"log"
	"sync"
	"time"
)

type Peer struct {
	addr Addr
	Sender
}

func NewPeer(addr Addr) (Peer, error) {
	s, err := NewSender(addr)
	if err != nil {
		return Peer{}, err
	}

	return Peer{Sender: s, addr: addr}, nil
}

func (p Peer) Addr() Addr {
	return p.addr
}

var _ PeerPool = &peerPool{}

type peerPool struct {
	logger *log.Logger
	mtx    sync.RWMutex
	peers  map[Addr]Peer

	shutdown, done chan struct{}
	processCounter uint8
}

// Set time parameters to zero to disable it
func newPeerPool(logger *log.Logger, peerDiscoveryTime, isAliveTime time.Duration) *peerPool {
	p := &peerPool{
		logger:   logger,
		peers:    make(map[Addr]Peer),
		shutdown: make(chan struct{}, 1),
		done:     make(chan struct{}, 2),
	}

	job := func(f func(), freq time.Duration) {
		defer func() { p.done <- struct{}{} }()
		t := time.NewTicker(freq)

		for {
			select {
			case <-t.C:
				f()
			case <-p.shutdown:
				return
			}
		}
	}

	// TODO: handle errors and resource leaks
	if peerDiscoveryTime != 0 {
		go job(p.discoverNewPeers, peerDiscoveryTime)
		p.processCounter++
	}

	if isAliveTime != 0 {
		go job(p.pingConnections, isAliveTime)
		p.processCounter++
	}

	return p
}

func (p *peerPool) Close() {
	close(p.shutdown)
	for i := uint8(0); i < p.processCounter; i++ {
		<-p.done
	}
	close(p.done)
}

func (p *peerPool) add(peer Peer) {
	p.peers[peer.addr] = peer
}

func (p *peerPool) Add(peer Peer) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.add(peer)
}

func (p *peerPool) pingConnections() {
	if p.NumberOfPeers() <= 0 {
		p.logger.Print("No active peers to ping")
		return
	}
	// Check health of connections
	notAlivePeers := make(chan Peer, p.NumberOfPeers())
	errorsChan := p.SendToPeers(func(peer Peer) error {
		if err := peer.SendIsAlive(); err != nil {
			notAlivePeers <- peer

			p.logger.Printf("On pinging a peer (%s): %s", peer.addr, err)
			return nil
		}

		return nil
	})

	go func() {
		for err := range errorsChan {
			if err != nil {
				p.logger.Printf("On sending a msg to a peer: %s", err)
			}
		}
		close(notAlivePeers)
	}()

	p.mtx.Lock()
	for peer := range notAlivePeers {
		delete(p.peers, peer.Addr())
	}
	p.mtx.Unlock()
}

func (p *peerPool) getNewAddresses() (addrs []Addr) {
	newAddrs := make(chan Addr, 20)
	errorsChan := p.SendToPeers(func(peer Peer) error {
		resp, err := peer.SendPeersDiscovery()
		if err != nil {
			p.logger.Printf("On getting a peer list from %s", peer.addr)
			return nil
		}

		p.mtx.RLock()
		for _, addr := range resp.addrs {
			if _, exists := p.peers[addr]; !exists {
				newAddrs <- addr
			}
		}
		p.mtx.RUnlock()

		return nil
	})

	go func() {
		for err := range errorsChan {
			if err != nil {
				p.logger.Printf("On sending a msg to a peer: %s", err)
			}
		}
		close(newAddrs)
	}()

	for addr := range newAddrs {
		addrs = append(addrs, addr)
	}

	return addrs
}

func (p *peerPool) discoverNewPeers() {
	newAddrs := p.getNewAddresses()
	p.mtx.Lock()
	for _, addr := range newAddrs {
		peer, err := NewPeer(addr)
		if err != nil {
			p.logger.Printf("On creating a peer connection: %s", err)
		}

		p.add(peer)
	}
	p.mtx.Unlock()
}

func (p *peerPool) Peers() []Peer {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	copyPeers := make([]Peer, len(p.peers))
	for _, peer := range p.peers {
		copyPeers = append(copyPeers, peer)
	}

	return copyPeers
}

func (p *peerPool) NumberOfPeers() int {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	return len(p.peers)
}

func (p *peerPool) SendToPeers(fun func(Peer) error) <-chan error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(p.peers))
	)

	wg.Add(len(p.peers))

	p.mtx.RLock()
	for _, s := range p.peers {
		go func(p Peer) {
			defer wg.Done()

			if err := fun(p); err != nil {
				errChan <- err
			}
		}(s)
	}
	p.mtx.RUnlock()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}
