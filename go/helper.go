package ethutils

import (
	"UniqueDice/contracts/ambr"
	"UniqueDice/contracts/uniquedice"
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type ClientManager struct {
	RpcClient *rpc.Client
	EthClient *ethclient.Client
}

//"http://127.0.0.1:8545"
func NewClientManager(addr string) (*ClientManager, error) {
	rpcClient, err := rpc.Dial(addr)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	return &ClientManager{
		RpcClient: rpcClient,
		EthClient: ethClient,
	}, nil
}

func (p *ClientManager) Close() {
	if p.RpcClient != nil {
		p.RpcClient.Close()
	}

	if p.EthClient != nil {
		p.EthClient.Close()
	}
}

// GetBlockNumber returns the block number.
func (p *ClientManager) GetBlockNumber(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := p.RpcClient.CallContext(ctx, &result, "eth_blockNumber")
	return (*big.Int)(&result), err
}

func (p *ClientManager) GetAccountBalance(ctx context.Context, account string) (uint64, error) {
	addr := common.HexToAddress(account)
	i, ex := p.GetBlockNumber(ctx)
	if ex != nil {
		return 0, ex
	}

	num, e2 := p.EthClient.BalanceAt(context.Background(), addr, i)
	if e2 != nil {
		return 0, e2
	}
	return num.Uint64(), nil
}

func (p *ClientManager) GetTransaction(ctx context.Context, tx string) (*types.Transaction, bool, error) {
	hash := common.HexToHash(tx)
	return p.EthClient.TransactionByHash(ctx, hash)
}

func (p *ClientManager) GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	return p.EthClient.TransactionByHash(ctx, hash)
}

//privatekey: 09A4E515FD15398BC70030565F45CA831B89C87B545EC238296ABF30939DDDCB
//account: 0xe14c47861b9f20a6BAb730f10e8BB5d4aB420cD4
func (p *ClientManager) DeployAmbrContract(ctx context.Context, keyJson string) error {
	//backend := backends.NewSimulatedBackend()
	//p.EthClient.si
	return nil
}

func StringToPrivateKey(hex string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hex)
}

const (
	uniqueDiceOwnerPrivateKey = "09A4E515FD15398BC70030565F45CA831B89C87B545EC238296ABF30939DDDCB"
	uniqueDiceContractAddress = "0x2ae9f6532799cadfd48ecd99beb3ae9deb396171"
)

func IsInUniqueDiceWhitelist(tokenAddr string) (bool, error) {
	//grant uniquedice
	conn, e := ethclient.Dial(KRinkebyUrl)
	if e != nil {
		return false, e
	}
	defer conn.Close()
	token, ex := uniquedice.NewUniqueDice(common.HexToAddress(uniqueDiceContractAddress), conn)
	if ex != nil {
		return false, ex
	}
	return token.IsTokenAllowed(&bind.CallOpts{}, common.HexToAddress(tokenAddr))
}

func GetTokens(tokenAddr string) (*big.Int, error) {
	//grant uniquedice
	conn, e := ethclient.Dial(KRinkebyUrl)
	if e != nil {
		return nil, e
	}
	defer conn.Close()
	token, ex := ambr.NewERC20(common.HexToAddress(tokenAddr), conn)
	if ex != nil {
		return nil, ex
	}
	return token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(uniqueDiceContractAddress))
}

func ModifyWhiteList(tokenAddr string, value bool) error {
	//grant uniquedice
	conn, e := ethclient.Dial(KRinkebyUrl)
	if e != nil {
		return e
	}
	defer conn.Close()
	privkey, err := StringToPrivateKey(uniqueDiceOwnerPrivateKey)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := uniquedice.NewUniqueDice(common.HexToAddress(uniqueDiceContractAddress), conn)
	if ex != nil {
		return ex
	}
	if value {
		_, ex2 := token.AddTokenToWhiteList(auth, common.HexToAddress(tokenAddr))
		return ex2
	} else {
		_, ex2 := token.RemoveTokenFromWhiteList(auth, common.HexToAddress(tokenAddr))
		return ex2
	}
}

func ModifyWhiteListByPrivateKey(tokenAddr string, privateKey string, value bool) error {
	//grant uniquedice
	conn, e := ethclient.Dial(KRinkebyUrl)
	if e != nil {
		return e
	}
	defer conn.Close()
	privkey, err := StringToPrivateKey(privateKey)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := uniquedice.NewUniqueDice(common.HexToAddress(uniqueDiceContractAddress), conn)
	if ex != nil {
		return ex
	}
	if value {
		_, ex2 := token.AddTokenToWhiteList(auth, common.HexToAddress(tokenAddr))
		return ex2
	} else {
		_, ex2 := token.RemoveTokenFromWhiteList(auth, common.HexToAddress(tokenAddr))
		return ex2
	}
}

func GrantDelegate(playerPrivKey string, tokenAddr string, tokenNumber string) error {
	conn, e := ethclient.Dial(KRinkebyUrl)
	if e != nil {
		return e
	}
	defer conn.Close()
	privkey, err := StringToPrivateKey(playerPrivKey)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(privkey)
	token, ex := ambr.NewERC20(common.HexToAddress(tokenAddr), conn)
	if ex != nil {
		return ex
	}

	n := new(big.Int)
	n.SetString(tokenNumber, 10)
	_, ex2 := token.Approve(auth, common.HexToAddress(uniqueDiceContractAddress), n)
	return ex2
}

func IsAddress(addr string) bool {
	return common.IsHexAddress(addr)
}
