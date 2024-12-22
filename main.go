package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guddu75/bttorrent/torrentfile"
)

func main() {

	inpath := os.Args[1]
	outpath := os.Args[2]

	tf,err := torrentfile.Open(inpath)

	if(err != nil){
		log.Fatal(err)
	}

	err = 

}
