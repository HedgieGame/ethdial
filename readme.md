ethdial

**Current Problems**

**TODO** The test run in main/main.go results in a failed transaction on Rinkeby.  This is suspected to be incorrect gas settings, but the debugger on etherscan.io isn't helping much.  This is at least better than the several days spent discovering that ganache-cli is returning a different Transaction ID than the one it reports to the screen.  They've acknowledged the bug have not released a fix.

**Overview**

This code acts as an abstraction interface between the Hedgie creature manipulation code and the blockchain which is to be used as a system of record for ownership, characteristics, etc.  The current state of the Ethereum VM disallows calling functions with sufficient variables to store the several characteristics of each Hedgie, so the approach is taken of taking a Hedgie struct and packing the several elements into a single 256-bit word, which is then stored and retrieved via a simple Solidity contract that acts on just an Address and a Value.  The contract is in peekpoke.sol (one potential improvement is to add a *collection* parameter, which will be hashed with the Hedgie.HID, allowing several different stores of values).
 
**TODO** Evaluate adding a *collection* parameter to dsicriminate between multipel storage, err, collections.

The two entry points are Hedgie.Peek() and Hedgie.Poke(), in the peek.go file.

**TODO** Believe it better to rewrite this to not depend on a Hedgie{} but to use parameters.

The file *pack.go* does the packing and unpacking of multiple big-ints to a single 256-bit and vice versa.

**How it works**

github.com/ethereum/go-ethereum is where the Ethereum code lives.  Perusal of that code would lead one to believe that they are doing something special with the *Go*-ish code structure to send Ether, call a contract function with no payment required (any pure reads), and call a contract function with payment required (a write).  The distinction of the last two matters, and the Ethereum repository code might lead one down the wrong path of believing they are called similarly.  If you dig deeply enough into this code, then you discover they are doing nothing more than making an RPC call just like is done using curl.  The curl approach is quite well described at https://github.com/ethereum/wiki/wiki/JSON-RPC and frankly would be the better approach as it's consistent and cleaner in how things are done.  I've started a side project to write something using the curl approach but this code using the functions provided in the Ethereum repo are what's available for now.
