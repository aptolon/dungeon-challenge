package main

import (
	"bufio"
	"dungeon-challenge/internal/config"
	"dungeon-challenge/internal/event"
	"dungeon-challenge/internal/processor"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	eventsFile, err := os.Open("events")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer eventsFile.Close()

	p := processor.NewProcessor(cfg)

	scanner := bufio.NewScanner(eventsFile)

	for scanner.Scan() {
		line := scanner.Text()

		ev, err := event.ParseEvent(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err := p.Process(ev); err != nil {
			fmt.Println(err)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}
	p.PrintReport()

}
