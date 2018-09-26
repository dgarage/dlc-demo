package matchmaker

// MatchMaker matches two parties
type MatchMaker struct {
	offers []*TfcOffer
}

// NewMatchMaker creates a MatchMaker
func NewMatchMaker() MatchMaker {
	return MatchMaker{}
}

// Offers returns TFC offers on board
func (mm *MatchMaker) Offers() []TfcOffer {
	var offers []TfcOffer
	for _, offer := range mm.offers {
		offers = append(offers, *offer)
	}
	return offers
}

// PutOffer puts a new TFC offer on board
func (mm *MatchMaker) PutOffer(offer TfcOffer) error {
	// TODO: validate offer
	mm.offers = append(mm.offers, &offer)
	return nil
}

// SearchOffers finds offers that match to foward conditions
func (mm *MatchMaker) SearchOffers(
	fconds FowardConditions) []TfcOffer {
	var matches []TfcOffer

	for _, offer := range mm.offers {
		if offer.Fconds() == fconds {
			matches = append(matches, *offer)
		}
	}

	return matches
}

// TakeOffer prepares DLC for a user and provides offer
func (mm *MatchMaker) TakeOffer(id int) (TfcOffer, error) {
	offer := mm.getOffer(id)

	var err error
	offer.dlc, err = offer.makeDlc(true, 1)

	if err != nil {
		return *offer, err
	}

	return *offer, nil
}

func (mm *MatchMaker) getOffer(id int) *TfcOffer {
	for _, o := range mm.offers {
		if o.ID() == id {
			return o
		}
	}
	return nil
}
