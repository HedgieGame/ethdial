package ethdial

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type FuncTran func(*EthDial)

type EthDial struct {
	EthAddr       common.Address
	EthCall       string
	EthCancel     context.CancelFunc
	EthChainId    *big.Int
	EthClient     *ethclient.Client
	EthCoin       *big.Int
	EthContext    context.Context
	EthData       string
	Error         error
	EthFrom       common.Address
	EthLimit      uint64
	EthPrice      *big.Int
	EthNonce      uint64
	Result        *big.Int
	EthSign       *ecdsa.PrivateKey
	EthTran       common.Hash
	EthTranFunc   FuncTran
	EthTranStatus string
	EthTranString string
	EthTranWait   int
	EthWait       int
	Transaction   string
}

// initialize struct with some defaults
func New() *EthDial {
	return &EthDial{
		EthCoin: big.NewInt(0),
		EthWait: 10}
}

// make a connection to the endpoint with a context timeout
func (e *EthDial) Dial(a string) *EthDial {
	e.EthClient, e.Error = ethclient.Dial(a)
	if e.Error == nil {
		e.EthContext, e.EthCancel = context.WithDeadline(context.Background(), time.Now().Add(time.Duration(e.EthWait)*time.Second))
		if e.EthChainId == nil {
			e.EthChainId, e.Error = e.EthClient.NetworkID(e.EthContext)
		}
		if e.Error == nil {
			if e.EthCall == "" {
				e = e.TxSend()
			} else {
				if e.EthSign == nil {
					e = e.TxCallFree()
				} else {
					e = e.TxCallPaid()
				}
			}
		}
		e.EthCancel()
	}
	if e.EthTranFunc != nil || e.EthTranWait > 0 {
		go e.TranWait()
	}
	return e
}

func (e *EthDial) TranWait() {
	e.EthTranStatus = "unknown"
	ctx := context.Background()
	e.EthTranWait = int(time.Now().Unix()) + e.EthTranWait
	for {
		time.Sleep(10 * time.Second)
		receipt, _ := e.EthClient.TransactionReceipt(ctx, e.EthTran)
		if receipt == nil {
			e.EthTranStatus = "pending"
		} else {
			switch receipt.Status {
			case 0:
				e.EthTranStatus = "failure"
			case 1:
				e.EthTranStatus = "success"
			default:
				e.EthTranStatus = "unknown"
			}
			break
		}
		if int(time.Now().Unix()) > e.EthTranWait {
			e.EthTranStatus = "timeout"
		}
	}
	if e.EthTranFunc != nil {
		e.EthTranFunc(e)
	}
}

func (e *EthDial) TxCallFree() *EthDial {
	msg := ethereum.CallMsg{
		From:     e.EthFrom,
		To:       &e.EthAddr,
		Gas:      e.EthLimit,
		GasPrice: e.EthPrice,
		Value:    e.EthCoin,
		Data:     common.Hex2Bytes(e.EthCall + e.EthData)}
	result, err := e.EthClient.CallContract(e.EthContext, msg, nil)
	e.Error = err
	if e.Error == nil {
		e.Result = new(big.Int).SetBytes(result)
	}
	return e
}

func (e *EthDial) TxCallPaid() *EthDial {
	e.EthFrom = crypto.PubkeyToAddress(e.EthSign.PublicKey)
	e.EthNonce, e.Error = e.EthClient.NonceAt(e.EthContext, e.EthFrom, nil)
	if e.Error == nil {
		utx := types.NewTransaction(e.EthNonce, e.EthAddr, e.EthCoin, e.EthLimit, e.EthPrice, common.Hex2Bytes(e.EthCall+e.EthData))
		stx, err := types.SignTx(utx, types.NewEIP155Signer(e.EthChainId), e.EthSign)
		e.Error = err
		if e.Error == nil {
			e.Error = e.EthClient.SendTransaction(e.EthContext, stx)
			if e.Error == nil {
				e.EthTran = stx.Hash()
				e.EthTranString = stx.Hash().String()
			}
		}
	}
	return e
}

func (e *EthDial) TxSend() *EthDial {
	e.EthFrom = crypto.PubkeyToAddress(e.EthSign.PublicKey)
	e.EthNonce, e.Error = e.EthClient.NonceAt(e.EthContext, e.EthFrom, nil)
	if e.Error == nil {
		utx := types.NewTransaction(e.EthNonce, e.EthAddr, e.EthCoin, e.EthLimit, e.EthPrice, nil)
		stx, err := types.SignTx(utx, types.NewEIP155Signer(e.EthChainId), e.EthSign)
		e.Error = err
		if e.Error == nil {
			e.Error = e.EthClient.SendTransaction(e.EthContext, stx)
			if e.Error == nil {
				e.EthTran = stx.Hash()
				e.EthTranString = stx.Hash().String()
			}
		}
	}
	return e
}

// set the 'signing' (private) address
func (e *EthDial) Sign(a string) *EthDial {
	if e.Error == nil {
		e.EthSign, e.Error = crypto.HexToECDSA(a)
	}
	return e
}

// set the ethereum network chainId
func (e *EthDial) Chain(a *big.Int) *EthDial {
	if e.Error == nil {
		e.EthChainId = a
	}
	return e
}

// set how long for the context to timeout
func (e *EthDial) Wait(a int) *EthDial {
	if e.Error == nil {
		e.EthWait = a
	}
	return e
}

// set the gasLimit
func (e *EthDial) Limit(a uint64) *EthDial {
	if e.Error == nil {
		e.EthLimit = a
	}
	return e
}

// set the amount of coin to be sent
func (e *EthDial) Coin(a *big.Int) *EthDial {
	if e.Error == nil {
		e.EthCoin = a
	}
	return e
}

// set the gasPrice
func (e *EthDial) Price(a *big.Int) *EthDial {
	if e.Error == nil {
		e.EthPrice = a
	}
	return e
}

// set the 'from' address
// TODO this can be calculated from the 'sign' address
func (e *EthDial) From(a string) *EthDial {
	if e.Error == nil {
		e.EthFrom = common.HexToAddress(a)
	}
	return e
}

// set the 'to' address
func (e *EthDial) Addr(a string) *EthDial {
	if e.Error == nil {
		e.EthAddr = common.HexToAddress(a)
	}
	return e
}

// set the called function in the 'addr' contract
func (e *EthDial) Call(a string) *EthDial {
	if e.Error == nil {
		e.EthCall = crypto.Keccak256Hash([]byte(a)).String()[2:10]
	}
	return e
}

// build a 'data' string as part of calling a contract function
func (e *EthDial) DataInt(a int) *EthDial {
	if e.Error == nil {
		e.EthData += leftPadHex(unHex(big.NewInt(int64(a)).Text(16)), 64)
	}
	return e
}

func (e *EthDial) DataBigInt(a *big.Int) *EthDial {
	if e.Error == nil {
		e.EthData += leftPadHex(unHex(a.Text(16)), 64)
	}
	return e
}

//
func (e *EthDial) Tran(wait int, fTran FuncTran) *EthDial {
	e.EthTranWait = wait
	e.EthTranFunc = fTran
	return e
}
