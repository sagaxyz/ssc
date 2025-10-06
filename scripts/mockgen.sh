#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -source=x/chainlet/types/expected_keepers.go -package testutil -destination x/chainlet/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/billing/types/expected_keepers.go -package testutil -destination x/billing/testutil/expected_keepers_mocks.go
#$mockgen_cmd -source=x/epochs/types/expected_keepers.go -package testutil -destination x/epoch/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/escrow/types/expected_keepers.go -package testutil -destination x/escrow/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/liquid/types/expected_keepers.go -package testutil -destination x/liquid/testutil/expected_keepers_mocks.go
