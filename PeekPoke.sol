pragma solidity ^0.4.21;

contract PeekPoke {
	address                      public   owner;
	mapping(uint256 => uint256)  private  data;

	function PeekPoke() public {
		owner = msg.sender;
	}

	function Peek(uint256 _name) public view returns(uint256) {
		return data[_name];
	}

	function Poke(uint256 _name, uint256 _data) public returns(bool) {
		require(msg.sender == owner);
		data[_name] = _data;
		return true;
	}
}
