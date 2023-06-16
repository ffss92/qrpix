package main

import (
	"log"

	"github.com/ffss92/qrpix"
)

const (
	code = "00020126580014br.gov.bcb.pix0136123e4567-e12b-12d1-a456-4266554400005204000053039865802BR5913Fulano de Tal6008BRASILIA62070503***63041D3D"
)

func main() {
	p := qrpix.NewParser()
	static, err := p.ParseStatic(code)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", static)
}
