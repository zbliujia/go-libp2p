package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	ll "github.com/ipfs/go-log/v2"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/zbliujia/go-libp2p/core/host"
	"io"
	"log"
	"net/http"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/zbliujia/go-libp2p"
	"github.com/zbliujia/go-libp2p/core/network"
	"github.com/zbliujia/go-libp2p/core/peer"
)

func main() {
	ll.SetAllLoggers(ll.LevelDebug)
	run()
}

func run() {
	// Create two "unreachable" libp2p hosts that want to communicate.
	// We are configuring them with no listen addresses to mimic hosts
	// that cannot be directly dialed due to problematic firewall/NAT
	// configurations.
	relay := flag.String("relay", "", "relay addrs")
	d := flag.String("d", "", "d id")
	o := flag.String("o", "", "other id")
	flag.Parse()
	println(o)
	if *relay == "" {
		fmt.Println("Please Use -relay ")
		return
	}

	if *d == "" {
		fmt.Println("Please Use -d ")
		return
	}

	unreachable1, err := libp2p.New(
		//libp2p.NoListenAddrs,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10000", "/ip4/0.0.0.0/udp/10001/quic-v1"),
		// Usually EnableRelay() is not required as it is enabled by default
		// but NoListenAddrs overrides this, so we're adding it in explictly again.
		libp2p.EnableRelay(),
		libp2p.EnableHolePunching(),
	)
	if err != nil {
		log.Printf("Failed to create unreachable1: %v", err)
		return
	}

	log.Println("First let's attempt to directly connect")

	unreachable1info := peer.AddrInfo{
		ID:    unreachable1.ID(),
		Addrs: unreachable1.Addrs(),
	}

	log.Printf("unreachable1info=%s", unreachable1info.String())

	relay1info, err := peer.AddrInfoFromString(*relay)
	if err != nil {
		log.Printf("Failed to create AddrInfoFromString: %v", err)
		return
	}
	log.Printf("relay1info=%s", relay1info.String())

	// Connect both unreachable1 and unreachable2 to relay1
	if err := unreachable1.Connect(context.Background(), *relay1info); err != nil {
		log.Printf("Failed to connect unreachable1 and relay1: %v", err)
		return
	}

	// Now create a new address for unreachable2 that specifies to communicate via
	// relay1 using a circuit relay
	//relayaddr, err := ma.NewMultiaddr(*relay + "/p2p-circuit/p2p/" + *d)
	//relayaddr, err := ma.NewMultiaddr("/p2p/" + relay1info.ID.String() + "/p2p-circuit/p2p/" + *d)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//id, err := peer.Decode(*d)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//// Open a connection to the previously unreachable host via the relay address
	//unreachable2relayinfo := peer.AddrInfo{
	//	ID:    id,
	//	Addrs: []ma.Multiaddr{relayaddr},
	//}
	//
	//log.Printf("unreachable2relayinfo=%s", unreachable2relayinfo.String())
	//if err := unreachable1.Connect(context.Background(), unreachable2relayinfo); err != nil {
	//	log.Printf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
	//	return
	//}

	log.Println("Yep, that worked!")
	//sendTestMessage(relay1info, o, unreachable1)
	unreachable2relayinfo, err := sendTestMessage(relay1info, d, unreachable1)
	if err != nil {
		log.Fatalln(err)
		return
	}

	proxyAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 9900))
	if err != nil {
		log.Fatalln(err)
		return
	}
	// Woohoo! we're connected!
	// Let's start talking!
	service := &ProxyService{
		relayAddr: *relay1info,
		host:      unreachable1,
		destAddr:  *unreachable2relayinfo,
		proxyAddr: proxyAddr,
	}

	print(peer.AddrInfo{
		ID:    unreachable1.ID(),
		Addrs: unreachable1.Addrs(),
	}.String())

	service.Serve()

}

func sendTestMessage(relay1info *peer.AddrInfo, d *string, unreachable1 host.Host) (*peer.AddrInfo, error) {
	relayaddr, err := ma.NewMultiaddr("/p2p/" + relay1info.ID.String() + "/p2p-circuit/p2p/" + *d)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	id, err := peer.Decode(*d)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Open a connection to the previously unreachable host via the relay address
	unreachable2relayinfo := peer.AddrInfo{
		ID:    id,
		Addrs: []ma.Multiaddr{relayaddr},
	}

	log.Printf("unreachable2relayinfo=%s", unreachable2relayinfo.String())
	if err := unreachable1.Connect(context.Background(), unreachable2relayinfo); err != nil {
		log.Printf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
		return nil, err
	}

	s, err := unreachable1.NewStream(network.WithUseTransient(context.Background(), "customprotocol"), unreachable2relayinfo.ID, "/customprotocol")
	if err != nil {
		log.Println("Whoops, this should have worked...: ", err)
		return nil, err
	}
	response := make([]byte, 102400)
	n, err := s.Read(response) // block until the handler closes the stream
	if err != nil {
		log.Println("read...: ", err, n)
	} else {
		log.Println("read...: ", string(response), n)
	}
	return &unreachable2relayinfo, err
}

type ProxyService struct {
	relayAddr peer.AddrInfo
	destAddr  peer.AddrInfo
	host      host.Host
	proxyAddr ma.Multiaddr
}

// Serve listens on the ProxyService's proxy address. This effectively
// allows to set the listening address as http proxy.
func (p *ProxyService) Serve() {
	_, serveArgs, _ := manet.DialArgs(p.proxyAddr)
	fmt.Println("proxy listening on ", serveArgs)
	if p.destAddr.ID != "" {
		http.ListenAndServe(serveArgs, p)
	}
}

// ServeHTTP implements the http.Handler interface. WARNING: This is the
// simplest approach to a proxy. Therefore, we do not do any of the things
// that should be done when implementing a reverse proxy (like handling
// headers correctly). For how to do it properly, see:
// https://golang.org/src/net/http/httputil/reverseproxy.go?s=3845:3920#L121
//
// ServeHTTP opens a stream to the dest peer for every HTTP request.
// Streams are multiplexed over single connections so, unlike connections
// themselves, they are cheap to create and dispose of.
func (p *ProxyService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("proxying request for %s to peer %s\n", r.URL, p.destAddr.ID)

	if p.host.Network().Connectedness(p.relayAddr.ID) == network.Connected {

	}

	if p.host.Network().Connectedness(p.destAddr.ID) == network.Connected {

	}

	if err := p.host.Connect(context.Background(), p.relayAddr); err != nil {
		log.Printf("Unexpected error here. Failed to connect unreachable1 and relay: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := p.host.Connect(context.Background(), p.destAddr); err != nil {
		log.Printf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// We need to send the request to the remote libp2p peer, so
	// we open a stream to it
	stream, err := p.host.NewStream(network.WithUseTransient(context.Background(), "proxy-example"), p.destAddr.ID, "/proxy-example/0.0.1")
	// If an error happens, we write an error for response.
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	// r.Write() writes the HTTP request to the stream.
	err = r.Write(stream)
	if err != nil {
		stream.Reset()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Now we read the response that was sent from the dest
	// peer
	buf := bufio.NewReader(stream)
	resp, err := http.ReadResponse(buf, r)
	if err != nil {
		stream.Reset()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	// Copy any headers
	for k, v := range resp.Header {
		for _, s := range v {
			w.Header().Add(k, s)
		}
	}

	// Write response status and headers
	w.WriteHeader(resp.StatusCode)

	// Finally copy the body
	io.Copy(w, resp.Body)
	resp.Body.Close()
}
