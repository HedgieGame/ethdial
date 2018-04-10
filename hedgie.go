package ethdial

import (
	"math/big"
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

type Hedgie struct {
	HID          int     `json:"hid"`   // id of the hedgie - same as the color?
	Color        string  `json:"color"` // 2^24 3x8-bit
	Name         *string `json:"name"`  //
	Fire         float64 `json:"fire"`
	Water        float64 `json:"water"`
	Earth        float64 `json:"earth"`
	Air          float64 `json:"air"`
	Intelligence float64 `json:"intelligence"`
	Luck         float64 `json:"luck"`
	Prudence     float64 `json:"prudence"`
	Charm        float64 `json:"charm"`
	Owner        *string `json:"owner"`
	Status       int     `json:"status"`
	Tier         int     `json:"tier"`
	Level        int     `json:"level"`
	ImgURL       *string `json:"imgURL"`
}

// these change if the contract changes
//var endpoint = "http://localhost:8545"
//var private = "12bf6f0806822a6763205d012a3302f73646b50da9f4b71826cd86f794ee5b3e"
//var contract = "0x014194F3D48c61bF768e70e2AD39c2d80c66f6ce"
//var gasLimit = uint64(1000000)
//var gasPrice = big.NewInt(1 * 1000000000) // 1 gwei
//var peekFunc = "Peek(uint256)"
//var pokeFunc = "Poke(uint256,uint256)"
// var endpoint = "https://mainnet.infura.io/pQZitksokILr3E3rp7u8"
var endpoint = "https://rinkeby.infura.io/pQZitksokILr3E3rp7u8"
var private = "12BF6F0806822A6763205D012A3302F73646B50DA9F4B71826CD86F794EE5B3E"
var contract = "0xBa4764def35E38397Fbdd7e6570a9Da97378a5c3"
var gasLimit = uint64(100000)
var gasPrice = big.NewInt(1 * 10000000000) // 10 gwei
var peekFunc = "Peek(uint256)"
var pokeFunc = "Poke(uint256,uint256)"

var sIndex = 1

// these are the number of bytes each struct element occupies in Eth storage
// index 0
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

func (hed *Hedgie) Peek() (*Hedgie, error) {
	// fetch the *Hedgie
	var err error
	e := New().
		Addr(contract).
		Call(peekFunc).
		DataInt(hed.HID).
		Dial(endpoint)
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
			hed.Status = int(status)
			tier, _ := new(big.Float).SetInt(b[10]).Int64()
			hed.Tier = int(tier)
		}
	}
	return hed, err
}

func (hed *Hedgie) Poke(wait int, fTran funcTran) (string, error) {
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
			Sign(private).
			Addr(contract).
			Limit(gasLimit).
			Price(gasPrice).
			Call(pokeFunc).
			DataInt(hed.HID).
			DataBigInt(packed).
			Tran(wait, fTran).
			Dial(endpoint)
		tid = e.EthTranString
		err = e.Error
	}
	return tid, err
}
