package main

import (
	"ethtx-spider/config"
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"time"
)

type transfer struct {
	Timestamp       int64    `json:"timestamp"`
	TransactionHash string   `json:"transactionHash"`
	BlockNumber     int64    `json:"-"`
	Contract        string   `json:"contract"`
	Value           string   `json:"value"`
	IntValue        int64    `json:"-"`
	Type            string   `json:"type"`
	IsEth           bool     `json:"isEth"`
	Priority        int64    `json:"-"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	Addresses       []string `json:"-"`
	UsdPrice        string   `json:"-"`
}

type ResponseBody struct {
	Transfers []transfer `json:"transfers"`
}

type Spider struct {
	hub    *Hub
	stopCh chan struct{}
}

func newSpider(h *Hub) *Spider {
	return &Spider{
		hub:    h,
		stopCh: make(chan struct{}),
	}
}

func (s *Spider) run() {
	s.hub.wg.Add(1)
	defer s.hub.wg.Done()

	var page = 1
	for {
		ro := &grequests.RequestOptions{
			Params: map[string]string{
				"refresh": "transfers",
				"data":    "0xF8C595D070d104377f58715ce2E6C93E49a87f3c",
				"page":    "tab=tab-transfers&holders=2&showTx=tokens&transfers=%d&pageSize=100",
				"showTx":  "tokens",
			},
		}
		ro.Params["page"] = fmt.Sprintf(ro.Params["page"], page)
		resp, err := grequests.Get(config.BaseUrl, ro)

		if err != nil {
			log.Println("Unable to make request: ", err)
			time.Sleep(time.Second * 30)
			continue
		}
		if resp.StatusCode != 200 {
			log.Println(fmt.Sprintf("Http StatusCode: %d", resp.StatusCode))
			time.Sleep(time.Second * 30)
			continue
		}

		var body ResponseBody
		err = resp.JSON(&body)
		if err != nil {
			log.Println("JSON unmarshal failed: ", err)
			time.Sleep(time.Second * 30)
			continue
		}
		if &body == nil {
			log.Println("Body is nil ")
			time.Sleep(time.Second * 30)
			continue
		}

		select {
		case <-s.hub.StopCh:
			log.Println("Spider stopped.")
			close(s.hub.ResponseBodyCh)
			return
		default:
		}

		select {
		case <-s.hub.StopCh:
			log.Println("Spider stopped.")
			close(s.hub.ResponseBodyCh)
			return
		case s.hub.ResponseBodyCh <- &body:
			log.Println(fmt.Sprintf("Page %d was sent to ResponseBodyCh", page))
			page++
		}
		time.Sleep(time.Second * 2)
	}
}
