package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//struct for changing page based on input
type welcome struct {
	Value string
}

func main() {
	//main sets up the page and changes it when someone submits to it
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("public/index.html")
		var value string
		if r.FormValue("request") == "" {

			value = ""
		} else {
			value = request(r.FormValue("request"))
		}
		x := welcome{Value: value}
		t.Execute(w, x)

	})
	http.ListenAndServe(":"+port, nil)
}

//request simply requests stock info from google searches
//the bulk of the function is error checking and regex to extract information
func request(stock string) string {
	resp, err := http.Get("https://www.google.com/search?q=" + stock + "+stock&safe=strict&rlz=1C1CHBF_enIE784IE784&sxsrf=ALeKk00R7a6xf6UTpwxc_R23lq_m9yQx0A:1600457550056&gbv=1&sei=TgtlX_eLA9uM1fAP7MacsAc")
	if err != nil {
		fmt.Println("get request for ticker failed")
		return "invalid input"
	}
	defer resp.Body.Close()
	output, _ := ioutil.ReadAll(resp.Body)

	outputstring := string(output)
	changepos := strings.LastIndex(outputstring, "lB8g7")
	changeval := ""
	if changepos == -1 {
		changepos = strings.LastIndex(outputstring, "AWuZUe")
		changeval = string(outputstring[changepos+8 : changepos+25])
	} else {

		changeval = string(outputstring[changepos+7 : changepos+25])
	}
	lastpos := strings.LastIndex(changeval, "<")
	changeval = changeval[0:lastpos]
	re := regexp.MustCompile(`Bp4i AP7Wnd">[0-9]`)
	position := re.FindStringIndex(outputstring)
	if position == nil {
		return "invalid input"
	}
	position2 := position[1]
	value := string(outputstring[position2-1 : position2+25])
	endpoint := strings.Index(value, " ")
	value = value[0:endpoint]
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return "invalid input"
	}
	write_to_db(stock, value)
	return stock + ":" + value + changeval
}

//establishes connection with db and writes information to it
func write_to_db(stock string, price string) {
	timeNow := time.Now()
	sqlTime := timeNow.Format("2006-01-02 15:04:05")
	priceFloat, _ := strconv.ParseFloat(price, 64)
	db, err := sql.Open() //sql credentials go here
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO Stocks (ticker,price,time) VALUES(?,?,?);", strings.ToUpper(stock), priceFloat, sqlTime)
	if err != nil {
		panic(err.Error())
	}
}
