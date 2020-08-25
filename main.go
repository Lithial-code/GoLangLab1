package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
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
	jsonFile, err := os.Open("data/InternsAtCompanies2.json")
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
	//printAllDetails(companies)   			//this was the first test
	//printCompanyDetails(companies)		//this was the second test
	//printInternData(interns, companies) //this is stage ones function
	printCompanySuburbData(companies)
}

//TimeTaken  time taken to run a function
func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	fmt.Printf("TIME: %s took %s\n", name, elapsed)
}
func printCompanySuburbData(companies Companies) {
	companyLen := len(companies.Companies)
	mapOfSuburbs := make(map[string][2]int)

	for i := 0; i < companyLen; i++ {
		array := strings.SplitAfter(companies.Companies[i].Address, ",")
		suburb := array[len(array)-1]
		//fmt.Printf("Company %s \n", suburb)
		companies.Companies[i].Suburb = suburb
	}
	//Sort by suburb, keeping original order or equal elements.
	sort.SliceStable(companies.Companies, func(i, j int) bool {
		return companies.Companies[i].Suburb < companies.Companies[j].Suburb
	})
	// for i := 0; i < companyLen; i++ {
	// 	fmt.Printf("Company %s \n", companies.Companies[i].Suburb)
	// }

	//counter := 0
	// average := 0
	// suburb := ""

	for i := 0; i < companyLen; i++ {
		applicableSuburb := companies.Companies[i].Suburb

		if val, ok := mapOfSuburbs[applicableSuburb]; !ok {
			//fmt.Printf("value: %s\n", applicableSuburb)
			val[0] += companies.Companies[i].StaffSize
			val[1]++
			mapOfSuburbs[applicableSuburb] = val

			//fmt.Printf("Suburb: %s | %d| %d\n", applicableSuburb, val[0], val[1])

		}
		// if suburb == "" || suburb != applicableSuburb {
		// 	if suburb != "" {
		// 		fmt.Printf("Suburb %s has an average of %d per company\n", suburb, average/counter)
		// 	}
		// 	suburb = applicableSuburb
		// 	counter = 1
		// 	average += companies.Companies[i].StaffSize
		// } else {
		// 	counter++
		// 	average += companies.Companies[i].StaffSize
		// }
	}
	values := make([]int, 0, len(mapOfSuburbs))
	for v := range mapOfSuburbs {
		temp := int(v[0] / v[1])
		values = append(values, temp)
	}
	sort.Ints(values)

	for k, v := range mapOfSuburbs {
		fmt.Printf("Suburb: %s has an average of %d staff company\n", k, v[0]/v[1])
	}

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
		fmt.Printf("Intern Name: %s\n", interns.Interns[i].InternName.FirstName)
		fmt.Printf("Last Name: %s\n", interns.Interns[i].InternName.LastName)
		if intern.CompanyID >= 0 && intern.CompanyID < companyLen {
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
		fmt.Println("Index: " + string(companies.Companies[i].Index))
		fmt.Println("Company: " + companies.Companies[i].CompanyName)
		fmt.Println("Staff Size: " + string(companies.Companies[i].StaffSize))
		fmt.Println("Address: " + companies.Companies[i].Address)
		fmt.Println("Phone: " + companies.Companies[i].Phone)
		fmt.Println("Email: " + companies.Companies[i].Email)
	}
}
