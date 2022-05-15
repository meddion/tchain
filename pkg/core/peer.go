package core

import (
	"log"
	"sync"
	"time"
)

type Peer struct {
	Sender
	addr Addr
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

	mtx            sync.RWMutex
	peers          map[Addr]Peer
	shutdown, done chan struct{}
}

func newPeerPool(logger *log.Logger) *peerPool {
	p := &peerPool{
		logger:   logger,
		shutdown: make(chan struct{}, 1),
		done:     make(chan struct{}, 2),
	}

	// TODO: handle errors and resource leaks
	go p.startDiscovery(time.Second * 10)
	go p.checkConnections(time.Second * 10)

	return p
}

func (p *peerPool) Close() error {
	close(p.shutdown)
	for i := 0; i < 2; i++ {
		<-p.done
	}
	close(p.done)

	return nil
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

func (p *peerPool) checkConnections(freq time.Duration) {
	defer func() {
		p.done <- struct{}{}
	}()

	t := time.NewTicker(freq)

	for {
		select {
		case <-t.C:
			if p.NumberOfPeers() < 1 {
				p.logger.Print("No active peers to ping")
				break
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

		case <-p.shutdown:
			return
		}
	}
}

func (p *peerPool) startDiscovery(freq time.Duration) {
	defer func() {
		p.done <- struct{}{}
	}()

	t := time.NewTicker(freq)

	for {
		select {
		case <-t.C:
			newAddrs := p.getNewAddresses()
			p.mtx.Lock()
			for _, addr := range newAddrs {
				peer, err := NewPeer(addr)
				if err != nil {
					p.logger.Printf("On creating a peer connection: %s", err)
					continue
				}

				p.peers[addr] = peer
			}
			p.mtx.Unlock()
		case <-p.shutdown:
			return
		}
	}
}

func (p *peerPool) Peers() []Peer {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	copyPeers := make([]Peer, 0, len(p.peers))
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
