package libp2phttp_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/zbliujia/go-libp2p"
	"github.com/zbliujia/go-libp2p/core/peer"
	libp2phttp "github.com/zbliujia/go-libp2p/p2p/http"
)

func ExampleHost_withAStockGoHTTPClient() {
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/0/http")},
	}

	// A server with a simple echo protocol
	server.SetHTTPHandler("/echo/1.0.0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/octet-stream")
		io.Copy(w, r.Body)
	}))
	go server.Serve()
	defer server.Close()

	var serverHTTPPort string
	var err error
	for _, a := range server.Addrs() {
		serverHTTPPort, err = a.ValueForProtocol(ma.P_TCP)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	// Make an HTTP request using the Go standard library.
	resp, err := http.Post("http://127.0.0.1:"+serverHTTPPort+"/echo/1.0.0/", "application/octet-stream", strings.NewReader("Hello HTTP"))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	// Output: Hello HTTP
}

func ExampleHost_listenOnHTTPTransportAndStreams() {
	serverStreamHost, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/udp/50124/quic-v1"))
	if err != nil {
		log.Fatal(err)
	}
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50124/http")},
		StreamHost:        serverStreamHost,
	}
	go server.Serve()
	defer server.Close()

	fmt.Println("Server listening on:", server.Addrs())
	// Output: Server listening on: [/ip4/127.0.0.1/udp/50124/quic-v1 /ip4/127.0.0.1/tcp/50124/http]
}

func ExampleHost_overLibp2pStreams() {
	serverStreamHost, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/udp/0/quic-v1"))
	if err != nil {
		log.Fatal(err)
	}

	server := libp2phttp.Host{
		StreamHost: serverStreamHost,
	}

	// A server with a simple echo protocol
	server.SetHTTPHandler("/echo/1.0.0", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/octet-stream")
		io.Copy(w, r.Body)
	}))
	go server.Serve()
	defer server.Close()

	clientStreamHost, err := libp2p.New(libp2p.NoListenAddrs)
	if err != nil {
		log.Fatal(err)
	}

	client := libp2phttp.Host{StreamHost: clientStreamHost}

	// Make an HTTP request using the Go standard library, but over libp2p
	// streams. If the server were listening on an HTTP transport, this could
	// also make the request over the HTTP transport.
	httpClient, err := client.NamespacedClient("/echo/1.0.0", peer.AddrInfo{ID: server.PeerID(), Addrs: server.Addrs()})

	// Only need to Post to "/" because this client is namespaced to the "/echo/1.0.0" protocol.
	resp, err := httpClient.Post("/", "application/octet-stream", strings.NewReader("Hello HTTP"))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	// Output: Hello HTTP
}

func ExampleHost_Serve() {
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50221/http")},
	}

	go server.Serve()
	defer server.Close()

	fmt.Println(server.Addrs())

	// Output: [/ip4/127.0.0.1/tcp/50221/http]
}

func ExampleHost_SetHTTPHandler() {
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50222/http")},
	}

	server.SetHTTPHandler("/hello/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Hello World"))
	}))

	go server.Serve()
	defer server.Close()

	port, err := server.Addrs()[0].ValueForProtocol(ma.P_TCP)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get("http://127.0.0.1:" + port + "/hello/1/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))

	// Output: Hello World
}

func ExampleHost_SetHTTPHandlerAtPath() {
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50224/http")},
	}

	server.SetHTTPHandlerAtPath("/hello/1", "/other-place/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Hello World"))
	}))

	go server.Serve()
	defer server.Close()

	port, err := server.Addrs()[0].ValueForProtocol(ma.P_TCP)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get("http://127.0.0.1:" + port + "/other-place/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))

	// Output: Hello World
}

func ExampleHost_NamespacedClient() {
	var client libp2phttp.Host

	// Create the server
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50221/http")},
	}

	server.SetHTTPHandlerAtPath("/hello/1", "/other-place/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Hello World"))
	}))

	go server.Serve()
	defer server.Close()

	// Create an http.Client that is namespaced to this protocol.
	httpClient, err := client.NamespacedClient("/hello/1", peer.AddrInfo{ID: server.PeerID(), Addrs: server.Addrs()})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := httpClient.Get("/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))

	// Output: Hello World
}

func ExampleHost_NamespaceRoundTripper() {
	var client libp2phttp.Host

	// Create the server
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50223/http")},
	}

	server.SetHTTPHandler("/hello/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Hello World"))
	}))

	go server.Serve()
	defer server.Close()

	// Create an http.Roundtripper for the server
	rt, err := client.NewConstrainedRoundTripper(peer.AddrInfo{ID: server.PeerID(), Addrs: server.Addrs()})
	if err != nil {
		log.Fatal(err)
	}

	// Namespace this roundtripper to a specific protocol
	rt, err = client.NamespaceRoundTripper(rt, "/hello/1", server.PeerID())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := (&http.Client{Transport: rt}).Get("/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))

	// Output: Hello World
}

func ExampleHost_NewConstrainedRoundTripper() {
	var client libp2phttp.Host

	// Create the server
	server := libp2phttp.Host{
		InsecureAllowHTTP: true, // For our example, we'll allow insecure HTTP
		ListenAddrs:       []ma.Multiaddr{ma.StringCast("/ip4/127.0.0.1/tcp/50225/http")},
	}

	server.SetHTTPHandler("/hello/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Hello World"))
	}))

	go server.Serve()
	defer server.Close()

	// Create an http.Roundtripper for the server
	rt, err := client.NewConstrainedRoundTripper(peer.AddrInfo{ID: server.PeerID(), Addrs: server.Addrs()})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := (&http.Client{Transport: rt}).Get("/hello/1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))

	// Output: Hello World
}

func ExampleWellKnownHandler() {
	var h libp2phttp.WellKnownHandler
	h.AddProtocolMeta("/hello/1", libp2phttp.ProtocolMeta{
		Path: "/hello-path/",
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	// Serve `.well-known/libp2p`. Note, this is handled automatically if you use the libp2phttp.Host.
	go http.Serve(listener, &h)

	// Get the `.well-known/libp2p` resource
	resp, err := http.Get("http://" + listener.Addr().String() + "/.well-known/libp2p")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody))
	// Output: {"/hello/1":{"path":"/hello-path/"}}

}
