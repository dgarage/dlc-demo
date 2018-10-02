// demo project main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"

	"matchmaker"
	"oracle"
	"rpc"
	"usr"
)

func main() {
	logfile, err := os.OpenFile("./demo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic("cannot open demo.log:" + err.Error())
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	// log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags + log.Lshortfile)
	// init
	demo, err := initial()
	if err != nil {
		log.Fatalf("initial error : %+v\n", err)
		return
	}
	err = set([]string{"set", "0"}, demo)
	if err != nil {
		log.Fatalf("set error : %+v\n", err)
		return
	}
	console(demo)
}

// Demo is dataset for demo
type Demo struct {
	mm     *matchmaker.MatchMaker
	rpc    *rpc.BtcRPC
	alice  *usr.User
	bob    *usr.User
	olivia *oracle.Oracle
	sc     *scenario
}

func initial() (*Demo, error) {
	s := time.Now()
	log.Printf("begin initial\n")
	d := &Demo{}
	// TODO bitcoin rpc of regtest
	d.rpc = rpc.NewBtcRPC("http://localhost:18443", "user", "pass")
	d.mm = matchmaker.NewMatchMaker()

	// regtest requires 432 blocks to make csv active
	res, err := d.rpc.Request("getblockcount")
	if err != nil {
		return nil, err
	}
	height, _ := res.Result.(float64)
	log.Printf("block count  : %.0f\n", height)
	if height < 432 {
		_, err = d.rpc.Request("generate", 432-height)
		if err != nil {
			return nil, err
		}
	}

	// demo balance
	res, err = d.rpc.Request("getbalance")
	if err != nil {
		return nil, err
	}
	total, _ := res.Result.(float64)
	log.Printf("total amount : %.8f BTC\n", total)

	params := chaincfg.RegressionNetParams
	// Olivia (Oracle)
	d.olivia, err = oracle.NewOracle("Olivia", params, d.rpc)
	if err != nil {
		return nil, err
	}
	// Alice (User)
	d.alice, err = usr.NewUser("Alice", params, d.rpc)
	if err != nil {
		return nil, err
	}
	// Bob (User)
	d.bob, err = usr.NewUser("Bob", params, d.rpc)
	if err != nil {
		return nil, err
	}
	log.Printf("end   initial %f sec\n", (time.Now()).Sub(s).Seconds())
	return d, nil
}

func console(demo *Demo) {
	cmds := listCmds()
	fmt.Print("$ ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			args := strings.Split(line, " ")
			if args[0] == "exit" || args[0] == "e" {
				fmt.Println("bye")
				break
			}
			flg := false
			for _, cmd := range cmds {
				flg = contains(cmd.n, args[0])
				if flg {
					err := cmd.f(args, demo)
					if err != nil {
						fmt.Printf("%s error : %+v\n", cmd.n[0], err)
					}
					break
				}
			}
			if !flg {
				fmt.Printf("unknown command. %v\n", args)
			}
		}
		fmt.Print("$ ")
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
