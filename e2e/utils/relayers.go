package utils

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
)

// startRelayers starts the networks relayers and waits for blocks to be produced.
func (icn *InterchainNetwork) startRelayers(ctx context.Context) error {
	if icn.links == nil {
		panic("no links set for network")
	}

	for _, linksForA := range icn.links {
		for _, link := range linksForA {
			err := link.Relayer.StartRelayer(ctx, icn.eRep, link.Path)
			if err != nil {
				return err
			}
		}
	}

	heighters := make([]testutil.ChainHeighter, len(icn.chains))
	for i, c := range icn.chains {
		heighters[i] = c
	}

	if err := testutil.WaitForBlocks(ctx, 5, heighters...); err != nil {
		return err
	}

	return nil
}

func (icn *InterchainNetwork) StopRelayers(ctx context.Context, eRep *testreporter.RelayerExecReporter) error {
	var err error
	for a, linksForA := range icn.links {
		for b, link := range linksForA {
			err = link.Relayer.StopRelayer(ctx, eRep)
			if err != nil {
				// NOTE: we join the errors instead of the early return here to try and stop all relayers that work.
				err = errors.Join(err, fmt.Errorf("error stopping relayer %d-%d: %w", a, b, err))
			}
		}
	}

	return err
}

// registerRelayerCleanup registers a cleanup routine to the testing instance, that stops all running relayers.
func (icn *InterchainNetwork) registerRelayerCleanup(t *testing.T, ctx context.Context) {
	t.Cleanup(
		func() {
			err := icn.StopRelayers(ctx, icn.eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)
}
