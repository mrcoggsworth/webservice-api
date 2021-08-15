package product

import (
	"fmt"
	"log"
	"sync"
)

var productMap struct {
	sync.RWMutex
	m map[int]Product
}

func init() {
	fmt.Println("loading products...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d Products loaded...\n", len(productMap.m))
}

func loadProductMap()(map[int]Product, error){
	return map[int]Product{}, nil
}