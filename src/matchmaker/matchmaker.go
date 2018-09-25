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
func (mm *MatchMaker) Offers() []*TfcOffer {
	return mm.offers
}

// PutOffer puts a new TFC offer on board
func (mm *MatchMaker) PutOffer(offer *TfcOffer) error {
	offers := mm.Offers()
	mm.offers = append(offers, offer)
	return nil
}

// SearchOffer finds offers that match to foward conditions
func (mm *MatchMaker) SearchOffers(
	fconds FowardConditions) []*TfcOffer {
	var matches []*TfcOffer

	for _, offer := range mm.Offers() {
		if offer.Fconds() == fconds {
			matches = append(matches, offer)
		}
	}

	return matches
}
