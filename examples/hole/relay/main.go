package main

import (
	"flag"
	"fmt"
	ll "github.com/ipfs/go-log/v2"
	"log"
	"time"

	"github.com/zbliujia/go-libp2p"
	"github.com/zbliujia/go-libp2p/p2p/protocol/circuitv2/relay"

	"github.com/zbliujia/go-libp2p/core/peer"
)

func main() {
	ll.SetAllLoggers(ll.LevelDebug)
	run()
}

func run() {
	// Create a host to act as a middleman to relay messages on our behalf
	port := flag.Int("p", 12000, "port")
	flag.Parse()
	relay1, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *port-1), fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", *port)))
	if err != nil {
		log.Printf("Failed to create relay1: %v", err)
		return
	}

	// Configure the host to offer the circuit relay service.
	// Any host that is directly dialable in the network (or on the internet)
	// can offer a circuit relay service, this isn't just the job of
	// "dedicated" relay services.
	// In circuit relay v2 (which we're using here!) it is rate limited so that
	// any node can offer this service safely
	_, err = relay.New(relay1)
	if err != nil {
		log.Printf("Failed to instantiate the relay: %v", err)
		return
	}

	relay1info := peer.AddrInfo{
		ID:    relay1.ID(),
		Addrs: relay1.Addrs(),
	}

	addrs, err := peer.AddrInfoToP2pAddrs(&relay1info)
	if err != nil {
		log.Printf("Failed to AddrInfoToP2pAddrs: %v", err)
		return
	}
	fmt.Println("relay1info:", addrs)

	for {
		log.Printf("sleep")
		time.Sleep(time.Minute)
	}

}
