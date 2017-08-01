package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/arjanvaneersel/gockan/gockan"
	"github.com/vuleetu/goconfig/config"
)

type Harvester struct {
	cfg      *config.Config
	name     string
	src, dst gockan.Repository
}

func NewHarvester(cfg *config.Config, name string, src, dst gockan.Repository) *Harvester {
	harvester := Harvester{cfg: cfg, name: name, src: src, dst: repo}
	return &harvester
}

func (h *Harvester) Start() {
	interval, err := h.cfg.Int(h.name, "harvest_interval")
	if err != nil {
		interval = 24 * 60
	}
	wait, err := h.cfg.Int(h.name, "harvest_wait")
	if err != nil {
		wait = rand.Int() % interval
	}

	go func() {
		log.Printf("[%s] harvest scheduled in %d minutes", h.name, wait)
		time.Sleep(5 * time.Second)
		for {
			//start := time.Now().Nanosecond()
			log.Printf("[%s] mirroring", h.name)
			err := MirrorRepo(h.src, h.dst, h.name+".gob")
			if err != nil {
				log.Print(err)
			}
			//end := time.Now().Nanosecond()
			next := 60
			log.Printf("[%s] done mirroring. next due in %d minutes",
				h.name, next)
			time.Sleep(60 * time.Minute)
		}
	}()
}
