package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"regexp"
)

func correctGroup(originalGroupName string) string {
	//Change to lower case
	newGroupName := strings.ToLower(originalGroupName)
	

	//Remove invalid characters - mimincs terraform
	reg := regexp.MustCompile(`[^0-9A-Za-z\\-]`)
	strs := reg.ReplaceAllString(newGroupName, "")

	//Check that the team name is a valid OCP / Kubernetes team name
	match, err := regexp.MatchString(`[a-z0-9]([-a-z0-9]*[a-z0-9])?`, strs)
	if err != nil {
		fmt.Println("Error ", err)
	}
	if !match{
		strs = "invalid-team-name"
	}

	//Dump out the new team name in the terminal. Sort out debugging :-)
	fmt.Println(strs)
	return strs
}

func addRoleBinding(username string, group string) {



	fmt.Printf("Attempting to create role binding for user %s in project %s\n",username, group)
	//out, err := exec.Command("oc","adm","policy","add-role-to-user","admin",username,"-n" ,group).Output()
	args := []string{"adm", "policy", "add-role-to-user", "admin", username, "-n", group}

	cmd := exec.Command("oc", args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal("cmd.Run() failed with $s\n", err)
	}
	fmt.Printf("combined out: \n%s\n", string(out))
}
func createUser(user string){
	args := []string{"create","user", user}

	cmd := exec.Command("oc", args...)
	out, _ := cmd.CombinedOutput()

//		if err != nil {
//			log.Fatal("cmd.Run() failed with $s\n", err)
//		}
	fmt.Printf("combined out: \n%s\n", string(out))
}
func createGroupProjects(group string) {

	//correctedGroup := correctGroup(group)
	correctedGroup := group
	//out, err := exec.Command("oc","adm","policy","add-role-to-user","admin",username,"-n" ,group).Output()
		args := []string{"new-project", correctedGroup}

		cmd := exec.Command("oc", args...)
		out, _ := cmd.CombinedOutput()
	
//		if err != nil {
//			log.Fatal("cmd.Run() failed with $s\n", err)
//		}
		fmt.Printf("combined out: \n%s\n", string(out))

}

func getUniqueGroups() {
}

func main() {

	if runtime.GOOS == "windows" {
		fmt.Println("Not supported on Windows")
		os.Exit(1)
	}

	// Open the CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Print the CSV data
	for _, row := range data {
		for _, col := range row {
			fmt.Printf("%s,", col)
		}
		fmt.Println()
	}



	//Get unique group names
	getUniqueGroups()

	for _, row := range data {
		//Create the user
		createUser(row[0])
		
		group := correctGroup(row[1])
		
		//Create the Goroup project
		createGroupProjects(group)
		
		//Create the role binding
		addRoleBinding(row[0], group)
	}

}
