package swarm

import (
	"runtime"
	"sync"
	"testing"
	"time"

	ci "github.com/ipfs/go-libp2p/testutil/ci"
	peer "gx/ipfs/QmY1xNhBfF9xA1pmD8yejyQAyd77K68qNN6JPM1CN2eiRu/go-libp2p-peer"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	ma "gx/ipfs/QmcobAGsCjYt5DXoq9et9L8yR8er7o7Cu3DTvpaq12jYSz/go-multiaddr"
)

func TestSimultOpen(t *testing.T) {

	t.Parallel()

	ctx := context.Background()
	swarms := makeSwarms(ctx, t, 2)

	// connect everyone
	{
		var wg sync.WaitGroup
		connect := func(s *Swarm, dst peer.ID, addr ma.Multiaddr) {
			// copy for other peer
			log.Debugf("TestSimultOpen: connecting: %s --> %s (%s)", s.local, dst, addr)
			s.peers.AddAddr(dst, addr, peer.PermanentAddrTTL)
			if _, err := s.Dial(ctx, dst); err != nil {
				t.Fatal("error swarm dialing to peer", err)
			}
			wg.Done()
		}

		log.Info("Connecting swarms simultaneously.")
		wg.Add(2)
		go connect(swarms[0], swarms[1].local, swarms[1].ListenAddresses()[0])
		go connect(swarms[1], swarms[0].local, swarms[0].ListenAddresses()[0])
		wg.Wait()
	}

	for _, s := range swarms {
		s.Close()
	}
}

func TestSimultOpenMany(t *testing.T) {
	// t.Skip("very very slow")

	addrs := 20
	rounds := 10
	if ci.IsRunning() || runtime.GOOS == "darwin" {
		// osx has a limit of 256 file descriptors
		addrs = 10
		rounds = 5
	}
	SubtestSwarm(t, addrs, rounds)
}

func TestSimultOpenFewStress(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// t.Skip("skipping for another test")
	t.Parallel()

	msgs := 40
	swarms := 2
	rounds := 10
	// rounds := 100

	for i := 0; i < rounds; i++ {
		SubtestSwarm(t, swarms, msgs)
		<-time.After(10 * time.Millisecond)
	}
}
