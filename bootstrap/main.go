package main

import (
	"log"
	"os"

	"github.com/amirhnajafiz/telescope/bootstrap/internal/config"
	"github.com/amirhnajafiz/telescope/bootstrap/internal/ipfs"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		panic(err)
	}

	// create a new IPFS client
	ipfsC := ipfs.NewClient(cfg.IPFSGateway)

	// read all dir names inside the bootstrap input directory
	dirs, err := os.ReadDir(cfg.IDP)
	if err != nil {
		panic(err)
	}

	// iterate over directories and upload them to IPFS
	for _, dir := range dirs {
		if cid, err := ipfsC.PutDIR(dir.Name()); err != nil {
			log.Printf("failed to upload %s: %v", dir.Name(), err)
		} else {
			log.Printf("uploaded %s to IPFS with CID: %s", dir.Name(), cid)
		}
	}
}
