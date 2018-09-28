package matchmaker

import (
	"math"
	"time"

	"github.com/btcsuite/btcutil"

	"dlc"
	"usr"
)

// TfcOffer consists of a Counterparty and FowardCondition and DLC
type TfcOffer struct {
	id     int
	cparty usr.User // Counterparty
	fconds FowardConditions
	dlc    *dlc.Dlc
}

// NewTfcOffer creates a TfcOffer
func NewTfcOffer(
	id int, cparty usr.User, fconds FowardConditions,
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
func (offer TfcOffer) Fconds() FowardConditions {
	return offer.fconds
}

// Dlc returns DLC
func (offer *TfcOffer) Dlc() dlc.Dlc {
	return *offer.dlc
}

func (offer *TfcOffer) makeDlc(isA bool, length int) (*dlc.Dlc, error) {
	settleAt := offer.fconds.settleAt
	namountBtc := offer.fconds.namount
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

// FowardConditions is a confition of TFC offer
type FowardConditions struct {
	namount  float64   // Notional Amount
	vols     float64   // Allowable Volatilities
	famount  float64   // Fund Amount
	rate     float64   // Foward rate
	settleAt time.Time // Settlement Datetime
}

// NewFowardConditions creates a FowardConditions
func NewFowardConditions(
	namount float64, rate float64, vols float64, settleAt time.Time,
) FowardConditions {
	fconds := FowardConditions{
		namount:  namount,
		vols:     vols,
		famount:  namount * namount,
		rate:     rate,
		settleAt: settleAt,
	}

	return fconds
}

func (fconds FowardConditions) Namount() float64 {
	return fconds.namount
}
