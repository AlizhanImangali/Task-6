package main

import (
	"database/sql"
	"encoding/json"

	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Response struct {
	XMLName xml.Name `xml:"rates"`
	Title   string   `xml:"title"`
	Date    string   `xml:"date"`

	Items []Item `xml:"item"`
}

type Item struct {
	XMLName  xml.Name `xml:"item"`
	Fullname string   `xml:"fullname"`
	// id     int     `xml:id`
	Title       string `xml:"title"`
	Description string `xml:"description"`

	// value  float32 `xml:"description,attr"`
}
type Rows struct {
	Count int `json:"count"`
}

type CodeAndDate struct {
	Data []string `json:"data"`
}

func main() {

	router := mux.NewRouter()
	//router.HandleFunc("/currency/save", getCurrecnyByDate).Methods("GET")
	router.HandleFunc("/currency/save", getCurrencyByDateAndCode).Methods("GET")
	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	//fmt.Println(result.Date)
	//sp_GetRates("cba", "1996-05-13")
}
func getCurrencyByDateAndCode(response http.ResponseWriter, request *http.Request) {
	date := request.FormValue("date")
	code := request.FormValue("code")

	layout := "02.01.2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		log.Printf("Parse problem: %v", err)
	}
	n_date := string(t.Format("2006.01.02"))
	sp_GetRates(n_date, code)

	u := &CodeAndDate{
		Data: []string{n_date, code}}
	r, _ := json.Marshal(u)
	fmt.Println(string(r) + ",")
}

/*func getCurrecnyByDate(response http.ResponseWriter, request *http.Request) {
	fmt.Println("request: ", request)
	date := request.FormValue("date")
	fmt.Println("date: ", date)
	query := "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date

	xmlBytes, err := getBytes(query)
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	var result Response
	xml.Unmarshal(xmlBytes, &result)

	rowsAdded, err := AddCurrency(&result)
	count := Rows{Count: rowsAdded}
	json.NewEncoder(response).Encode(count)
	fmt.Println("Rows Added: ", rowsAdded)
}*/

func getBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Read body: %v", err)
	}
	return data, err
}
func DB() *sql.DB {
	connStr := "user=postgres password=1234 dbname=Test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected")
	//defer db.Close()
	return db
}

func sp_GetRates(date string, code string) {
	var db = DB()
	rows, err := db.Query("Select a_date ,code from r_currency Where a_date= $1 and code=$2", date, code)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err := rows.Scan(&date, &code)
		if err != nil {
			panic(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(date, code) /*
		var cd []CodeAndDate
		j, _ := json.Marshal(cd)
		log.Println(j)*/
}
func AddCurrency(rates *Response) (rowsAdded int, err error) {
	db := DB()
	var result sql.Result
	rowsAdded = 0
	for _, rate := range rates.Items {
		/*layout := "02-01-2006"
		input := rates.Date
		t, _ := time.Parse(layout, input)*/
		layout := "02.01.2006"
		t, err := time.Parse(layout, rates.Date)
		result, err = db.Exec("INSERT INTO r_currency(title,code,value, a_date) Values($1,$2,$3,$4)", rate.Fullname, rate.Title, rate.Description, t)
		if err != nil {
			panic(err)
		}
		rowsAdded++
	}
	defer db.Close()
	if result != nil {
		fmt.Println(result)
	}

	return
}

/*func sp_GetRates(a_date string) {
	var db = DB()
	rows, err := db.Query("Select * from r_currency where a_date=$1 ", a_date)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &title, &code, &value, &a_date)
		if err != nil {
			panic(err)
		}
		log.Println(id, title, code, value, a_date)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
} */
/*func (example Example) Select(code int , a_date string){
	rows , err := example.db.Exec("call sp_GetRates(? , ? )", , code , a_date)
	if err != nil {
		return
	}
	else{

	}
}*/
