package matchmaker

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"

	"rpc"
	"tfc"
	"usr"
)

func TestMatchMakerOffers(t *testing.T) {
	// Prepare test data
	cparty, err := testCounterParty("Bob")
	if err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	settleAt := testSettleAt()
	fconds := tfc.NewFowardConditions(100, 0.1, 0.5, settleAt)
	offer := *NewTfcOffer(1, *cparty, fconds)

	mm := NewMatchMaker()

	if err = mm.PutOffer(offer); err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	offers := mm.Offers()

	if len(offers) > 1 {
		t.Errorf("offers shouldn't be empty")
		return
	}
	if offers[0].ID() != offer.ID() {
		t.Errorf("Invalid offers found")
		return
	}
}

func TestMatchMakerSearchOffers(t *testing.T) {
	// Prepare test data
	cparty, err := testCounterParty("Bob")
	if err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	settleAt := testSettleAt()
	fconds1 := tfc.NewFowardConditions(100, 0.1, 0.5, settleAt)
	offer1 := *NewTfcOffer(1, *cparty, fconds1)
	fconds2 := tfc.NewFowardConditions(200, 0.1, 0.5, settleAt)
	offer2 := *NewTfcOffer(2, *cparty, fconds2)

	mm := NewMatchMaker()

	for _, offer := range []TfcOffer{offer1, offer2} {
		if err = mm.PutOffer(offer); err != nil {
			t.Errorf("unexpecter error: %v", err)
			return
		}
	}

	offers := mm.SearchOffers(fconds1)

	if len(offers) != 1 {
		t.Errorf("Invalid number of offers matched. expected: %d, actual: %d", 1, len(offers))
		return
	}
	if offers[0].ID() != offer1.ID() {
		t.Errorf("Invalid offer found.")
		return
	}
}

func TestMatchMakerTakeOffer(t *testing.T) {
	// Prepare test data
	cparty, err := testCounterParty("Bob")
	if err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	settleAt := testSettleAt()
	fconds := tfc.NewFowardConditions(100, 0.1, 0.5, settleAt)
	offer := *NewTfcOffer(1, *cparty, fconds)

	mm := NewMatchMaker()

	if err = mm.PutOffer(offer); err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	var takenOffer TfcOffer
	takenOffer, err = mm.TakeOffer(offer.ID())

	if err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	dlc := takenOffer.Dlc()

	if &dlc == nil {
		t.Errorf("DLC should be prepared by matchmaker")
		return
	}

	dlcSettleAt := dlc.GameDate()
	if dlcSettleAt != settleAt {
		t.Errorf("Invalid settle at. expected: %s, actual: %s",
			settleAt, dlcSettleAt)
	}

	t.Logf("DLC made by matchmaker: %+v", dlc)
}

func testSettleAt() time.Time {
	n := time.Now()
	tomorrow := n.AddDate(0, 0, 1)
	return tomorrow
}

func testCounterParty(name string) (*usr.User, error) {
	chainParams := chaincfg.RegressionNetParams
	rpc := rpc.NewBtcRPC("http://localhost:18443", "user", "pass")
	return usr.NewUser(name, chainParams, rpc)
}
