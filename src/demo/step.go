// step.go
package main

import (
	"fmt"
	"log"
	"time"

	"dlc"
	"matchmaker"
	"tfc"
	"usr"
)

func stepBobPutTfcOfferOnBoard(num int, d *Demo) error {
	var err error
	var date time.Time
	date, err = time.Parse("2006-01-02", "2018-10-31")
	cparty := d.bob
	// TODO: adjust values
	vol := 0.5
	frate := 0.1

	amts := []float64{1, 2}

	for i, amt := range amts {
		ID := i + 1
		fconds := tfc.NewFowardConditions(amt, frate, vol, date)
		offer := *matchmaker.NewTfcOffer(ID, *cparty, fconds)

		if err = d.mm.PutOffer(offer); err != nil {
			return err
		}
	}

	return nil
}

func stepAliceSendOfferToBob(num int, d *Demo) error {
	fmt.Printf("step%d : Take TFC Offer on MatchMaker board\n", num)
	tfcoffer := d.mm.Offers()[0]
	var err error
	tfcoffer, err = d.mm.TakeOffer(tfcoffer.ID())
	if err != nil {
		return err
	}

	dlc := tfcoffer.Dlc()
	err = aliceSendOfferToBob(num, d, &dlc)
	return err
}

func aliceSendOfferToBob(num int, d *Demo, dlc *dlc.Dlc) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)

	log.Printf("step%d : Alice GetOfferData\n", num)
	var odata []byte
	var err error
	odata, err = d.alice.GetOfferData(dlc)
	if err != nil {
		return err
	}
	log.Printf("step%d : Alice SetOracleKeys\n", num)
	keys, err := d.olivia.Keys(d.alice.GameDate())
	if err != nil {
		return err
	}
	err = d.alice.SetOracleKeys(keys)
	if err != nil {
		return err
	}
	log.Printf("step%d : Alice -> Bob\n", num)
	// dump(odata)
	log.Printf("step%d : Bob SetOfferData\n", num)
	err = d.bob.SetOfferData(odata)
	if err != nil {
		return err
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}

func stepBobSendAcceptToAlice(num int, d *Demo) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)
	log.Printf("step%d : Bob SetOracleKeys\n", num)
	keys, err := d.olivia.Keys(d.bob.GameDate())
	if err != nil {
		return err
	}
	err = d.bob.SetOracleKeys(keys)
	if err != nil {
		return err
	}
	log.Printf("step%d: Bob GetAcceptData\n", num)
	adata, err := d.bob.GetAcceptData()
	if err != nil {
		return err
	}
	log.Printf("step%d : Bob -> Alice\n", num)
	// dump(adata)
	log.Printf("step%d : Alice SetAcceptData\n", num)
	err = d.alice.SetAcceptData(adata)
	if err != nil {
		return err
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}

func stepAliceSendSignToBob(num int, d *Demo) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)
	log.Printf("step%d : Alice GetSignData\n", num)
	sdata, err := d.alice.GetSignData()
	if err != nil {
		return err
	}
	log.Printf("step%d : Alice -> Bob\n", num)
	// dump(sdata)
	log.Printf("step%d : Bob SetSignData\n", num)
	err = d.bob.SetSignData(sdata)
	if err != nil {
		return err
	}
	err = d.bob.SendFundTx()
	if err != nil {
		return err
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}

func stepAliceAndBobSetOracleSign(num int, d *Demo) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)
	date := d.alice.GameDate()
	sigs, err := d.olivia.Signs(date)
	if err != nil {
		return err
	}
	log.Printf("step%d : Alice & Bob SetOracleSigns\n", num)
	err = d.alice.SetOracleSigns(sigs)
	if err != nil {
		return err
	}
	date = d.bob.GameDate()
	sigs, err = d.olivia.Signs(date)
	if err != nil {
		return err
	}
	err = d.bob.SetOracleSigns(sigs)
	if err != nil {
		return err
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}

func stepAliceOrBobSendSettlementTx(num int, demo *Demo) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)
	users := []*usr.User{demo.alice, demo.bob}
	for _, user := range users {
		err := user.SendSettlementTx()
		if err != nil {
			fmt.Printf("SendSettlementTx error : %+v\n", err)
			continue
		}
		err = user.SendSettlementTxTo(int64(10))
		if err != nil {
			return err
		}
		break
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}

func stepAliceOrBobSendRefundTx(num int, demo *Demo) error {
	s := time.Now()
	log.Printf("begin step%d\n", num)
	user := demo.alice
	err := user.SendRefundTx()
	if err != nil {
		return err
	}
	log.Printf("end   step%d %f sec\n", num, (time.Now()).Sub(s).Seconds())
	return nil
}
