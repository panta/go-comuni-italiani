package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	comuni "github.com/panta/go-comuni-italiani"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("can't read file: %v", err)
	}

	comuniSlice := []comuni.Comune{}
	if err := json.Unmarshal(data, &comuniSlice); err != nil {
		log.Fatalf("can't unmarshal JSON: %v", err)
	}

	for _, comune := range comuniSlice {
		fmt.Printf("%30s  ISTAT:%s %-30s %-4s %-3s\n",
			comune.Denominazione, comune.CodiceComune,
			comune.DenominazioneRegione, comune.CodiceProvincia, comune.SiglaAutomobilistica)
	}
}
