module github.com/zbliujia/go-libp2p

go 1.20

retract v0.26.1 // Tag was applied incorrectly due to a bug in the release workflow.

require (
	github.com/benbjohnson/clock v1.3.5
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0
	github.com/flynn/noise v1.0.0
	github.com/google/gopacket v1.1.19
	github.com/gorilla/websocket v1.5.0
	github.com/hashicorp/golang-lru/arc/v2 v2.0.5
	github.com/hashicorp/golang-lru/v2 v2.0.5
	github.com/ipfs/go-cid v0.4.1
	github.com/ipfs/go-datastore v0.6.0
	github.com/ipfs/go-ds-badger v0.3.0
	github.com/ipfs/go-ds-leveldb v0.5.0
	github.com/ipfs/go-log/v2 v2.5.1
	github.com/jbenet/go-temp-err-catcher v0.1.0
	github.com/klauspost/compress v1.16.7
	github.com/libp2p/go-buffer-pool v0.1.0
	github.com/libp2p/go-flow-metrics v0.1.0
	github.com/libp2p/go-libp2p-asn-util v0.3.0
	github.com/libp2p/go-libp2p-testing v0.12.0
	github.com/libp2p/go-msgio v0.3.0
	github.com/libp2p/go-nat v0.2.0
	github.com/libp2p/go-netroute v0.2.1
	github.com/libp2p/go-reuseport v0.4.0
	github.com/libp2p/go-yamux/v4 v4.0.1
	github.com/libp2p/zeroconf/v2 v2.2.0
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b
	github.com/minio/sha256-simd v1.0.1
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-base32 v0.1.0
	github.com/multiformats/go-multiaddr v0.11.0
	github.com/multiformats/go-multiaddr-dns v0.3.1
	github.com/multiformats/go-multiaddr-fmt v0.1.0
	github.com/multiformats/go-multibase v0.2.0
	github.com/multiformats/go-multicodec v0.9.0
	github.com/multiformats/go-multihash v0.2.3
	github.com/multiformats/go-multistream v0.5.0
	github.com/multiformats/go-varint v0.0.7
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pion/datachannel v1.5.5
	github.com/pion/ice/v2 v2.3.6
	github.com/pion/logging v0.2.2
	github.com/pion/stun v0.6.0
	github.com/pion/webrtc/v3 v3.2.9
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_model v0.4.0
	github.com/quic-go/quic-go v0.38.1
	github.com/quic-go/webtransport-go v0.5.3
	github.com/stretchr/testify v1.8.4
	github.com/wlynxg/anet v0.0.0-20230919012416-aa9594d8ec61
	go.uber.org/fx v1.20.0
	go.uber.org/goleak v1.2.0
	go.uber.org/mock v0.3.0
	go.uber.org/zap v1.25.0
	golang.org/x/crypto v0.12.0
	golang.org/x/exp v0.0.0-20230817173708-d852ddb80c63
	golang.org/x/sync v0.3.0
	golang.org/x/sys v0.11.0
	golang.org/x/tools v0.12.1-0.20230815132531-74c255bcf846
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/pprof v0.0.0-20230817174616-7a8ec2ada47b // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/huin/goupnp v1.2.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/koron/go-ssdp v0.0.4 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.55 // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/onsi/ginkgo/v2 v2.11.0 // indirect
	github.com/pion/dtls/v2 v2.2.7 // indirect
	github.com/pion/interceptor v0.1.17 // indirect
	github.com/pion/mdns v0.0.7 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.10 // indirect
	github.com/pion/rtp v1.7.13 // indirect
	github.com/pion/sctp v1.8.7 // indirect
	github.com/pion/sdp/v3 v3.0.6 // indirect
	github.com/pion/srtp/v2 v2.0.15 // indirect
	github.com/pion/transport/v2 v2.2.1 // indirect
	github.com/pion/turn/v2 v2.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-20 v0.3.3 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
)
