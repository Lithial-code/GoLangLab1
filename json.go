package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Companies struct {
	Companies []Company `json:"companies"`
}
type Company struct {
	Index       int    `json:"index"`
	CompanyName string `json:"company"`
	StaffSize   int    `json:"staff_size"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}
type Interns struct {
	Interns []Intern `json:"interns"`
}
type Name struct {
	FirstName string `json:"first"`
	LastName  string `json:"last"`
}
type Intern struct {
	InternName Name   `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	CompanyID  int    `json:"company_id"`
}

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

//TODO
func printInternData(interns Interns, companies Companies) {
	internsWithoutJobs := 0
	companyLen := len(companies.Companies)
	internsLen := len(interns.Interns)

	for i := 0; i < internsLen; i++ {
		var intern = interns.Interns[i]
		var companyID = intern.CompanyID
		if intern.CompanyID > 0 && intern.CompanyID < companyLen {
			fmt.Printf("Intern Name: %s\n", interns.Interns[i].InternName.FirstName)
			fmt.Printf("Last Name: %s\n", interns.Interns[i].InternName.LastName)
			fmt.Printf("Company Name: %s\n", companies.Companies[companyID].CompanyName)
			fmt.Printf("Company Email: %s\n", companies.Companies[companyID].Email)
		}
	}
	fmt.Printf("Unclaimed Interns: \n")
	for j := 0; j < companyLen; j++ {
		var intern = interns.Interns[j]
		if intern.CompanyID < 0 || intern.CompanyID > companyLen {
			fmt.Printf("Intern Name: %s\n", interns.Interns[j].InternName.FirstName)
			fmt.Printf("Last Name: %s\n", interns.Interns[j].InternName.LastName)
			internsWithoutJobs++
		}
	}
	fmt.Printf("***********************************************\n")
	fmt.Printf("Number of companies %d", companyLen)
	fmt.Printf("***********************************************\n")
	fmt.Printf("Interns without jobs: %d", internsWithoutJobs)

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
