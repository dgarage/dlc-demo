package matchmaker

import (
	"testing"
	"time"

	"github.com/btcsuite/btcd/chaincfg"

	"rpc"
	"usr"
)

func TestMatchMaker(t *testing.T) {
	mm := NewMatchMaker()
	offer, err := demoTfcOffer()

	if err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	if err = mm.PutTfcOffer(offer); err != nil {
		t.Errorf("unexpecter error: %v", err)
		return
	}

	offers := mm.Offers()

	if len(offers) < 1 {
		t.Errorf("offers shouldn't be empty")
		return
	}
}

func demoTfcOffer() (TfcOffer, error) {
	var offer TfcOffer
	chainParams := chaincfg.RegressionNetParams
	rpc := rpc.NewBtcRPC("http://localhost:18443", "user", "pass")
	cparty, err := usr.NewUser("Bob", chainParams, rpc)
	if err != nil {
		return offer, err
	}
	fconds := NewFowardConditions(100, 0.01, time.Now())
	offer = NewTfcOffer(111, *cparty, fconds)
	return offer, nil
}
