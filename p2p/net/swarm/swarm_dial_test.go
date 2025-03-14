package swarm

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"net"
	"sort"
	"testing"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/zbliujia/go-libp2p/core/crypto"
	"github.com/zbliujia/go-libp2p/core/network"
	"github.com/zbliujia/go-libp2p/core/peer"
	"github.com/zbliujia/go-libp2p/core/peerstore"
	"github.com/zbliujia/go-libp2p/core/test"
	"github.com/zbliujia/go-libp2p/p2p/host/eventbus"
	"github.com/zbliujia/go-libp2p/p2p/host/peerstore/pstoremem"
	libp2pquic "github.com/zbliujia/go-libp2p/p2p/transport/quic"
	"github.com/zbliujia/go-libp2p/p2p/transport/quicreuse"
	"github.com/zbliujia/go-libp2p/p2p/transport/tcp"
	"github.com/zbliujia/go-libp2p/p2p/transport/websocket"
	libp2pwebtransport "github.com/zbliujia/go-libp2p/p2p/transport/webtransport"

	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/stretchr/testify/require"
)

func TestAddrsForDial(t *testing.T) {
	mockResolver := madns.MockResolver{IP: make(map[string][]net.IPAddr)}
	ipaddr, err := net.ResolveIPAddr("ip4", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	mockResolver.IP["example.com"] = []net.IPAddr{*ipaddr}

	resolver, err := madns.NewResolver(madns.WithDomainResolver("example.com", &mockResolver))
	if err != nil {
		t.Fatal(err)
	}

	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	id, err := peer.IDFromPrivateKey(priv)
	require.NoError(t, err)

	ps, err := pstoremem.NewPeerstore()
	require.NoError(t, err)
	ps.AddPubKey(id, priv.GetPublic())
	ps.AddPrivKey(id, priv)
	t.Cleanup(func() { ps.Close() })

	tpt, err := websocket.New(nil, &network.NullResourceManager{})
	require.NoError(t, err)
	s, err := NewSwarm(id, ps, eventbus.NewBus(), WithMultiaddrResolver(resolver))
	require.NoError(t, err)
	defer s.Close()
	err = s.AddTransport(tpt)
	require.NoError(t, err)

	otherPeer := test.RandPeerIDFatal(t)

	ps.AddAddr(otherPeer, ma.StringCast("/dns4/example.com/tcp/1234/wss"), time.Hour)

	ctx := context.Background()
	mas, _, err := s.addrsForDial(ctx, otherPeer)
	require.NoError(t, err)

	require.NotZero(t, len(mas))
}

func TestDedupAddrsForDial(t *testing.T) {
	mockResolver := madns.MockResolver{IP: make(map[string][]net.IPAddr)}
	ipaddr, err := net.ResolveIPAddr("ip4", "1.2.3.4")
	if err != nil {
		t.Fatal(err)
	}
	mockResolver.IP["example.com"] = []net.IPAddr{*ipaddr}

	resolver, err := madns.NewResolver(madns.WithDomainResolver("example.com", &mockResolver))
	if err != nil {
		t.Fatal(err)
	}

	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	id, err := peer.IDFromPrivateKey(priv)
	require.NoError(t, err)

	ps, err := pstoremem.NewPeerstore()
	require.NoError(t, err)
	ps.AddPubKey(id, priv.GetPublic())
	ps.AddPrivKey(id, priv)
	t.Cleanup(func() { ps.Close() })

	s, err := NewSwarm(id, ps, eventbus.NewBus(), WithMultiaddrResolver(resolver))
	require.NoError(t, err)
	defer s.Close()

	tpt, err := tcp.NewTCPTransport(nil, &network.NullResourceManager{})
	require.NoError(t, err)
	err = s.AddTransport(tpt)
	require.NoError(t, err)

	otherPeer := test.RandPeerIDFatal(t)

	ps.AddAddr(otherPeer, ma.StringCast("/dns4/example.com/tcp/1234"), time.Hour)
	ps.AddAddr(otherPeer, ma.StringCast("/ip4/1.2.3.4/tcp/1234"), time.Hour)

	ctx := context.Background()
	mas, _, err := s.addrsForDial(ctx, otherPeer)
	require.NoError(t, err)

	require.Equal(t, 1, len(mas))
}

func newTestSwarmWithResolver(t *testing.T, resolver *madns.Resolver) *Swarm {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)
	id, err := peer.IDFromPrivateKey(priv)
	require.NoError(t, err)
	ps, err := pstoremem.NewPeerstore()
	require.NoError(t, err)
	ps.AddPubKey(id, priv.GetPublic())
	ps.AddPrivKey(id, priv)
	t.Cleanup(func() { ps.Close() })
	s, err := NewSwarm(id, ps, eventbus.NewBus(), WithMultiaddrResolver(resolver))
	require.NoError(t, err)
	t.Cleanup(func() {
		s.Close()
	})

	// Add a tcp transport so that we know we can dial a tcp multiaddr and we don't filter it out.
	tpt, err := tcp.NewTCPTransport(nil, &network.NullResourceManager{})
	require.NoError(t, err)
	err = s.AddTransport(tpt)
	require.NoError(t, err)

	connmgr, err := quicreuse.NewConnManager(quic.StatelessResetKey{})
	require.NoError(t, err)
	quicTpt, err := libp2pquic.NewTransport(priv, connmgr, nil, nil, &network.NullResourceManager{})
	require.NoError(t, err)
	err = s.AddTransport(quicTpt)
	require.NoError(t, err)

	wtTpt, err := libp2pwebtransport.New(priv, nil, connmgr, nil, &network.NullResourceManager{})
	require.NoError(t, err)
	err = s.AddTransport(wtTpt)
	require.NoError(t, err)

	wsTpt, err := websocket.New(nil, &network.NullResourceManager{})
	require.NoError(t, err)
	err = s.AddTransport(wsTpt)
	require.NoError(t, err)

	return s
}

func TestAddrResolution(t *testing.T) {
	ctx := context.Background()

	p1 := test.RandPeerIDFatal(t)
	p2 := test.RandPeerIDFatal(t)
	addr1 := ma.StringCast("/dnsaddr/example.com")
	addr2 := ma.StringCast("/ip4/192.0.2.1/tcp/123")

	p2paddr2 := ma.StringCast("/ip4/192.0.2.1/tcp/123/p2p/" + p1.String())
	p2paddr3 := ma.StringCast("/ip4/192.0.2.1/tcp/123/p2p/" + p2.String())

	backend := &madns.MockResolver{
		TXT: map[string][]string{"_dnsaddr.example.com": {
			"dnsaddr=" + p2paddr2.String(), "dnsaddr=" + p2paddr3.String(),
		}},
	}
	resolver, err := madns.NewResolver(madns.WithDefaultResolver(backend))
	require.NoError(t, err)

	s := newTestSwarmWithResolver(t, resolver)

	s.peers.AddAddr(p1, addr1, time.Hour)

	tctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	mas, _, err := s.addrsForDial(tctx, p1)
	require.NoError(t, err)

	require.Len(t, mas, 1)
	require.Contains(t, mas, addr2)

	addrs := s.peers.Addrs(p1)
	require.Len(t, addrs, 2)
	require.Contains(t, addrs, addr1)
	require.Contains(t, addrs, addr2)
}

func TestAddrResolutionRecursive(t *testing.T) {
	ctx := context.Background()

	p1, err := test.RandPeerID()
	if err != nil {
		t.Error(err)
	}
	p2, err := test.RandPeerID()
	if err != nil {
		t.Error(err)
	}
	addr1 := ma.StringCast("/dnsaddr/example.com")
	addr2 := ma.StringCast("/ip4/192.0.2.1/tcp/123")
	p2paddr1 := ma.StringCast("/dnsaddr/example.com/p2p/" + p1.String())
	p2paddr2 := ma.StringCast("/dnsaddr/example.com/p2p/" + p2.String())
	p2paddr1i := ma.StringCast("/dnsaddr/foo.example.com/p2p/" + p1.String())
	p2paddr2i := ma.StringCast("/dnsaddr/bar.example.com/p2p/" + p2.String())
	p2paddr1f := ma.StringCast("/ip4/192.0.2.1/tcp/123/p2p/" + p1.String())

	backend := &madns.MockResolver{
		TXT: map[string][]string{
			"_dnsaddr.example.com": {
				"dnsaddr=" + p2paddr1i.String(),
				"dnsaddr=" + p2paddr2i.String(),
			},
			"_dnsaddr.foo.example.com": {
				"dnsaddr=" + p2paddr1f.String(),
			},
			"_dnsaddr.bar.example.com": {
				"dnsaddr=" + p2paddr2i.String(),
			},
		},
	}
	resolver, err := madns.NewResolver(madns.WithDefaultResolver(backend))
	if err != nil {
		t.Fatal(err)
	}

	s := newTestSwarmWithResolver(t, resolver)

	pi1, err := peer.AddrInfoFromP2pAddr(p2paddr1)
	require.NoError(t, err)

	tctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	s.Peerstore().AddAddrs(pi1.ID, pi1.Addrs, peerstore.TempAddrTTL)
	_, _, err = s.addrsForDial(tctx, p1)
	require.NoError(t, err)

	addrs1 := s.Peerstore().Addrs(pi1.ID)
	require.Len(t, addrs1, 2)
	require.Contains(t, addrs1, addr1)
	require.Contains(t, addrs1, addr2)

	pi2, err := peer.AddrInfoFromP2pAddr(p2paddr2)
	require.NoError(t, err)

	s.Peerstore().AddAddrs(pi2.ID, pi2.Addrs, peerstore.TempAddrTTL)
	_, _, err = s.addrsForDial(tctx, p2)
	// This never resolves to a good address
	require.Equal(t, ErrNoGoodAddresses, err)

	addrs2 := s.Peerstore().Addrs(pi2.ID)
	require.Len(t, addrs2, 1)
	require.Contains(t, addrs2, addr1)
}

func TestAddrsForDialFiltering(t *testing.T) {
	q1 := ma.StringCast("/ip4/1.2.3.4/udp/1/quic-v1")
	q1v1 := ma.StringCast("/ip4/1.2.3.4/udp/1/quic-v1")
	wt1 := ma.StringCast("/ip4/1.2.3.4/udp/1/quic-v1/webtransport/")

	q2 := ma.StringCast("/ip4/1.2.3.4/udp/2/quic-v1")
	q2v1 := ma.StringCast("/ip4/1.2.3.4/udp/2/quic-v1")
	wt2 := ma.StringCast("/ip4/1.2.3.4/udp/2/quic-v1/webtransport/")

	q3 := ma.StringCast("/ip4/1.2.3.4/udp/3/quic-v1")

	t1 := ma.StringCast("/ip4/1.2.3.4/tcp/1")
	ws1 := ma.StringCast("/ip4/1.2.3.4/tcp/1/ws")

	unSpecQ := ma.StringCast("/ip4/0.0.0.0/udp/2/quic-v1")
	unSpecT := ma.StringCast("/ip6/::/tcp/2/")

	resolver, err := madns.NewResolver(madns.WithDefaultResolver(&madns.MockResolver{}))
	require.NoError(t, err)
	s := newTestSwarmWithResolver(t, resolver)
	ourAddrs := s.ListenAddresses()

	testCases := []struct {
		name   string
		input  []ma.Multiaddr
		output []ma.Multiaddr
	}{
		{
			name:   "quic-filtered",
			input:  []ma.Multiaddr{q1, q1v1, q2, q2v1, q3},
			output: []ma.Multiaddr{q1v1, q2v1, q3},
		},
		{
			name:   "webtransport-filtered",
			input:  []ma.Multiaddr{q1, q1v1, wt1, wt2},
			output: []ma.Multiaddr{q1v1, wt2},
		},
		{
			name:   "all",
			input:  []ma.Multiaddr{q1, q1v1, wt1, q2, q2v1, wt2, t1, ws1},
			output: []ma.Multiaddr{q1v1, q2v1, t1},
		},
		{
			name:   "our-addrs-filtered",
			input:  append([]ma.Multiaddr{q1}, ourAddrs...),
			output: []ma.Multiaddr{q1},
		},
		{
			name:   "unspecified-filtered",
			input:  []ma.Multiaddr{q1v1, t1, unSpecQ, unSpecT},
			output: []ma.Multiaddr{q1v1, t1},
		},
	}

	ctx := context.Background()
	p1 := test.RandPeerIDFatal(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s.Peerstore().ClearAddrs(p1)
			s.Peerstore().AddAddrs(p1, tc.input, peerstore.PermanentAddrTTL)
			result, _, err := s.addrsForDial(ctx, p1)
			require.NoError(t, err)
			sort.Slice(result, func(i, j int) bool { return bytes.Compare(result[i].Bytes(), result[j].Bytes()) < 0 })
			sort.Slice(tc.output, func(i, j int) bool { return bytes.Compare(tc.output[i].Bytes(), tc.output[j].Bytes()) < 0 })
			if len(result) != len(tc.output) {
				t.Fatalf("output mismatch got: %s want: %s", result, tc.output)
			}
			for i := 0; i < len(result); i++ {
				if !result[i].Equal(tc.output[i]) {
					t.Fatalf("output mismatch got: %s want: %s", result, tc.output)
				}
			}
		})
	}
}

func TestBlackHoledAddrBlocked(t *testing.T) {
	resolver, err := madns.NewResolver()
	if err != nil {
		t.Fatal(err)
	}
	s := newTestSwarmWithResolver(t, resolver)
	defer s.Close()

	n := 3
	s.bhd.ipv6 = &blackHoleFilter{n: n, minSuccesses: 1, name: "IPv6"}

	// all dials to the address will fail. RFC6666 Discard Prefix
	addr := ma.StringCast("/ip6/0100::1/tcp/54321/")

	p, err := test.RandPeerID()
	if err != nil {
		t.Error(err)
	}
	s.Peerstore().AddAddr(p, addr, peerstore.PermanentAddrTTL)

	// do 1 extra dial to ensure that the blackHoleDetector state is updated since it
	// happens in a different goroutine
	for i := 0; i < n+1; i++ {
		s.backf.Clear(p)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		conn, err := s.DialPeer(ctx, p)
		if err == nil || conn != nil {
			t.Fatalf("expected dial to fail")
		}
		cancel()
	}
	s.backf.Clear(p)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	conn, err := s.DialPeer(ctx, p)
	require.Nil(t, conn)
	var de *DialError
	if !errors.As(err, &de) {
		t.Fatalf("expected to receive an error of type *DialError, got %s of type %T", err, err)
	}
	require.ErrorIs(t, err, ErrDialRefusedBlackHole)
}
