package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
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

//Keeping a bunch of stuff in here for now because i haven't worked out how to move it
func main() {
	defer TimeTaken(time.Now(), "main")
	one := "data/InternsAtCompanies1.json"
	two := "data/InternsAtCompanies2.json"
	processData(one)
	processData(two)
}

//ProcessData this is a function that does all my main work so main is tidy
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
	printInternData(interns, companies) //this is stage ones function

	//printAllDetails(companies)   		//this was the first test
	//printCompanyDetails(companies)	//this was the second test
}

//TimeTaken  time taken to run a function
func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	fmt.Printf("TIME: %s took %s\n", name, elapsed)
	fmt.Printf("--------------------------------------------------\n")
}

//function that does all the heavy listing. Takes in the loaded json in the form of the interns and companies variables
func printInternData(interns Interns, companies Companies) {
	//used to count these things below
	internsWithoutJobs := 0
	companiesWithoutInterns := 0
	//so we don't have to count these more than once
	companyLen := len(companies.Companies)
	internsLen := len(interns.Interns)

	//loop the interns and print the relevant information.
	for i := 0; i < internsLen; i++ {
		var intern = interns.Interns[i]
		var companyID = intern.CompanyID
		//fmt.Printf("Intern Name: %s\n", interns.Interns[i].InternName.FirstName)
		//fmt.Printf("Last Name: %s\n", interns.Interns[i].InternName.LastName)
		if intern.CompanyID >= 0 && intern.CompanyID < companyLen {
			//fmt.Printf("Company Name: %s\n", companies.Companies[companyID].CompanyName)
			//fmt.Printf("Company Email: %s\n", companies.Companies[companyID].Email)
			//tell us this company has at least one intern
			companies.Companies[companyID].HasIntern = true
		}
	}
	//loop again to check for unclaimed interns
	//TODO maybe add this to a struct and do the same as how we do with companies
	fmt.Printf("Unclaimed Interns: \n")
	fmt.Printf("If not 0 its because they have no company associated with them in the dataset\n")
	for j := 0; j < internsLen; j++ {
		var intern = interns.Interns[j]
		if intern.CompanyID < 0 || intern.CompanyID > companyLen {
			fmt.Printf("Intern Name: %s %s\n", interns.Interns[j].InternName.FirstName, interns.Interns[j].InternName.LastName)
			//fmt.Printf("Last Name: %s\n", interns.Interns[j].InternName.LastName)
			internsWithoutJobs++
		}
	}
	fmt.Printf("\nCompanies without Interns: \n")
	fmt.Printf("If not 0 its because they had no interns associated with them in the dataset \n")

	for k := 0; k < companyLen; k++ {
		var company = companies.Companies[k]
		if !company.HasIntern {
			fmt.Printf("Company: %s\n", company.CompanyName)
			companiesWithoutInterns++
		}
	}
	fmt.Printf("--------------------------------------------------\n")
	fmt.Printf("Companies without interns %d\n", companiesWithoutInterns)
	fmt.Printf("--------------------------------------------------\n")
	fmt.Printf("Interns without jobs: %d\n", internsWithoutJobs)
	fmt.Printf("--------------------------------------------------\n")

}

//this loops through all the companies to print their ids and names. Used for debugging
func printCompanyDetails(companies Companies) {
	for i := 0; i < len(companies.Companies); i++ {
		fmt.Println("Company ID: " + strconv.Itoa(i))
		fmt.Println("Company Name: " + companies.Companies[i].CompanyName)
	}
}

//this loops through all the companies and lists all their details
func printAllDetails(companies Companies) {
	for i := 0; i < len(companies.Companies); i++ {
		fmt.Println("Index: " + string(rune(companies.Companies[i].Index)))
		fmt.Println("Company: " + companies.Companies[i].CompanyName)
		fmt.Println("Staff Size: " + string(rune(companies.Companies[i].StaffSize)))
		fmt.Println("Address: " + companies.Companies[i].Address)
		fmt.Println("Phone: " + companies.Companies[i].Phone)
		fmt.Println("Email: " + companies.Companies[i].Email)
	}
}
