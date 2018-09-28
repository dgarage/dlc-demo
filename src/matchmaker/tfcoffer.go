package matchmaker

import (
	"math"

	"github.com/btcsuite/btcutil"

	"dlc"
	"tfc"
	"usr"
)

// TfcOffer consists of a Counterparty and FowardCondition and DLC
type TfcOffer struct {
	id     int
	cparty usr.User // Counterparty
	fconds tfc.FowardConditions
	dlc    *dlc.Dlc
}

// NewTfcOffer creates a TfcOffer
func NewTfcOffer(
	id int, cparty usr.User, fconds tfc.FowardConditions,
) *TfcOffer {
	offer := &TfcOffer{
		id:     111,
		cparty: cparty,
		fconds: fconds,
	}
	return offer
}

// ID returns the id of the offer
func (offer *TfcOffer) ID() int {
	return offer.id
}

// Fconds returns foward conditions of the offer
func (offer TfcOffer) Fconds() tfc.FowardConditions {
	return offer.fconds
}

// Dlc returns DLC
func (offer *TfcOffer) Dlc() dlc.Dlc {
	return *offer.dlc
}

func (offer *TfcOffer) makeDlc(isA bool, length int) (*dlc.Dlc, error) {
	settleAt := offer.fconds.SettleAt()
	namountBtc := offer.fconds.Namount()
	namountSat := int64(namountBtc * btcutil.SatoshiPerBitcoin)
	// TODO: need to calculate fees?
	fefee := int64(10)                      // fund transaction estimate fee satoshi/byte
	sefee := int64(10)                      // settlement transaction estimate fee satoshi/byte
	sfee := dlc.DlcSettlementTxSize * sefee // settlement transaction size is 345 bytes
	d, err := dlc.NewDlc(half(namountSat), half(namountSat), fefee,
		sefee, half(sfee), half(sfee), isA)
	if err != nil {
		return nil, err
	}
	d.SetGameConditions(settleAt, length)
	return d, nil
}

func half(value int64) int64 {
	return int64(math.Ceil(float64(value) / float64(2)))
}
