package main

import (
	"os"
	"github.com/arjanvaneersel/synonyms/thesaurus"
	"bufio"
	"log"
	"fmt"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for " + word + ".", err)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}