package main

import (
	"ethtx-spider/config"
	"log"
	"sync"
	"time"
)

type Hub struct {
	wg             *sync.WaitGroup
	ResponseBodyCh chan *ResponseBody
	TxCh           chan *config.TxTable
	sp             *startPoint
	StopCh         chan struct{}
}

type startPoint struct {
	Ts     int64
	TxHash string
}

func newHub() *Hub {
	// 查询数据库中最后的一条数据，从这条数据开始爬取
	db := config.DB()
	defer db.Close()

	var sp startPoint
	err := db.QueryRow("SELECT `tx`,`timestamp` FROM `tx` ORDER BY `timestamp` DESC LIMIT 1").Scan(&sp.TxHash, &sp.Ts)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &Hub{
		wg:             new(sync.WaitGroup),
		ResponseBodyCh: make(chan *ResponseBody, 1000),
		TxCh:           make(chan *config.TxTable, 1000),
		sp:             &sp,
		StopCh:         make(chan struct{}),
	}
}

func (h *Hub) run() {
	var canExit bool
	for {
		select {
		case body := <-h.ResponseBodyCh:
			for _, t := range body.Transfers {
				if t.Timestamp >= h.sp.Ts && t.TransactionHash != h.sp.TxHash {
					select {
					case h.TxCh <- &config.TxTable{
						Tx:        t.TransactionHash,
						Timestamp: t.Timestamp,
						From:      t.From,
						To:        t.To,
						Value:     t.Value,
						T:         time.Unix(t.Timestamp, 0).Format("2006-01-02 15:04:05"),
					}:

					}
				} else if !canExit {
					canExit = true
					close(h.StopCh)
				}
			}
			if canExit {
				log.Println("Hub stopped.")
				close(h.TxCh)
				h.wg.Wait()
				return
			}
		}
	}
}
