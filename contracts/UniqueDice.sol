pragma solidity ^0.4.18;


/**
 * @title SafeMath
 * @dev Math operations with safety checks that throw on error
 */
library SafeMath {

    /**
    * @dev Multiplies two numbers, throws on overflow.
    */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }
        uint256 c = a * b;
        assert(c / a == b);
        return c;
    }

    /**
    * @dev Integer division of two numbers, truncating the quotient.
    */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        // assert(b > 0); // Solidity automatically throws when dividing by 0
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold
        return c;
    }

    /**
    * @dev Substracts two numbers, throws on overflow (i.e. if subtrahend is greater than minuend).
    */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        assert(b <= a);
        return a - b;
    }

    /**
    * @dev Adds two numbers, throws on overflow.
    */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        assert(c >= a);
        return c;
    }
}


/**
 * @title ERC20Basic
 * @dev Simpler version of ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/179
 */
contract ERC20Basic {
    function totalSupply() public view returns (uint256);
    function balanceOf(address who) public view returns (uint256);
    function transfer(address to, uint256 value) public returns (bool);
    event Transfer(address indexed from, address indexed to, uint256 value);
}

/**
 * @title ERC20 interface
 * @dev see https://github.com/ethereum/EIPs/issues/20
 */
contract ERC20 is ERC20Basic {
    function allowance(address owner, address spender) public view returns (uint256);
    function transferFrom(address from, address to, uint256 value) public returns (bool);
    function approve(address spender, uint256 value) public returns (bool);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

contract UniqueDice {
    using SafeMath for uint256;
    
    address public owner;
    //event OwnershipTransferred(address indexed _from, address indexed _to);

    constructor() public {
        //store the contract owner address to set some 
        //permission operation
        owner = msg.sender;
    }

    mapping(address => uint64) Amounts;
    mapping(address => bool) AllowedTokens;

    function isBetWin(uint64 rolls) private view returns
        (bool) {
        uint256 seed = uint256(keccak256(abi.encodePacked(
            (block.timestamp).add
            (block.difficulty).add
            ((uint256(keccak256(abi.encodePacked(block.coinbase)))) / (now)).add
            (block.gaslimit).add
            ((uint256(keccak256(abi.encodePacked(msg.sender)))) / (now)).add
            (block.number)
        )));

        if((seed - ((seed / 100) * 100)) < rolls)
            return(true);
        else
            return(false);
    }

    function computePayoutOnWin(uint64 rolls, uint64 amount) private pure returns 
        (uint64) {
        uint64 accurateX = 98 * 10000 / (rolls - 1);
        uint64 payoutOnWin = accurateX * amount / 10000;
        return payoutOnWin;
    }

    function doPlay(ERC20 token, address fromAddr, uint64 rolls, uint64 payoutOnWin) private returns 
        (uint64) {
        bool win = isBetWin(rolls);
        if (!win) {
            return 0;
        }

        bool ok = token.transfer(fromAddr, payoutOnWin);
        //bool ok = token.transferFrom(this, fromAddr, payoutOnWin);
        if (ok) {
            return payoutOnWin;
        }
        
        return 0;
    }

    //a white list that support contract owner control the use
    function addTokenToWhiteList(address tokenAddr) public {
        require(msg.sender == owner, "only owner can operation this");
        AllowedTokens[tokenAddr] = true;
    }

    //sometime we have to remove it due to someother reasons
    function removeTokenFromWhiteList(address tokenAddr) public {
        require(msg.sender == owner, "only owner can operation this");
        AllowedTokens[tokenAddr] = false;
    }

    //support for ownership transfer for debugging usage
    function transferOwnership(address newOwner) public returns
        (bool) {
        require(msg.sender == owner, "Operation must be invoked by contract owner");
        owner = newOwner;
        return true;
    }

    function isTokenAllowed(address tokenAddr) public view returns 
        (bool) {
        bool value = AllowedTokens[tokenAddr];
        return value;
    }
	
	//withdraw some token to target address
	function withdraw(address tokenAddr, uint64 amount, address target) public returns 
		(bool) {
		require(msg.sender == owner, "Operation must be invoked by contract owner");
		ERC20 token = ERC20(tokenAddr);
		bool ok = token.transfer(target, amount);
        
		return ok;
	}
	
    //the only play entrance
    function delegatePlay(address tokenAddr, uint64 rolls, uint64 amount) public returns
        (uint64) {
        require(rolls >= 1 && rolls < 100, "invalid rolls");
        require(amount > 0, "0 not allowed");
        require(isTokenAllowed(tokenAddr), "token not allowed");
        //require(computePayoutOnWin(rolls, amount));
        //require.....
        uint64 payoutOnWin = computePayoutOnWin(rolls, amount);
        ERC20 token = ERC20(tokenAddr);
        require(token.balanceOf(this) > payoutOnWin, "has no enough token to handle this play");

        bool ok = token.transferFrom(msg.sender, this, amount);
        if (ok) {
            Amounts[msg.sender] += amount;
            return doPlay(token, msg.sender, rolls, payoutOnWin);
        }

        return 0;
    } //end delegatePlay

    /*
    function approve(address tokenAddr, uint allowed) public {
        ERC20 token = ERC20(tokenAddr);
        token.approve(msg.sender, allowed);
    }
    */
}