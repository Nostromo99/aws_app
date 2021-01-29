package main

import(
	"html/template"
	"net/http"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"os"


)
type welcome struct{
	Value string
}
func main(){
	// ticker:=info{"0.0"}
	port := os.Getenv("PORT")
    if port == "" {
	port = "5000"
    }
	http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		t,_:=template.ParseFiles("public/index.html")
		var value string
		if r.FormValue("request")==""{
			
			value=""
		}else{
		value=request(r.FormValue("request"))
	}
		x:=welcome{Value:value}
		t.Execute(w,x)
		
	
	})
    http.ListenAndServe(":"+port, nil)
}
func request(stock string) (string) {
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
	re:=regexp.MustCompile(`Bp4i AP7Wnd">[0-9]`)
	position:=re.FindStringIndex(outputstring)
	if position==nil{
		return "invalid input"
	}
	position2:=position[1]
	value := string(outputstring[position2-1 : position2+25])
	endpoint := strings.Index(value, " ")
	value = value[0:endpoint]
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return "invalid input"
	}

	return stock+":"+value+changeval
}