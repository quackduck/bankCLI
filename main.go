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
	AmountCents int    `json:"amount_cents"`
	Memo        string `json:"memo"`
	Date        string `json:"date"`
	Type        string `json:"type"`
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

	var org struct {
		PublicMsg string `json:"public_message"`
	}
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
