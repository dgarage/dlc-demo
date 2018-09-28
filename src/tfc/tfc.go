package tfc

import "time"

// FowardConditions is a confition of TFC offer
type FowardConditions struct {
	namount    float64   // Notional Amount
	fundRate   float64   // Fund Rate
	famount    float64   // Fund Amount
	fowardRate float64   // Foward rate
	settleAt   time.Time // Settlement Datetime
}

// NewFowardConditions creates a FowardConditions
func NewFowardConditions(
	namount float64, fundRate float64, fowardRate float64, settleAt time.Time,
) FowardConditions {
	fconds := FowardConditions{
		namount:    namount,
		fundRate:   fundRate,
		famount:    namount * fundRate,
		fowardRate: fowardRate,
		settleAt:   settleAt,
	}

	return fconds
}

// Namount returns notional amount
func (fconds FowardConditions) Namount() float64 {
	return fconds.namount
}

// FowardRate returns forward rate
func (fconds FowardConditions) FowardRate() float64 {
	return fconds.fowardRate
}

// FundRate returns fund rate
func (fconds FowardConditions) FundRate() float64 {
	return fconds.fundRate
}

// Vols returns allowable volatilities
// func (fconds FowardConditions) Vols() float64 {
// 	return fconds.vols
// }

// SettleAt returns settlement datetime
func (fconds FowardConditions) SettleAt() time.Time {
	return fconds.settleAt
}
