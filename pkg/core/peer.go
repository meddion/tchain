package core

import "sync"

type Peer struct {
	Sender
}

type peerPool struct {
	lck   sync.RWMutex
	peers []Peer
}

func (p *peerPool) Peers() []Peer {
	p.lck.RLock()
	defer p.lck.RUnlock()

	copyPeers := make([]Peer, len(p.peers))
	copy(copyPeers, p.peers)

	return copyPeers
}

func (p *peerPool) NumberOfPeers() int {
	p.lck.RLock()
	defer p.lck.RUnlock()

	return len(p.peers)
}

func (p *peerPool) SendToPeers(fun func(Sender) error) <-chan error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, len(p.peers))
	)

	wg.Add(len(p.peers))

	p.lck.RLock()
	for _, s := range p.peers {
		go func(s Sender) {
			defer wg.Done()

			if err := fun(s); err != nil {
				errChan <- err
			}
		}(s)
	}
	p.lck.RUnlock()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}
