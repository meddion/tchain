package core

import (
	"log"
	"strings"
	"sync"
	"time"
)

type Peer struct {
	Sender
}

func NewPeer(ip, port string) (Peer, error) {
	s, err := NewSender(ip, port)
	if err != nil {
		return Peer{}, err
	}

	return Peer{s}, nil
}

type peerPool struct {
	logger *log.Logger

	mtx   sync.RWMutex
	peers map[string]Peer
}

func newPeerPool() *peerPool {
	p := &peerPool{}
	// TODO: handle errors and resource leaks
	go p.startDiscovery()

	return p
}

func getEndpoints() []string {
	return nil
}

func (p *peerPool) startDiscovery() {
	t := time.NewTicker(time.Second * 10)
	for range t.C {
		endpoints := getEndpoints()
		p.mtx.Lock()
		for _, e := range endpoints {
			if _, exists := p.peers[e]; !exists {
				addrPair := strings.Split(e, ".")
				if len(addrPair) < 2 {
					p.logger.Printf("On extracting ip and port from an endpoint: %s", addrPair)
					continue
				}
				peer, err := NewPeer(addrPair[0], addrPair[1])
				if err != nil {
					p.logger.Printf("On creating a peer connection: %s", err)
					continue
				}

				p.peers[e] = peer
			}
		}
		p.mtx.Unlock()
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

func (p *peerPool) SendToPeers(fun func(Sender) error) <-chan error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(p.peers))
	)

	wg.Add(len(p.peers))

	p.mtx.RLock()
	for _, s := range p.peers {
		go func(s Sender) {
			defer wg.Done()

			if err := fun(s); err != nil {
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
