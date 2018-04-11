package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"bitbucket.org/billyharvey/ethdial"
	hedgie "github.com/HedgieGame/hedgie-server/app/domain"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	//peek hedgie
	hed, err := ethdial.Peek(12345)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//	spew.Dump(hed)

	hStatus := []hedgie.HedgieStatus{hedgie.StatusAvail, hedgie.StatusPending, hedgie.StatusSold}
	hTier := []hedgie.HedgieTierLevel{hedgie.HedgieTier1, hedgie.HedgieTier2, hedgie.HedgieTier3, hedgie.HedgieTier4, hedgie.HedgieTier5, hedgie.HedgieTier6, hedgie.HedgieTier7}
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
	hed.Status = hStatus[rand.Intn(3)]
	hed.Tier = hTier[rand.Intn(7)]
	//	spew.Dump(hed)

	// poke hedgie
	tid, err := ethdial.Poke(hed, 100, done)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tid, "pending")
	for {
	}
}

func done(e *ethdial.EthDial) {
	fmt.Println(e.EthTranString, e.EthTranStatus)
	os.Exit(0)
}
