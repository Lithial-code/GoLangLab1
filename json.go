package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

//Companies  array of companies from json
type Companies struct {
	Companies []Company `json:"companies"`
}

//Company company structure for json
type Company struct {
	Index       int    `json:"index"`
	CompanyName string `json:"company"`
	StaffSize   int    `json:"staff_size"`
	Address     string `json:"address"`
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
	InternName   Name   `json:"name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	CompanyID    int    `json:"company_id"`
	CompanyValid bool
}

//Keeping a bunch of stuff in here for now because i haven't worked out how to move it
func main() {
	jsonFile, err := os.Open("InternsAtCompanies2.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened the file")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var companies Companies
	var interns Interns
	json.Unmarshal(byteValue, &companies)
	json.Unmarshal(byteValue, &interns)
	//printAllDetails(companies)
	//printCompanyDetails(companies)
	printInternData(interns, companies)
}

//function that does all teh heavy listing. Takes in the loaded json in the form of the interns and companies variables
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
		fmt.Printf("Intern Name: %s\n", interns.Interns[i].InternName.FirstName)
		fmt.Printf("Last Name: %s\n", interns.Interns[i].InternName.LastName)
		if intern.CompanyID > 0 && intern.CompanyID < companyLen {
			fmt.Printf("Company Name: %s\n", companies.Companies[companyID].CompanyName)
			fmt.Printf("Company Email: %s\n", companies.Companies[companyID].Email)
			//tell us this company has at least one intern
			companies.Companies[companyID].HasIntern = true
		}
	}
	//loop again to check for unclaimed interns
	//TODO maybe add this to a struct and do the same as how we do with companies
	fmt.Printf("\nUnclaimed Interns: \n")
	for j := 0; j < internsLen; j++ {
		var intern = interns.Interns[j]
		if intern.CompanyID < 0 || intern.CompanyID > companyLen {
			fmt.Printf("Intern Name: %s\n", interns.Interns[j].InternName.FirstName)
			fmt.Printf("Last Name: %s\n", interns.Interns[j].InternName.LastName)
			internsWithoutJobs++
		}
	}
	fmt.Printf("\nCompanies without Interns: \n")
	for k := 0; k < companyLen; k++ {
		var company = companies.Companies[k]
		if !company.HasIntern {
			fmt.Printf("Company: %s\n", company.CompanyName)
			companiesWithoutInterns++
		}
	}
	fmt.Printf("***********************************************\n")
	fmt.Printf("Companies without interns %d\n", companiesWithoutInterns)
	fmt.Printf("***********************************************\n")
	fmt.Printf("Interns without jobs: %d\n", internsWithoutJobs)

}
func printCompanyDetails(companies Companies) {
	for i := 0; i < len(companies.Companies); i++ {
		fmt.Println("Company ID: " + strconv.Itoa(i))
		fmt.Println("Company Name: " + companies.Companies[i].CompanyName)
	}
}
func printAllDetails(companies Companies) {
	for i := 0; i < len(companies.Companies); i++ {
		fmt.Println("Index: " + string(companies.Companies[i].Index))
		fmt.Println("Company: " + companies.Companies[i].CompanyName)
		fmt.Println("Staff Size: " + string(companies.Companies[i].StaffSize))
		fmt.Println("Address: " + companies.Companies[i].Address)
		fmt.Println("Phone: " + companies.Companies[i].Phone)
		fmt.Println("Email: " + companies.Companies[i].Email)
	}
}
