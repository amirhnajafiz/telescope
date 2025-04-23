package main

import (
	"log"
	"os"
	"path"

	"github.com/amirhnajafiz/telescope/bootstrap/internal/config"
	"github.com/amirhnajafiz/telescope/bootstrap/internal/ipfs"
)

func main() {
	// load configs
	cfg, err := config.LoadConfigs()
	if err != nil {
		log.Fatalf("failed to load configs: %v", err)
	}

	// create a new IPFS client
	ipfsC := ipfs.NewClient(cfg.IPFSGateway)

	// read all dir names inside the bootstrap input directory
	dirs, err := os.ReadDir(cfg.IDP)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	// open the bootstrap data file
	file, err := os.OpenFile(cfg.DataPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// iterate over directories and upload them to IPFS
	for _, dir := range dirs {
		// build the path to the directory
		path := path.Join(cfg.IDP, dir.Name())

		if cid, err := ipfsC.PutDIR(path); err != nil {
			log.Printf("failed to upload %s: %v", dir.Name(), err)
		} else {
			log.Printf("uploaded %s to IPFS with CID: %s", dir.Name(), cid)

			// write the CID to the bootstrap data file
			if _, err := file.WriteString(dir.Name() + ": " + cid + "\n"); err != nil {
				log.Printf("failed to write to file: %v", err)
			}
		}
	}
}
