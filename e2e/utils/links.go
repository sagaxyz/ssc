package utils

import (
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// RelayerLinks stores the mapping of the different chain ids that are being
// set up for the test network and their respective connection links.
type RelayerLinks map[uint8]map[uint8]interchaintest.InterchainLink

// insert adds a new link to the map of relayer links and avoids nil access by initializing empty maps.
func (rl RelayerLinks) insert(idxA, idxB uint8, link interchaintest.InterchainLink) error {
	// prepare insertion into nested map and instantiate the sub-map if it's nil
	linksForA, found := rl[idxA]
	if found {
		if _, found := linksForA[idxB]; found {
			panic(fmt.Sprintf("overwriting link for path %d-%d", idxA, idxB))
		}
	} else {
		rl[idxA] = make(map[uint8]interchaintest.InterchainLink)
	}

	rl[idxA][idxB] = link

	return nil
}

func (icn *InterchainNetwork) GetLinks() RelayerLinks {
	return icn.links
}

// GetLink retrieves the link between chains with the indices A and B.
//
// CONTRACT: A < B
func (icn *InterchainNetwork) GetLink(iChainA, iChainB uint8) interchaintest.InterchainLink {
	if iChainA >= iChainB {
		panic(fmt.Sprintf("expected A < B; got A=%d, B=%d", iChainA, iChainB))
	}

	linksForA, found := icn.links[iChainA]
	if !found {
		panic(fmt.Sprintf("first chain index not found in links map: %d", iChainA))
	}

	link, found := linksForA[iChainB]
	if !found {
		panic(fmt.Sprintf("second chain index not found in links map: %d", iChainB))
	}

	return link
}

// getLinksToCreate either returns the configured values from the given configuration and validates those,
// or puts together the list of links that would connect all available connections.
func getLinksToCreate(config networkConfig) []RelayerPath {
	var linksToCreate []RelayerPath
	if len(config.relayerPaths) == 0 {
		for a := range config.nChains {
			for b := a + 1; b < config.nChains; b++ {
				if a == b {
					continue
				}

				linksToCreate = append(linksToCreate, [2]uint8{a, b})
			}
		}
	} else {
		linksToCreate = config.relayerPaths
	}

	return linksToCreate
}

func (icn *InterchainNetwork) createAndAddRelayers(
	t *testing.T,
	chains []ibc.Chain,
) error {
	if icn.links != nil {
		panic("relayer links already set")
	}

	linksToCreate := getLinksToCreate(icn.config)
	links := make(RelayerLinks)
	for _, link := range linksToCreate {
		// Create a new relayer instance for each connection
		relayerFactory := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t, zaptest.Level(zap.ErrorLevel)))
		relayer := relayerFactory.Build(t, icn.client, icn.networkName)

		l := interchaintest.InterchainLink{
			Chain1:  chains[link[0]],
			Chain2:  chains[link[1]],
			Relayer: relayer,
			Path:    fmt.Sprintf("link-path-%d-%d", link[0], link[1]),
		}

		icn.interchain = icn.interchain.
			AddRelayer(relayer, fmt.Sprintf("relayer%d-%d", link[0], link[1])).
			AddLink(l)

		err := links.insert(link[0], link[1], l)
		if err != nil {
			return err
		}
	}

	icn.links = links
	return nil
}
