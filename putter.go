package main

import (
	"fmt"
	"github.com/ethplorer-spider/config"
	"log"
)

type Putter struct {
	hub    *Hub
	stopCh chan struct{}
}

func newPutter(h *Hub) *Putter {
	return &Putter{
		hub:    h,
		stopCh: make(chan struct{}),
	}
}

func (p *Putter) run() {
	p.hub.wg.Add(1)
	defer p.hub.wg.Done()
	db := config.DB()
	defer db.Close()

	for {
		select {
		case tx := <-p.hub.TxCh:
			if tx == nil {
				log.Println("Putter stopped.")
				return
			}
			_, err := db.Exec(
				"INSERT INTO `tx`(`tx`,`timestamp`,`from`,`to`,`value`,`t`) values(?,?,?,?,?,?)",
				tx.Tx,
				tx.Timestamp,
				tx.From,
				tx.To,
				tx.Value,
				tx.T,
			)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(fmt.Sprintf("Inserted TxHash: %s", tx.Tx))
		}
	}
}
