package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"bitbucket.org/billyharvey/ethdial"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	hed := &ethdial.Hedgie{HID: 12345}

	//peek hedgie
	hed, err := hed.Peek()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//	spew.Dump(hed)

	// change some values
	hed.Air = float64(rand.Intn(100000))
	hed.Charm = float64(rand.Intn(100000))
	hed.Earth = float64(rand.Intn(100000))
	hed.Fire = float64(rand.Intn(100000))
	hed.Intelligence = float64(rand.Intn(100000))
	hed.Luck = float64(rand.Intn(100000))
	hed.Prudence = float64(rand.Intn(100000))
	hed.Water = float64(rand.Intn(100000))
	hed.Level = rand.Intn(256)
	hed.Status = rand.Intn(256)
	hed.Tier = rand.Intn(256)
	//	spew.Dump(hed)

	// poke hedgie
	tid, err := hed.Poke(100, done)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tid)
	for {
	}
}

func done(e *ethdial.EthDial) {
	fmt.Println(e.EthTranString, e.EthTranStatus)
	os.Exit(0)
}
