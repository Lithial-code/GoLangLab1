package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//Companies array of companies from json
type Companies struct {
	Companies []Company `json:"companies"`
}

//Company company structure for json
type Company struct {
	Index       int    `json:"index"`
	CompanyName string `json:"company"`
	StaffSize   int    `json:"staff_size"`
	Address     string `json:"address"`
	Suburb      string
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	HasIntern   bool
}

//Interns array of Interns for json
type Interns struct {
	Interns []Intern `json:"interns"`
}

//Name names structure for interns in json
type Name struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
}

//Intern Intern structure for json load in
type Intern struct {
	InternName Name   `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	CompanyID  int    `json:"company_id"`
}

//Suburb struct for adding proper suburb information to the program
type Suburb struct {
	Name        string
	Count       int
	StaffTotal  int
	AverageSize float64
}

//Suburbs array used to hold my collection of suburbs
var Suburbs []Suburb

//details on a simple local postgres location
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "school"
)

//Using comments to change my data sets.
//TODO implement basic console UI if i get time at the end
func main() {
	//defer TimeTaken(time.Now(), "main")
	one := "data/InternsAtCompanies1.json"
	//two := "data/InternsAtCompanies2.json"
	//processData(two)
	db := createConnection()
	defer db.Close()
	createTable()

	processData(one)

}

//TODO clean up old comments and such
//ProcessData this is a function that does all my main work so main is tidy
//Processes one json file
func processData(location string) {
	defer TimeTaken(time.Now(), location)
	jsonFile, err := os.Open(location)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully opened %s\n", location)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var companies Companies
	var interns Interns
	json.Unmarshal(byteValue, &companies)
	json.Unmarshal(byteValue, &interns)
	postCompaniesToDatabase(companies)
	postInternsToDatabase(interns)
	megaSplitQuery()
}

// this is the query for part 3
func megaSplitQuery() {
	defer TimeTaken(time.Now(), "Mega Query")
	sqlStatement := `SELECT split_part(address, ',', 2) As "Suburbs",
	COUNT( DISTINCT companyname) AS "CompaniesArea",
	SUM(staffsize) AS "StaffTotal",
	ROUND(AVG(staffsize),2) AS "Average"
	FROM companies
	GROUP BY "Suburbs"
	ORDER BY "Average" DESC`
	query(sqlStatement, "mega query")
}

//function for posting the interns to the database
func postInternsToDatabase(interns Interns) {
	internlen := len(interns.Interns)
	//this was to remember what i named everything
	/*		internID INT SERIAL,
			phone TEXT,
			email TEXT,
			companyID INT,
			PRIMARY KEY(internID),
			FOREIGN KEY(companyID)
			REFERENCES companies(companyID)`*/
	sqlstatement := `INSERT INTO interns (name,phone, email, companyID) 
					 VALUES ($1, $2, $3,$4) 
					 RETURNING internID
					 `
	db := createConnection()
	//prepares to close database when done
	defer db.Close()
	var id int64
	//execute the sql statement and return a response
	for i := 0; i < internlen; i++ {
		companyID := 0
		if interns.Interns[i].CompanyID != 0 {
			companyID = interns.Interns[i].CompanyID
		}
		err := db.QueryRow(sqlstatement,
			interns.Interns[i].InternName.FirstName+" "+interns.Interns[i].InternName.LastName,
			interns.Interns[i].Phone,
			interns.Interns[i].Email,
			companyID,
		).Scan(&id)
		if err != nil {
			//logged things like this because i couldn't work out a way to do it differently.
			//2 interns are missing from the data because the company doesnt exist and i didnt want to break my tables for them
			log.Printf("%s | %d | This company doesn't exist so we're going to skip it", interns.Interns[i].InternName.FirstName, companyID)
			log.Printf("Unable to execute insert into interns. %v", err) //this could be logged as fatal but it would kill all the rest of the inserts
		}
	}
}

//posts the companies to the database
func postCompaniesToDatabase(companies Companies) {
	companyLen := len(companies.Companies)

	/*companyID INT,
	companyname TEXT,
	staffsize int,
	address TEXT,
	suburb TEXT,
	phone TEXT,
	email TEXT,
	hasintern BOOL,
	PRIMARY KEY(companyID)*/
	sqlstatement := `INSERT INTO companies (companyid, companyname, staffsize, address,suburb, phone, email) VALUES ($1, $2, $3,$4,$5,$6,$7) RETURNING companyid`
	db := createConnection()
	//prepares to close database when done
	defer db.Close()
	var id int64
	//execute the sql statement and return a response
	for i := 0; i < companyLen; i++ {
		//suburbify - stolen from stage 2
		array := strings.SplitAfter(companies.Companies[i].Address, ",")
		suburb := array[len(array)-1]
		err := db.QueryRow(sqlstatement,
			companies.Companies[i].Index,
			companies.Companies[i].CompanyName,
			companies.Companies[i].StaffSize,
			companies.Companies[i].Address,
			suburb,
			companies.Companies[i].Phone,
			companies.Companies[i].Email).Scan(&id)
		if err != nil {
			log.Fatalf("Unable to execute insert into companies. %v", err)
		}
	}
}

//TimeTaken  time taken to run a function
//Used with the defer feature in go to start a timer and return the value at the end of the block
func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	fmt.Printf("TIME: %s took %s\n", name, elapsed)
	fmt.Printf("--------------------------------------------------\n")
}

//CreateConnection create connection with postgres db
func createConnection() *sql.DB {
	// load .env file

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// check the connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	// return the connection
	return db
}

//CreateTable creates a default table in the database
func createTable() {
	//got sick of breaking my databases manually
	sqlStatement :=
		`DROP TABLE IF EXISTS interns;
		 DROP TABLE IF EXISTS companies;
		`
	executeQuery(sqlStatement, "drop tables")
	/*	Index       int    `json:"index"`
		CompanyName string `json:"company"`
		StaffSize   int    `json:"staff_size"`
		Address     string `json:"address"`
		Suburb      string
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		HasIntern   bool*/
	//create the base notes table for if it doesn't exist
	sqlStatement = `CREATE TABLE IF NOT EXISTS companies (
		companyID INT,
		companyname TEXT,
		staffsize int,
		address TEXT,
		suburb TEXT,
		phone TEXT,
		email TEXT,
		PRIMARY KEY(companyID)
	);`
	executeQuery(sqlStatement, "create companies")
	/*	InternName Name   `json:"name"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
		CompanyID  int    `json:"company_id"`*/
	//create the base notes table for if it doesn't exist
	sqlStatement = `CREATE TABLE IF NOT EXISTS interns (
			internID INT GENERATED ALWAYS AS IDENTITY,
			name TEXT,
			phone TEXT,
			email TEXT,
			companyID INT,
			PRIMARY KEY(internID),
			FOREIGN KEY(companyID) 
			REFERENCES companies(companyID)
		);`
	executeQuery(sqlStatement, "create interns")

}

//easy way to do database executions. probably bad tho
func executeQuery(sqlStatement string, qID string) {
	db := createConnection()
	//prepares to close database when done
	defer db.Close()
	//execute the sql statement and return a response
	res, err := db.Exec(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute %s | . %v", qID, err)
	}
	//print the response maybe
	fmt.Printf("%s\n ", res)
}

//the grunt work of the mega query
func query(sqlStatement string, qID string) {
	db := createConnection()
	//prepares to close database when done
	defer db.Close()
	//execute the sql statement and return a response
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute %s | . %v", qID, err)
	}
	var suburbs []Suburb
	//dont forget to close the rows
	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var suburb Suburb

		// unmarshal the row object to user
		err = rows.Scan(&suburb.Name, &suburb.Count, &suburb.StaffTotal, &suburb.AverageSize)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		// append the user in the users slice
		suburbs = append(suburbs, suburb)

	}
	//print the things
	for _, item := range suburbs {
		fmt.Printf("%s : %.2f\n", item.Name, item.AverageSize)
	}
}
