package main

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type fooHandler struct {
	Message string
}

type foo struct {
	Message string `json:"message,omitempty"`
	Name    string `json:"firstName,omitempty"`
	SurName string `json:"lastName,omitempty"`
	Age     int    `json:"age,omitempty"`
}

type Product struct {
	ProductId      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	Sku            string `json:"sku"`
	Upc            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product

func init() {

	productsJSON := `[
		{
			"productId": 1,
			"manufacturer": "Johns-Jenkins",
			"sku": "p5z343vdS",
			"upc": "939581000000",
			"pricePerUnit": "497.45",
			"quantityOnHand": 9703,
			"productName": "sticky note"
		},
		{
			"productId": 2,
			"manufacturer": "Hessel, Schimmel and Feeny",
			"sku": "i7v300kmx",
			"upc": "740979000000",
			"pricePerUnit": "282.29",
			"quantityOnHand": 9217,
			"productName": "leg warmers"
		},
		{
			"productId": 3,
			"manufacturer": "Swaniawski, Bartoletti and Bruen",
			"sku": "q0L657ys7",
			"upc": "111173000000",
			"pricePerUnit": "436.26",
			"quantityOnHand": 5905,
			"productName": "lamp shade"
		}		
	]`

	if err := json.Unmarshal([]byte(productsJSON), &productList); err != nil {
		log.Fatal(err)
	}
}
func getNextID() int {
	highestId := -1
	for _, product := range productList {
		if highestId < product.ProductId {
			highestId = product.ProductId
		}
	}
	return highestId + 1
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		if err = json.Unmarshal(bs, &newProduct); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}

		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		newProduct.ProductId = getNextID()
		productList = append(productList, newProduct)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(f.Message))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is from the bar handler function...\n"))
	fs := foo{Message: "Hello from Bar struct", Name: "Chris", SurName: "Scogin", Age: 32}

	bs, err := json.Marshal(&fs)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(bs)

	f := foo{}

	if err := json.Unmarshal(bs, &f); err != nil {
		fmt.Println(err)
	}

	fmt.Println(f.Message, "My name is", f.Name, f.SurName, "and I am", f.Age, "years old...")
	fmt.Println(string(bs))
}

func main() {

	fh := fooHandler{
		Message: "Hello from Go!",
	}

	http.HandleFunc("/products", productsHandler)
	http.Handle("/foo", &fh)
	http.HandleFunc("/bar", barHandler)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Println(err)
	}

}
