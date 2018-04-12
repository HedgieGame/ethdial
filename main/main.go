package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"

	"bitbucket.org/billyharvey/ethdial"
	hedgie "github.com/HedgieGame/hedgie-server/app/domain"
	"github.com/davecgh/go-spew/spew"
)

var Config ethdial.Config

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	Config.Endpoint = "https://rinkeby.infura.io/pQZitksokILr3E3rp7u8"
	Config.Private = "12BF6F0806822A6763205D012A3302F73646B50DA9F4B71826CD86F794EE5B3E"
	Config.Contract = "0xBa4764def35E38397Fbdd7e6570a9Da97378a5c3"
	Config.GasLimit = uint64(350000)
	Config.GasPrice = big.NewInt(15000000000)
	Config.PeekFunc = "Peek(uint256)"
	Config.PokeFunc = "Poke(uint256,uint256)"

	//peek hedgie
	hed, err := ethdial.Peek(&Config, 12345)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(hed)

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
	spew.Dump(hed)

	// poke hedgie
	tid, err := ethdial.Poke(&Config, hed, 100, done)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tid, "pending")
	select {}
}

func done(e *ethdial.EthDial) {
	fmt.Println(e.EthTranString, e.EthTranStatus)
	//peek hedgie
	hed, err := ethdial.Peek(&Config, 12345)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	spew.Dump(hed)
	os.Exit(0)
}
