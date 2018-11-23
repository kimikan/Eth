package ethutils_test

import (
	"UniqueDice/ethutils"
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_connect(t *testing.T) {
	c, e := ethutils.NewClientManager(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer c.Close()

	str := "0xe14c47861b9f20a6BAb730f10e8BB5d4aB420cD4"
	addr := common.HexToAddress(str)
	fmt.Println(addr)

	i, ex := c.GetBlockNumber(context.Background())
	fmt.Println(i.Uint64(), ex)

	num, e2 := c.EthClient.BalanceAt(context.Background(), addr, i)
	fmt.Println(num, e2)

	t.Error(c, e)
}

func Test_transaction(t *testing.T) {
	c, e := ethutils.NewClientManager(ethutils.KRinkebyUrl)
	if e != nil {
		t.Error(e)
	}
	defer c.Close()
	str := "0x63479a307dddab6ff4c5788264468f1ec4762caa6aea80982c0d1559ec37d299"
	hash := common.HexToHash(str)
	tx, is, ex := c.EthClient.TransactionByHash(context.Background(), hash)
	/*receipt, err := c.EthClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		log.Fatal(err)
	} */

	fmt.Println(tx.To())

	fmt.Println(tx, is, ex)
	t.Error(ex)
	//c.RpcClient.
}

func Test_deploy(t *testing.T) {

}
