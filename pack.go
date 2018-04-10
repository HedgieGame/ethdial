package ethdial

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

func EnPack(numb []*big.Int, span []int) (*big.Int, error) {
	// pack big.Ints into a 256-bit pack
	// []numb are the collection of big.ints
	// []span are the number of bytes each big.int occupies
	fmt.Println("toenpack", numb)
	var err error
	pack := new(big.Int)
	if len(numb) == len(span) {
		for i := range numb {
			modu := 1 << uint(8*span[i])
			if numb[i].Cmp(big.NewInt(int64(modu))) == 1 {
				err = errors.New(numb[i].String() + " is larger than 2^" + strconv.Itoa(8*span[i]) + " (" + strconv.Itoa(modu) + ")")
				break
			}
			pack.Mul(pack, big.NewInt(int64(modu)))
			pack.Add(pack, big.NewInt(numb[i].Int64()))
		}
	} else {
		err = errors.New("lengths do not match")
	}
	return pack, err
}

func UnPack(pack *big.Int, span []int) ([]*big.Int, error) {
	// unpack big.Ints from a 256-bit pack
	// []span are the number of bytes each int occupies
	var err error
	var numb []*big.Int
	if len(span) > 0 {
		for i := range span {
			remn := new(big.Int).Set(pack)
			modu := big.NewInt(int64(1 << uint(8*span[len(span)-1-i])))
			remn.Mod(remn, modu)
			numb = append([]*big.Int{remn}, numb...)
			pack.Div(pack, modu)
		}
	} else {
		err = errors.New("length of spans is nil")
	}
	fmt.Println("unpacked", numb)
	return numb, err
}
