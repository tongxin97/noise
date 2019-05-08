package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/cynthiatong/noise"
	"github.com/cynthiatong/noise/log"
	"github.com/cynthiatong/noise/network"
	"github.com/cynthiatong/noise/protocol"
	kad "github.com/cynthiatong/noise/skademlia"
)

var (
	ip       = "127.0.0.1"
	bsAddr   = []string{"127.0.0.1:7000"}
	node     *noise.Node
	numPeers = 50
	peerFile = "peers.txt"
)

func randID(ids []kad.ID) kad.ID {
	n := rand.Int() % len(ids)
	return ids[n]
}

func main() {
	portFlag := flag.Uint("p", 7000, "")
	flag.Parse()
	rand.Seed(time.Now().Unix())
	reader := bufio.NewReader(os.Stdin)

	if *portFlag == 4000 {
		peers := kad.LoadIDs(peerFile)
		bsAddrs := make([]kad.ID, 10)
		for i := 0; i < 10; i++ {
			bsAddrs[i] = randID(peers)
		}
		node = network.InitNetworkNode(ip, *portFlag, kad.IDAddresses(bsAddrs))
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			switch input {
			case "fn\n":
				target := randID(peers)
				found := kad.FindNode(node, target, kad.BucketSize(), 8)
				log.Info().Msgf("closest peers to target %s: %+v", target.Address(), kad.IDAddresses(found))
			}
		}
	} else {
		peers := []kad.ID{}
		for i := 0; i < numPeers; i++ {
			node = network.InitNetworkNode(ip, *portFlag, kad.IDAddresses(peers))
			*portFlag++
			peers = append(peers, protocol.NodeID(node).(kad.ID))
		}
		kad.PersistIDs(peerFile, peers)
		select {}
	}
}
