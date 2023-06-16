package main

import (
	"fmt"
	"log"

	"github.com/ffss92/qrpix"
)

func main() {
	static := qrpix.NewStatic(
		"123e4567-e12b-12d1-a456-426655440000",
		"Fulano de Tal",
		"BRASILIA",
		"***",
	)

	brCode, err := static.BRCode()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generated code:", brCode)

	if err := static.SaveFile("_examples/file/pix.png"); err != nil {
		log.Fatal(err)
	}
}
