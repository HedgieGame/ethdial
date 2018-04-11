ethdial

**Current Problems**

The test run in main/main.go results in a failed transaction on Rinkeby.  This is suspected to be incorrect gas settings, but the debugger on ethersca.io isn't helping much.  This is at least better than the several days spent discovering that ganache-cli is returning a different Transaction ID than the one it reports to the screen.  They've acknowledged the bug bug have not released a fix.

**Overview**

This code acts as an abstraction interface between the Hedgie creature manipulation code and the blockchain which is to be used as a system of record for ownership, characteristics, etc.  The current state of the Ethereum VM disallows calling functions with sufficient variables to store the several characteristics of each Hedgie, so an approach is taken her of taking a Hedgie struct and packing the several elements into a single 256-bit word, which is then stored and retrieved via a simple Solidity contract that acts on just an Address and a Value.  The contract is in peekpoke.sol (one potential improvement is to add a *collection* parameter, which will be hashed with the Hedgie.HID, allowing several different stores of values).

The two entry points are Hedgie.Peek() and Hedgie.Poke(), in the peek.go file.
*** Believe it better to rewrite this to not depend on a Hedgie{} but to use parameters. ***
