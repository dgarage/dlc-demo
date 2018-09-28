// scenario.go
package main

import (
	"fmt"
	"log"
	"strconv"
)

type scenario struct {
	memo  string
	steps []func(int, *Demo) error
	pos   int
}

func (s *scenario) step(d *Demo) error {
	if s.pos < 0 || len(s.steps) <= s.pos {
		fmt.Printf("This scenario is over.\n")
		return nil
	}
	// if s.pos == 0 {
	// 	fmt.Printf("This scenario start.\n")
	// }
	err := s.steps[s.pos](s.pos+1, d)
	if err != nil {
		return err
	}
	s.pos++
	if len(s.steps) == s.pos {
		fmt.Printf("This scenario finish.\n")
	}
	return nil
}

func set(args []string, d *Demo) error {
	var err error
	idx := 0
	if len(args) > 1 {
		idx, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	list := []func(*Demo) (*scenario, error){}
	list = append(list, scenario0)
	if idx < 0 || len(list) <= idx {
		return fmt.Errorf("out of range. %d,%d", idx, len(list))
	}
	err = faucet(nil, d)
	if err != nil {
		return err
	}
	d.sc, err = list[idx](d)
	if err != nil {
		return err
	}
	d.alice.ClearDlc()
	d.bob.ClearDlc()
	log.Printf("set the scenario.\n")
	log.Printf("%s\n", d.sc.memo)
	return nil
}

func step(args []string, d *Demo) error {
	if d.sc == nil {
		return fmt.Errorf("scenario is nil")
	}
	return d.sc.step(d)
}

//----------------------------------------------------------------

func scenario0(d *Demo) (*scenario, error) {
	sc := &scenario{}
	sc.memo = "normal"
	sc.steps = append(sc.steps, stepBobPutTfcOfferOnBoard)
	sc.steps = append(sc.steps, stepAliceSendOfferToBob)
	sc.steps = append(sc.steps, stepBobSendAcceptToAlice)
	sc.steps = append(sc.steps, stepAliceSendSignToBob)
	sc.steps = append(sc.steps, stepAliceAndBobSetOracleSign)
	sc.steps = append(sc.steps, stepAliceOrBobSendSettlementTx)
	return sc, nil
}

//----------------------------------------------------------------
