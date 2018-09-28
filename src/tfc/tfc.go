package tfc

import "time"

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

// Namount returns notional amount
func (fconds FowardConditions) Namount() float64 {
	return fconds.namount
}

// SettleAt returns settlement datetime
func (fconds FowardConditions) SettleAt() time.Time {
	return fconds.settleAt
}
