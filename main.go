package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "embed"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/runtime/interop"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
)

// go:embed config.yaml
var embedFile string

func main() {

	clConfig := params.MinimalSpecConfig()

	if err := params.SetActive(clConfig); err != nil {
		panic(err)
	}

	genesisTime := uint64(time.Now().Add(time.Second).Unix())
	config := params.BeaconConfig()

	gen := interop.GethTestnetGenesis(genesisTime, config)
	// HACK: fix this in prysm?
	gen.Config.DepositContractAddress = gethcommon.HexToAddress(config.DepositContractAddress)

	priv, pub, err := interop.DeterministicallyGenerateKeys(0, 100)
	if err != nil {
		panic(err)
	}

	depositData, roots, err := interop.DepositDataFromKeysWithExecCreds(priv, pub, 100)
	if err != nil {
		panic(err)
	}

	opts := make([]interop.PremineGenesisOpt, 0)
	opts = append(opts, interop.WithDepositData(depositData, roots))

	block := gen.ToBlock()
	log.Printf("Genesis block hash: %s", block.Hash())

	v := version.Fulu

	state, err := interop.NewPreminedGenesis(context.Background(), genesisTime, 0, 100, v, block, opts...)
	if err != nil {
		panic(err)
	}
	fmt.Println(state)
}
