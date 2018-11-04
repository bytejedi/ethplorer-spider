package main

func main() {
	hub := newHub()
	putter := newPutter(hub)
	spider := newSpider(hub)
	go putter.run()
	go spider.run()
	hub.run()
}
