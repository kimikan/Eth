package uniquedice_test

import (
	"UniqueDice/contracts/ambr"
	"UniqueDice/contracts/uniquedice"
	"UniqueDice/ethutils"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

const key = `
{
  "address": "1a9ec3b0b807464e6d3398a59d6b0a369bf422fa",
  "crypto": {
    "cipher": "aes-128-ctr",
    "ciphertext": "a471054846fb03e3e271339204420806334d1f09d6da40605a1a152e0d8e35f3",
    "cipherparams": {
      "iv": "44c5095dc698392c55a65aae46e0b5d9"
    },
    "kdf": "scrypt",
    "kdfparams": {
      "dklen": 32,
      "n": 262144,
      "p": 1,
      "r": 8,
      "salt": "e0a5fbaecaa3e75e20bccf61ee175141f3597d3b1bae6a28fe09f3507e63545e"
    },
    "mac": "cb3f62975cf6e7dfb454c2973bdd4a59f87262956d5534cdc87fb35703364043"
  },
  "id": "e08301fb-a263-4643-9c2b-d28959f66d6a",
  "version": 3
}
`

const ContractOwnerPrivateKey = "09A4E515FD15398BC70030565F45CA831B89C87B545EC238296ABF30939DDDCB"
const AmbrContractAddress = "0x965f6fe367Ca00dA87c67669F56442747aA13CD8"
const MajiaAddress = "0x0F13249F2d73f65b9ec319EFaC1D43BAD91eb4C2"
const MajiaPrivateKey = "AB7486B269CF5A77F91196631A19EA1BABF89238EF6815F51946A72535E976FD"
const UniqueDiceContractAddress = "0x2ae9f6532799cadfd48ecd99beb3ae9deb396171"

func Test_deployambr(t *testing.T) {
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	// Deploy a new awesome contract for the binding demo
	address, tx, token, err := ambr.DeployAmbrToken(auth, conn)
	if err != nil {
		t.Error("Failed to deploy new token contract: ", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	// Don't even wait, check its presence in the local pending state
	time.Sleep(1250 * time.Millisecond) // Allow it to be processed by the local node :P

	t.Error(token)
}

func Test_deploy_uniquedice(t *testing.T) {
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	// Deploy a new awesome contract for the binding demo
	address, tx, token, err := uniquedice.DeployUniqueDice(auth, conn)
	if err != nil {
		t.Error("Failed to deploy new token contract: ", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	// Don't even wait, check its presence in the local pending state
	time.Sleep(1250 * time.Millisecond) // Allow it to be processed by the local node :P

	t.Error(token)
}

func Test_deploy_uniquedice_others(t *testing.T) {
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	// Deploy a new awesome contract for the binding demo
	address, tx, token, err := uniquedice.DeploySafeMath(auth, conn)
	if err != nil {
		t.Error("Failed to deploy new token contract: ", err)
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())

	// Don't even wait, check its presence in the local pending state
	time.Sleep(1250 * time.Millisecond) // Allow it to be processed by the local node :P

	t.Error(token)
}

func Test_getAmbrTotal(t *testing.T) {
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}

	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := ambr.NewAmbrToken(common.HexToAddress(AmbrContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}

	tx, ee := token.SimpleToken(auth)
	if ee != nil {
		t.Error(ee)
	}

	fmt.Println(tx)

	total, ex := token.TotalSupply(&bind.CallOpts{})
	if ex != nil {
		t.Error(ex)
	}
	fmt.Println("total: ", total)

	name, ex2 := token.Name(&bind.CallOpts{})
	if ex2 != nil {
		t.Error(ex2)
	}
	t.Error(name)
}

func Test_transferambrToMajia(t *testing.T) {
	//majia-addr: 0x0F13249F2d73f65b9ec319EFaC1D43BAD91eb4C2
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := ambr.NewAmbrToken(common.HexToAddress(AmbrContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}

	bi := big.NewInt(10000000000)
	bi = bi.Mul(bi, big.NewInt(1000000000))
	bi = bi.Mul(bi, big.NewInt(10000000000))
	tx, ex2 := token.Transfer(auth, common.HexToAddress(MajiaAddress), bi)
	if ex2 == nil {
		fmt.Println("ok, transfer done")
	}
	t.Error(tx, ex2)
}

func Test_transferambrToUniqueDice(t *testing.T) {
	//majia-addr: 0x0F13249F2d73f65b9ec319EFaC1D43BAD91eb4C2
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := ambr.NewAmbrToken(common.HexToAddress(AmbrContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}

	bi := big.NewInt(10000000000)
	bi = bi.Mul(bi, big.NewInt(1000000000))
	bi = bi.Mul(bi, big.NewInt(10000000000))
	tx, ex2 := token.Transfer(auth, common.HexToAddress(UniqueDiceContractAddress), bi)
	if ex2 == nil {
		fmt.Println("ok, transfer done")
	}
	t.Error(tx, ex2)
}

//grant uniquedice contract to delegate majia on ambr token.
func Test_grantUniquediceToDelegateMajia(t *testing.T) {
	//grant uniquedice
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()
	privkey, err := ethutils.StringToPrivateKey(MajiaPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := ambr.NewAmbrToken(common.HexToAddress(AmbrContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}

	bi := big.NewInt(1000000)
	bi = bi.Mul(bi, big.NewInt(1000000000))
	bi = bi.Mul(bi, big.NewInt(10000000000))
	tx, ex2 := token.Approve(auth, common.HexToAddress(UniqueDiceContractAddress), bi)

	t.Error(tx, ex2)
}

func Test_addAmbrToWhiteList(t *testing.T) {
	//grant uniquedice
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()
	privkey, err := ethutils.StringToPrivateKey(ContractOwnerPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)

	token, ex := uniquedice.NewUniqueDice(common.HexToAddress(UniqueDiceContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}

	tx, ex2 := token.AddTokenToWhiteList(auth, common.HexToAddress(AmbrContractAddress))

	t.Error(tx, ex2)
}

//play with majia
func Test_playGameWithAmbr(t *testing.T) {
	//majia-addr: 0x0F13249F2d73f65b9ec319EFaC1D43BAD91eb4C2
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()

	privkey, err := ethutils.StringToPrivateKey(MajiaPrivateKey)
	if err != nil {
		t.Error(err)
	}
	auth := bind.NewKeyedTransactor(privkey)
	//auth.GasLimit = 1111111

	/*
		gp, _ := conn.SuggestGasPrice(context.Background())
		auth.GasPrice = gp.Mul(gp, big.NewInt(5))
	*/
	fmt.Println(auth)
	token, ex := uniquedice.NewUniqueDice(common.HexToAddress(UniqueDiceContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}
	ok, ee := token.IsTokenAllowed(&bind.CallOpts{}, common.HexToAddress(AmbrContractAddress))
	if ee != nil {
		t.Error(ee)
	}
	fmt.Println("Support: ", ok)

	tx, ex2 := token.DelegatePlay(auth, common.HexToAddress(AmbrContractAddress), 98, 8888)

	t.Error(tx, ex2)
}

func Test_getMajiaAmbrToken(t *testing.T) {
	conn, e := ethclient.Dial(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer conn.Close()
	token, ex := ambr.NewAmbrToken(common.HexToAddress(AmbrContractAddress), conn)
	if ex != nil {
		t.Error(ex)
	}
	tx, ex2 := token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(MajiaAddress))

	t.Error(tx, ex2)
}

/*
func Test_deploy(t *testing.T) {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	alloc := make(core.GenesisAlloc)
	alloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(133700000)}
	sim := backends.NewSimulatedBackend(alloc)

	// deploy contract
	addr, _, contract, err := contracts.DeployERC20(auth, sim)
	if err != nil {
		log.Fatalf("could not deploy contract: %v", err)
	}

	// interact with contract
	fmt.Printf("Contract deployed to %s\n", addr.String())
	deadlineCampaign, _ := contract.DeadlineCampaign(nil)
	fmt.Printf("Pre-mining Campaign Deadline: %s\n", deadlineCampaign)

	fmt.Println("Mining...")
	// simulate mining
	sim.Commit()

	postDeadlineCampaign, _ := contract.DeadlineCampaign(nil)
	fmt.Printf("Post-mining Campaign Deadline: %s\n", time.Unix(postDeadlineCampaign.Int64(), 0))

	// create a project
	numOfProjects, _ := contract.NumberOfProjects(nil)
	fmt.Printf("Number of Projects before: %d\n", numOfProjects)

	fmt.Println("Adding new project...")
	contract.SubmitProject(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasLimit: 2381623,
		Value:    big.NewInt(10),
	}, "test project", "http://www.example.com")

	fmt.Println("Mining...")
	sim.Commit()

	numOfProjects, _ = contract.NumberOfProjects(nil)
	fmt.Printf("Number of Projects after: %d\n", numOfProjects)
	info, _ := contract.GetProjectInfo(nil, auth.From)
	fmt.Printf("Project Info: %v\n", info)

	// instantiate deployed contract
	fmt.Printf("Instantiating contract at address %s...\n", auth.From.String())
	instContract, err := NewWinnerTakesAll(addr, sim)
	if err != nil {
		log.Fatalf("could not instantiate contract: %v", err)
	}
	numOfProjects, _ = instContract.NumberOfProjects(nil)
	fmt.Printf("Number of Projects of instantiated Contract: %d\n", numOfProjects)
}
*/
