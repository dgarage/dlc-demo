package matchmaker

import (
	"dlc"
	"time"
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

func (offer *TfcOffer) Fconds() FowardConditions {
	return offer.fconds
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
	namount float64, rate float64, settleAt time.Time,
) FowardConditions {
	vols := 0.5
	fconds := FowardConditions{
		namount:  namount,
		vols:     vols,
		famount:  namount * namount,
		rate:     rate,
		settleAt: settleAt,
	}

	return fconds
}
