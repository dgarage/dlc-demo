// cmds.go
package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/olekukonko/tablewriter"

	"matchmaker"
	"usr"
)

type cmd struct {
	n []string
	f func([]string, *Demo) error
}

func listCmds() []*cmd {
	list := []*cmd{}
	list = append(list, &cmd{[]string{"step", "s"}, step})
	list = append(list, &cmd{[]string{"set"}, set})
	list = append(list, &cmd{[]string{"generate", "g"}, generate})
	list = append(list, &cmd{[]string{"getrawtransaction", "grt"}, getrawtransaction})
	list = append(list, &cmd{[]string{"decodescript", "ds"}, decodescript})
	list = append(list, &cmd{[]string{"balance", "b"}, balance})
	list = append(list, &cmd{[]string{"fee"}, txfee})
	list = append(list, &cmd{[]string{"faucet"}, faucet})
	list = append(list, &cmd{[]string{"setval"}, setval})
	// commands for tfc demo
	list = append(list, &cmd{[]string{"wallet"}, walletRoot})
	list = append(list, &cmd{[]string{"offers"}, offersRoot})
	list = append(list, &cmd{[]string{"offer"}, offerRoot})
	list = append(list, &cmd{[]string{"contracts"}, contractsRoot})
	list = append(list, &cmd{[]string{"contract"}, contractRoot})
	return list
}

func generate(args []string, d *Demo) error {
	var err error
	nblocks := 1
	if len(args) > 1 {
		nblocks, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	if nblocks < 1 {
		return fmt.Errorf("nblocks is less than or equal to zero. %d", nblocks)
	}
	res, err := d.rpc.Request("generate", nblocks)
	if err != nil {
		return err
	}
	fmt.Printf("generate %d\n", nblocks)
	bs, err := json.Marshal(res.Result)
	if err != nil {
		return err
	}
	dump(bs)
	return nil
}

func getrawtransaction(args []string, d *Demo) error {
	if len(args) < 2 {
		return fmt.Errorf("illegal parameter")
	}
	txid := args[1]
	res, err := d.rpc.Request("getrawtransaction", txid, 1)
	if err != nil {
		return err
	}
	fmt.Printf("getrawtransaction %s 1\n", txid)
	bs, err := json.Marshal(res.Result)
	if err != nil {
		return err
	}
	dump(bs)
	return nil
}

func decodescript(args []string, d *Demo) error {
	if len(args) < 2 {
		return fmt.Errorf("illegal parameter")
	}
	hexstring := args[1]
	res, err := d.rpc.Request("decodescript", hexstring)
	if err != nil {
		return err
	}
	fmt.Printf("decodescript %s\n", hexstring)
	bs, err := json.Marshal(res.Result)
	if err != nil {
		return err
	}
	dump(bs)
	return nil
}

func balance(args []string, d *Demo) error {
	amta := d.alice.GetBalance()
	amtb := d.bob.GetBalance()
	fmt.Printf("alice amount : %.8f BTC\n", float64(amta)/btcutil.SatoshiPerBitcoin)
	fmt.Printf("bob   amount : %.8f BTC\n", float64(amtb)/btcutil.SatoshiPerBitcoin)
	return nil
}

func faucet(args []string, d *Demo) error {
	var err error
	satoshi := int(1 * btcutil.SatoshiPerBitcoin)
	if len(args) > 1 {
		satoshi, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	if satoshi < 1 {
		return fmt.Errorf("satoshi is less than or equal to zero. %d", satoshi)
	}
	s := time.Now()
	log.Printf("begin faucet\n")
	_, err = d.rpc.Request("generate", 1)
	if err != nil {
		return err
	}
	lowest := int64(satoshi)
	users := []*usr.User{d.alice, d.bob}
	for _, user := range users {
		amt := user.GetBalance()
		if amt < lowest {
			_, err = d.rpc.Request("sendtoaddress", user.GetAddress(), float64(lowest-amt)/btcutil.SatoshiPerBitcoin)
			if err != nil {
				return err
			}
			_, err = d.rpc.Request("generate", 1)
			if err != nil {
				return err
			}
		}
	}
	// balance(nil, d)
	log.Printf("end   faucet %f sec\n", (time.Now()).Sub(s).Seconds())
	return nil
}

func txfee(args []string, d *Demo) error {
	if len(args) < 2 {
		return fmt.Errorf("illegal parameter")
	}
	txid := args[1]
	res, err := d.rpc.Request("getrawtransaction", txid)
	if err != nil {
		return err
	}
	str, _ := res.Result.(string)
	bs, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	tx, err := bsToMsgTx(bs)
	if err != nil {
		return err
	}
	iamt := int64(0)
	for _, txin := range tx.TxIn {
		amt, err := getAmount(d, txin)
		if err != nil {
			return err
		}
		iamt += amt
	}
	oamt := int64(0)
	for _, txout := range tx.TxOut {
		oamt += txout.Value
	}
	fmt.Printf("input:%d output:%d fee:%d size:%d efee:%f\n",
		iamt, oamt, iamt-oamt, len(bs), float64(iamt-oamt)/float64(len(bs)))
	return nil
}

func getAmount(d *Demo, txin *wire.TxIn) (int64, error) {
	op := txin.PreviousOutPoint
	res, err := d.rpc.Request("getrawtransaction", op.Hash.String())
	if err != nil {
		return 0, err
	}
	str, _ := res.Result.(string)
	bs, err := hex.DecodeString(str)
	if err != nil {
		return 0, err
	}
	tx, err := bsToMsgTx(bs)
	if err != nil {
		return 0, err
	}
	if uint32(len(tx.TxOut)) <= op.Index {
		return 0, fmt.Errorf("out of range : %d,%d", len(tx.TxOut), op.Index)
	}
	txout := tx.TxOut[op.Index]
	return txout.Value, nil
}

func bsToMsgTx(bs []byte) (*wire.MsgTx, error) {
	var tx *wire.MsgTx
	tx = &wire.MsgTx{}
	buf := &bytes.Buffer{}
	_, err := buf.Write(bs)
	if err != nil {
		return nil, err
	}
	err = tx.Deserialize(buf)
	if err != nil {
		tx = &wire.MsgTx{}
		buf := &bytes.Buffer{}
		_, err := buf.Write(bs)
		if err != nil {
			return nil, err
		}
		err = tx.DeserializeNoWitness(buf)
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}

func dump(bs []byte) {
	var buf bytes.Buffer
	err := json.Indent(&buf, bs, "", "  ")
	if err != nil {
		fmt.Printf("dump Error : %+v\n", err)
		return
	}
	fmt.Printf("%s\n", buf.String())
}

func setval(args []string, d *Demo) error {
	if len(args) < 3 {
		return fmt.Errorf("illegal parameter")
	}
	err := d.olivia.SetVals(args[1], args[2])
	if err != nil {
		return err
	}
	return nil
}

func walletRoot(args []string, d *Demo) error {
	var err error
	switch subcmd := args[1]; subcmd {
	case "balance":
		err = showBalance(d)
	}
	return err
}

func showBalance(d *Demo) error {
	bSat := d.alice.GetBalance()
	bBtc := float64(bSat) / btcutil.SatoshiPerBitcoin

	fmt.Printf("Amount: %.8f BTC\n", bBtc)
	return nil
}

// commands for TFC demo

func offersRoot(args []string, d *Demo) error {
	var err error
	switch subcmd := args[1]; subcmd {
	case "list":
		err = listTfcOffers(d)
	}
	return err
}

// command: offers list
func listTfcOffers(d *Demo) error {
	// Add dummy data
	if len(d.mm.Offers()) < 1 {
		err := stepBobPutTfcOfferOnBoard(0, d)
		if err != nil {
			return err
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Notional Amount", "Forward Rate", "Settle Date"})
	for _, o := range d.mm.Offers() {
		id := strconv.Itoa(o.ID())
		fconds := o.Fconds()
		namount := fmt.Sprintf("%.8f BTC", fconds.Namount())
		frate := fmt.Sprintf("%.8f JPY/BTC", fconds.Rate())
		// tFormat := "2006-01-02 15:04:05"
		tFormat := "2006-01-02"
		settleAt := fconds.SettleAt().Format(tFormat)
		trow := []string{id, namount, frate, settleAt}
		table.Append(trow)
	}
	table.Render()
	return nil
}

func offerRoot(args []string, d *Demo) error {
	var err error
	switch subcmd := args[1]; subcmd {
	case "take":
		err = takeTfcOffer(args[2:], d)
	}
	return err
}

// command: offer take 111
func takeTfcOffer(args []string, d *Demo) error {
	var err error
	var offerID int
	offerID, err = strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	var tfcoffer matchmaker.TfcOffer
	tfcoffer, err = d.mm.TakeOffer(offerID)
	if err != nil {
		return err
	}

	fmt.Println("Sending DLC to counterparty")
	fmt.Println("TODO: add more logs")
	dlc := tfcoffer.Dlc()
	if err = aliceSendOfferToBob(1, d, &dlc); err != nil {
		return err
	}

	fmt.Println("Waiting for counterpaty to accept")
	fmt.Println("TODO: add more logs")
	if err = stepBobSendAcceptToAlice(2, d); err != nil {
		return err
	}

	fmt.Println("Sending sign to counterparty")
	fmt.Println("TODO: add more logs")
	if err = stepAliceSendSignToBob(3, d); err != nil {
		return err
	}

	// TODO: Save TFC locally
	fconds := tfcoffer.Fconds()
	d.alice.Fconds = &fconds

	// implicitly change the oracle status for demo
	date := dlc.GameDate().Format("20060102")
	fixingRate := "30"
	if err = d.olivia.SetVals(date, fixingRate); err != nil {
		return err
	}

	return nil
}

func contractsRoot(args []string, d *Demo) error {
	var err error
	switch subcmd := args[1]; subcmd {
	case "list":
		err = listContracts(d)
	}
	return err
}

// command: contracts list
func listContracts(d *Demo) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Notional Amount", "Status"})

	fmt.Println("TODO: Add more fields and dummy records")
	fconds := d.alice.Fconds
	id := "TODO: Generate dummy ID"
	namount := strconv.FormatFloat(fconds.Namount(), 'f', -1, 64)
	status := "Fixed"
	trow := []string{id, namount, status}
	table.Append(trow)

	table.Render()
	return nil
}

func contractRoot(args []string, d *Demo) error {
	var err error
	switch subcmd := args[1]; subcmd {
	case "settle":
		err = settleContract(args[2:], d)
	}
	return err
}

// command: contract settle
func settleContract(args []string, d *Demo) error {
	fmt.Println("TODO: add more logs")

	var err error
	if err = stepAliceAndBobSetOracleSign(4, d); err != nil {
		return err
	}
	if err = stepAliceOrBobSendSettlementTx(5, d); err != nil {
		return err
	}
	return nil
}
