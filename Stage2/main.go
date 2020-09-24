package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
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

type Suburb struct {
	Name        string
	Companies   []Company
	AverageSize int
}

var Suburbs []Suburb

//Keeping a bunch of stuff in here for now because i haven't worked out how to move it
func main() {
	defer TimeTaken(time.Now(), "main")
	// one := "data/InternsAtCompanies1.json"
	two := "data/InternsAtCompanies2.json"
	// processData(one)
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
	//printInternData(interns, companies) //this is stage ones function
	printCompanySuburbData(companies)

	//printAllDetails(companies)   		//this was the first test
	//printCompanyDetails(companies)	//this was the second test
}

//TimeTaken  time taken to run a function
func TimeTaken(t time.Time, name string) {
	elapsed := time.Since(t)
	fmt.Printf("TIME: %s took %s\n", name, elapsed)
	fmt.Printf("--------------------------------------------------\n")
}

func printCompanySuburbData(companies Companies) {
	companyLen := len(companies.Companies)

	//create suburbs
	for i := 0; i < companyLen; i++ {
		array := strings.SplitAfter(companies.Companies[i].Address, ",")
		suburb := array[len(array)-1]
		if i == 0 {
			var newSuburb Suburb
			newSuburb.Name = suburb
			Suburbs = append(Suburbs, newSuburb)
		} else {
			found := false
			for _, sub := range Suburbs {
				if suburb == sub.Name {
					found = true
					break
				}
			}
			if !found {
				var newSuburb Suburb
				newSuburb.Name = suburb
				Suburbs = append(Suburbs, newSuburb)
			}
		}
	}
	suburbsLen := len(Suburbs)
	//this one adds the companies
	for j := 0; j < companyLen; j++ {
		array := strings.SplitAfter(companies.Companies[j].Address, ",")
		suburb := array[len(array)-1]

		for k := 0; k < suburbsLen; k++ {
			if Suburbs[k].Name == suburb {
				Suburbs[k].Companies = append(Suburbs[k].Companies, companies.Companies[j])
			}
		}
	}
	//this one averages the staff
	for i := 0; i < suburbsLen; i++ {
		numOfCompanies := len(Suburbs[i].Companies)

		// for _, company := range Suburbs[i].Companies {
		// 	fmt.Println(company.CompanyName)
		// }
		for j := 0; j < numOfCompanies; j++ {
			if Suburbs[i].Companies[j].StaffSize == 0 {
				fmt.Println("staff size divide by zero immenient")
			}
			Suburbs[i].AverageSize += Suburbs[i].Companies[j].StaffSize
		}
		if Suburbs[i].AverageSize == 0 {
			fmt.Printf("suburb effect %s\n", Suburbs[i].Name)
			for _, company := range Suburbs[i].Companies {
				fmt.Println(company.CompanyName)
			}
			//fmt.Println("average size divide by zero immenient")
		}
	}
	//sort the suburbs
	//I used this sorting function first. After getting it to work i moved on to quick sort
	// sort.Slice(Suburbs, func(i, j int) bool {
	// 	return Suburbs[i].AverageSize > Suburbs[j].AverageSize
	// })
	//Better sorting, now with quick sort
	Suburbs = quicksort(Suburbs)
	fmt.Println("This is a list of all the suburbs sorted by average staff size from high to low")
	for _, v := range Suburbs {
		fmt.Printf("%s : %d\n", v.Name, v.AverageSize)
	}

}

//Quick sort function. Takes in a suburb array and spits it back out sorted.
//is approx 30% faster at sorting the datasets we have
func quicksort(a []Suburb) []Suburb {
	if len(a) < 2 {
		return a
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]

	for i := range a {
		if a[i].AverageSize > a[right].AverageSize {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]

	quicksort(a[:left])
	quicksort(a[left+1:])

	return a
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
	for j := 0; j < internsLen; j++ {
		var intern = interns.Interns[j]
		if intern.CompanyID < 0 || intern.CompanyID > companyLen {
			fmt.Printf("Intern Name: %s %s\n", interns.Interns[j].InternName.FirstName, interns.Interns[j].InternName.LastName)
			//fmt.Printf("Last Name: %s\n", interns.Interns[j].InternName.LastName)
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
