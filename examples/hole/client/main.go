package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/zbliujia/go-libp2p"
	"github.com/zbliujia/go-libp2p/core/network"
	"github.com/zbliujia/go-libp2p/core/peer"
)

func main() {
	run()
}

func run() {
	// Create two "unreachable" libp2p hosts that want to communicate.
	// We are configuring them with no listen addresses to mimic hosts
	// that cannot be directly dialed due to problematic firewall/NAT
	// configurations.
	relay := flag.String("relay", "", "relay addrs")
	d := flag.String("d", "", "d id")
	flag.Parse()

	if *relay == "" {
		fmt.Println("Please Use -relay ")
		return
	}

	if *d == "" {
		fmt.Println("Please Use -d ")
		return
	}

	unreachable1, err := libp2p.New(
		libp2p.NoListenAddrs,
		// Usually EnableRelay() is not required as it is enabled by default
		// but NoListenAddrs overrides this, so we're adding it in explictly again.
		libp2p.EnableRelay(),
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
	relayaddr, err := ma.NewMultiaddr("/p2p/" + relay1info.ID.String() + "/p2p-circuit/p2p/" + *d)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := peer.Decode(*d)
	if err != nil {
		log.Println(err)
		return
	}
	// Open a connection to the previously unreachable host via the relay address
	unreachable2relayinfo := peer.AddrInfo{
		ID:    id,
		Addrs: []ma.Multiaddr{relayaddr},
	}

	log.Printf("unreachable2relayinfo=%s", unreachable2relayinfo.String())
	if err := unreachable1.Connect(context.Background(), unreachable2relayinfo); err != nil {
		log.Printf("Unexpected error here. Failed to connect unreachable1 and unreachable2: %v", err)
		return
	}

	log.Println("Yep, that worked!")

	// Woohoo! we're connected!
	// Let's start talking!

	// Because we don't have a direct connection to the destination node - we have a relayed connection -
	// the connection is marked as transient. Since the relay limits the amount of data that can be
	// exchanged over the relayed connection, the application needs to explicitly opt-in into using a
	// relayed connection. In general, we should only do this if we have low bandwidth requirements,
	// and we're happy for the connection to be killed when the relayed connection is replaced with a
	// direct (holepunched) connection.
	s, err := unreachable1.NewStream(network.WithUseTransient(context.Background(), "customprotocol"), unreachable2relayinfo.ID, "/customprotocol")
	if err != nil {
		log.Println("Whoops, this should have worked...: ", err)
		return
	}
	response := make([]byte, 102400)
	n, err := s.Read(response) // block until the handler closes the stream
	if err != nil {
		log.Println("read...: ", err, n)
	} else {
		log.Println("read...: ", string(response), n)
	}
}
