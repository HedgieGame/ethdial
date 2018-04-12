package ethdial

// TODO create a config.json for the contract variables
// TODO which perhaps should exist above this code in the stack(?)
// TODO in a single consistent place for configuration

import (
	"math/big"

	hedgie "github.com/HedgieGame/hedgie-server/app/domain"
)

// the HID is unique, but so is the color - one can serve as the other
// and neither needs to be stored - it's the lookup key itself
// the owner doesn't need and probably shouldn't be stored with the Hedgie
// store a mapping from owner <-> color
// store a mapping from name  <-> color
// If we allow small and cap letters and digits and space and dash that is 64 characters
// which can be stored in 6 bits, allowing a 42 character name to be stored in 256 bits
// 7 bits allows 37 character names
// 8 bits allow 32 character names
// if we allow UTF-8 it will be variable but the ASCII set is each 8 bits

// these change if the contract changes
//var endpoint = "http://localhost:8545"
//var private = "12bf6f0806822a6763205d012a3302f73646b50da9f4b71826cd86f794ee5b3e"
//var contract = "0x014194F3D48c61bF768e70e2AD39c2d80c66f6ce"
//var gasLimit = uint64(1000000)
//var gasPrice = big.NewInt(1 * 1000000000) // 1 gwei
//var peekFunc = "Peek(uint256)"
//var pokeFunc = "Poke(uint256,uint256)"
// var endpoint = "https://mainnet.infura.io/pQZitksokILr3E3rp7u8"

type Config struct {
	Endpoint string
	Private  string
	Contract string
	GasLimit uint64
	GasPrice *big.Int
	PeekFunc string
	PokeFunc string
	BossFunc string
}

// these are the number of bytes each struct element occupies in Eth storage
var sAir = 3
var sCharm = 3
var sEarth = 3
var sFire = 3
var sIntelligence = 3
var sLuck = 3
var sPrudence = 3
var sWater = 3
var sLevel = 1
var sStatus = 1
var sTier = 1

// this will need to be stored in a separate mapping since owner is an eth address
var sOwner = 32

// if we decide to store name it will also need to be a separate mapping
var sName = 32

func Boss(config *Config) (string, error) {
	// fetch the boss
	var boss string
	e := New().
		Addr(config.Contract).
		Call(config.BossFunc).
		Dial(config.Endpoint)
	if e.Error == nil {
		boss = e.Result.Text(16)
	}
	return boss, e.Error
}

func Peek(config *Config, hid int) (*hedgie.Hedgie, error) {
	// fetch the *Hedgie
	var err error
	hed := &hedgie.Hedgie{HID: hid}
	e := New().
		Addr(config.Contract).
		Call(config.PeekFunc).
		DataInt(hid).
		Dial(config.Endpoint)
	if e.Error == nil {
		b, err := UnPack(e.Result, []int{sAir, sCharm, sEarth, sFire, sIntelligence, sLuck, sPrudence, sWater, sLevel, sStatus, sTier})
		if err == nil {
			hed.Air, _ = new(big.Float).SetInt(b[0]).Float64()
			hed.Charm, _ = new(big.Float).SetInt(b[1]).Float64()
			hed.Earth, _ = new(big.Float).SetInt(b[2]).Float64()
			hed.Fire, _ = new(big.Float).SetInt(b[3]).Float64()
			hed.Intelligence, _ = new(big.Float).SetInt(b[4]).Float64()
			hed.Luck, _ = new(big.Float).SetInt(b[5]).Float64()
			hed.Prudence, _ = new(big.Float).SetInt(b[6]).Float64()
			hed.Water, _ = new(big.Float).SetInt(b[7]).Float64()
			level, _ := new(big.Float).SetInt(b[8]).Int64()
			hed.Level = int(level)
			status, _ := new(big.Float).SetInt(b[9]).Int64()
			switch status {
			case 0:
				hed.Status = hedgie.StatusAvail
			case 1:
				hed.Status = hedgie.StatusPending
			case 2:
				hed.Status = hedgie.StatusSold
			}
			tier, _ := new(big.Float).SetInt(b[10]).Int64()
			switch tier {
			case 1:
				hed.Tier = hedgie.HedgieTier1
			case 2:
				hed.Tier = hedgie.HedgieTier2
			case 3:
				hed.Tier = hedgie.HedgieTier3
			case 4:
				hed.Tier = hedgie.HedgieTier4
			case 5:
				hed.Tier = hedgie.HedgieTier5
			case 6:
				hed.Tier = hedgie.HedgieTier6
			case 7:
				hed.Tier = hedgie.HedgieTier7
			}
		}
	}
	return hed, err
}

func Poke(config *Config, hed *hedgie.Hedgie, wait int, fTran funcTran) (string, error) {
	// store the *Hedgie
	var err error
	var tid string
	bAir := new(big.Int)
	bCharm := new(big.Int)
	bEarth := new(big.Int)
	bFire := new(big.Int)
	bIntelligence := new(big.Int)
	bLuck := new(big.Int)
	bPrudence := new(big.Int)
	bWater := new(big.Int)
	big.NewFloat(hed.Air).Int(bAir)
	big.NewFloat(hed.Charm).Int(bCharm)
	big.NewFloat(hed.Earth).Int(bEarth)
	big.NewFloat(hed.Fire).Int(bFire)
	big.NewFloat(hed.Intelligence).Int(bIntelligence)
	big.NewFloat(hed.Luck).Int(bLuck)
	big.NewFloat(hed.Prudence).Int(bPrudence)
	big.NewFloat(hed.Water).Int(bWater)
	bLevel := big.NewInt(int64(hed.Level))
	bStatus := big.NewInt(int64(hed.Status))
	bTier := big.NewInt(int64(hed.Tier))
	packed, err := EnPack(
		[]*big.Int{bAir, bCharm, bEarth, bFire, bIntelligence, bLuck, bPrudence, bWater, bLevel, bStatus, bTier},
		[]int{sAir, sCharm, sEarth, sFire, sIntelligence, sLuck, sPrudence, sWater, sLevel, sStatus, sTier})
	if err == nil {
		e := New().
			Sign(config.Private).
			Addr(config.Contract).
			Limit(config.GasLimit).
			Price(config.GasPrice).
			Call(config.PokeFunc).
			DataInt(hed.HID).
			DataBigInt(packed).
			Tran(wait, fTran).
			Dial(config.Endpoint)
		tid = e.EthTranString
		err = e.Error
	}
	return tid, err
}
