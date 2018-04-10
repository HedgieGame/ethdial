pragma solidity ^0.4.21;

contract PeekPoke {
	address                      private  owner;
	mapping(uint256 => uint256)  private  tokenData;

	function PeekPoke() public {
		owner = msg.sender;
	}

	function Peek(uint256 _tokenName) public view returns(uint256) {
		return tokenData[_tokenName];
	}

	function Poke(uint256 _tokenName, uint256 _tokenData) public returns(bool) {
		require(msg.sender == owner);
		tokenData[_tokenName] = _tokenData;
		return true;
	}
}
