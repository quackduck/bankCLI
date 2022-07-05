package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	chalk "github.com/jwalton/gchalk"
)

type Transaction struct {
	//ID             string `json:"id"`
	//Object         string `json:"object"`
	//Href           string `json:"href"`
	AmountCents int    `json:"amount_cents"`
	Memo        string `json:"memo"`
	Date        string `json:"date"`
	//OrganizationID string `json:"organization_id"`
	Type string `json:"type"`
	//Pending        bool   `json:"pending"`
	//AchTransfer    struct {
	//	ID          string `json:"id"`
	//	Object      string `json:"object"`
	//	Href        string `json:"href"`
	//	Memo        string `json:"memo"`
	//	Transaction struct {
	//	} `json:"transaction"`
	//	AmountCents string `json:"amount_cents"`
	//	Date        string `json:"date"`
	//	Status      string `json:"status"`
	//	Beneficiary struct {
	//		Name string `json:"name"`
	//	} `json:"beneficiary"`
	//} `json:"ach_transfer"`
	//Check struct {
	//	ID          string `json:"id"`
	//	Object      string `json:"object"`
	//	Href        string `json:"href"`
	//	Memo        string `json:"memo"`
	//	Transaction struct {
	//	} `json:"transaction"`
	//	AmountCents int    `json:"amount_cents"`
	//	Date        string `json:"date"`
	//	Status      string `json:"status"`
	//} `json:"check"`
	//Donation struct {
	//	ID          string `json:"id"`
	//	Object      string `json:"object"`
	//	Href        string `json:"href"`
	//	Memo        string `json:"memo"`
	//	Transaction struct {
	//	} `json:"transaction"`
	//	AmountCents int `json:"amount_cents"`
	//	Donor       struct {
	//		Name string `json:"name"`
	//	} `json:"donor"`
	//	Date   string `json:"date"`
	//	Status string `json:"status"`
	//} `json:"donation"`
	//Invoice struct {
	//	ID          string `json:"id"`
	//	Object      string `json:"object"`
	//	Href        string `json:"href"`
	//	Memo        string `json:"memo"`
	//	Transaction struct {
	//	} `json:"transaction"`
	//	AmountCents string `json:"amount_cents"`
	//	Sponsor     struct {
	//		ID   string `json:"id"`
	//		Name string `json:"name"`
	//	} `json:"sponsor"`
	//	Date   string `json:"date"`
	//	Status string `json:"status"`
	//} `json:"invoice"`
	//Transfer struct {
	//	ID          string `json:"id"`
	//	Object      string `json:"object"`
	//	Href        string `json:"href"`
	//	Memo        string `json:"memo"`
	//	Transaction struct {
	//	} `json:"transaction"`
	//	AmountCents string `json:"amount_cents"`
	//	Date        string `json:"date"`
	//	Status      string `json:"status"`
	//} `json:"transfer"`
}

// to string for transaction
func (t Transaction) String() string {
	base := fmt.Sprintf("%10.2f\t%s\t%s\033[K", float64(t.AmountCents)/100.0, t.Date, t.Memo)
	if t.AmountCents < 0 {
		return chalk.BgHex("47000F")(base)
		//return chalk.BgHex("9C0021")(base)
		//return chalk.BgBrightRed(base)
	} else {
		return chalk.BgHex("094700")(base)
		//return chalk.BgHex("0C5900")(base)
		//return chalk.BgBrightGreen(base)
	}
}

type Org struct {
	PublicMsg string `json:"public_message"`
}

//type Check struct {
//	ID          string      `json:"id"`
//	Object      string      `json:"object"`
//	Href        string      `json:"href"`
//	Memo        string      `json:"memo"`
//	Transaction Transaction `json:"transaction"`
//	AmountCents int         `json:"amount_cents"`
//	Date        time.Time   `json:"date"`
//	Status      string      `json:"status"`
//}
//
//func (c Check) String() string {
//	return fmt.Sprintf("{%s, %0.2f, %s}", c.Memo, float64(c.AmountCents)/100.0, c.Date.Format("January 2 2006"))
//}

var (
	lastRequest []Transaction
	reader      = bufio.NewReader(os.Stdin)
)

func main() {
	// get id from the user
	fmt.Print("Enter organization id: ")
	var id string
	fmt.Scanln(&id)

	getProfile(id)

	fmt.Println(`Menu:
1. Get last n transactions
2. Filter last request
3. Exit`)

	for {
		choice := input(chalk.Cyan("$ "))

		switch choice {
		case "1":
			getTransactions(id)
		case "2":
			filterTransactions()
		case "3":
			fmt.Println("Bai!")
			return
		}
	}

}

func getProfile(id string) {
	url := fmt.Sprintf("https://bank.hackclub.com/api/v3/organizations/" + id)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var org Org
	err = json.Unmarshal(body, &org)
	if err != nil {
		panic(err)
	}
	bef, _, _ := strings.Cut(org.PublicMsg, "\n")
	fmt.Println(bef)
	//fmt.Print(string(markdown.Render(org.PublicMsg, 100, 0)))
}

func getTransactions(id string) {

	// number of transactions to get
	var n int
	fmt.Print("Enter number of transactions to get: ")
	fmt.Scanln(&n)

	if n < 1 || n > 500 {
		fmt.Println("Invalid number of transactions: must be between 1 and 500 inclusive.")
		return
	}

	// get the transactions
	url := fmt.Sprintf("https://bank.hackclub.com/api/v3/organizations/" + id + "/transactions?per_page=" + fmt.Sprintf("%d", n))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// unmarshal the response body
	var transactions []Transaction
	err = json.Unmarshal(body, &transactions)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// save the transactions
	lastRequest = transactions

	if input("Print transactions? (y/n): ") == "y" {
		for _, transaction := range transactions {
			fmt.Println(transaction)
		}
	}
}

func filterTransactions() {
	// ask for filter term
	term := input("Enter filter term: ")

	for _, transaction := range lastRequest {
		if strings.Contains(transaction.String(), term) {
			//chalk.Hex("DC143C")(term)
			fun := chalk.Hex("FF5C5C")
			//fun := chalk.Hex("DC143C")
			fmt.Println(strings.ReplaceAll(transaction.String(), term, fun(term)))
		}
	}
}

func input(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
