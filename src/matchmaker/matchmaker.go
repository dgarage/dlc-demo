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

// PutTfcOffer puts a new TFC offer on board
func (mm *MatchMaker) PutTfcOffer(offer TfcOffer) error {
	offers := mm.Offers()
	mm.offers = append(offers, &offer)
	return nil
}

// SearchTfcOffer finds matches offers
func (mm *MatchMaker) SearchTfcOffer(fconds FowardConditions) ([]*TfcOffer, error) {
	var matches []*TfcOffer

	for _, offer := range mm.Offers() {
		if offer.FowardConditions == fconds {
			matches = append(matches, offer)
		}
	}

	return matches, nil
}
