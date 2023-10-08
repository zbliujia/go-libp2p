package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zbliujia/go-libp2p/p2p/protocol/circuitv2/client"
	"io"
	"log"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/zbliujia/go-libp2p"
	"github.com/zbliujia/go-libp2p/core/network"
	"github.com/zbliujia/go-libp2p/core/peer"
)

func main() {
	run()
}

func getPeerInfoByDest(relay string) (*peer.AddrInfo, error) {
	relayAddr, err := ma.NewMultiaddr(relay)
	if err != nil {
		return nil, err
	}
	pid, err := relayAddr.ValueForProtocol(ma.P_P2P)
	if err != nil {
		return nil, err
	}
	relayPeerID, err := peer.Decode(pid)
	if err != nil {
		return nil, err
	}

	relayPeerAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", pid))
	relayAddress := relayAddr.Decapsulate(relayPeerAddr)
	peerInfo := &peer.AddrInfo{
		ID:    relayPeerID,
		Addrs: []ma.Multiaddr{relayAddress},
	}
	return peerInfo, err
}

func run() {

	relay := flag.String("relay", "", "relay addrs")
	flag.Parse()

	if *relay == "" {
		fmt.Println("Please Use -relay ")
		return
	}

	unreachable2, err := libp2p.New(
		libp2p.NoListenAddrs,
		libp2p.EnableRelay(),
	)
	if err != nil {
		log.Printf("Failed to create unreachable2: %v", err)
		return
	}

	relay1info, err := peer.AddrInfoFromString(*relay)
	//relay1info, err = getPeerInfoByDest(*relay)
	if err != nil {
		log.Printf("Failed to create relay1info: %v", err)
		return
	}

	log.Printf("relay1info=%+v", relay1info)

	if err := unreachable2.Connect(context.Background(), *relay1info); err != nil {
		log.Printf("Failed to connect unreachable2 and relay1: %v", err)
		return
	}

	// Now, to test the communication, let's set up a protocol handler on unreachable2
	unreachable2.SetStreamHandler("/customprotocol", func(s network.Stream) {
		log.Println("Awesome! We're now communicating via the relay!")
		io.WriteString(s, "hello world")
		// End the example
		s.Close()
	})

	// Hosts that want to have messages relayed on their behalf need to reserve a slot
	// with the circuit relay service host
	// As we will open a stream to unreachable2, unreachable2 needs to make the
	// reservation
	// 这个东西 会断开的 需要保护
	_, err = client.Reserve(context.Background(), unreachable2, *relay1info)
	if err != nil {
		log.Printf("unreachable2 failed to receive a relay reservation from relay1. %v", err)
		return
	}

	// Now create a new address for unreachable2 that specifies to communicate via
	// relay1 using a circuit relay
	relayaddr, err := ma.NewMultiaddr("/p2p/" + relay1info.ID.String() + "/p2p-circuit/p2p/" + unreachable2.ID().String())
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Now let's attempt to connect the hosts via the relay node")

	// Open a connection to the previously unreachable host via the relay address
	unreachable2relayinfo := peer.AddrInfo{
		ID:    unreachable2.ID(),
		Addrs: []ma.Multiaddr{relayaddr},
	}
	log.Printf("unreachable2relayinfo=%+v", unreachable2relayinfo)

	log.Println("Yep, that worked!")

	for {
		log.Printf("sleep")
		time.Sleep(time.Minute)
	}

}
