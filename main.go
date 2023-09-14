package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
)

// Person represents a person's data.
type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	persons     = make(map[int]Person)
	currentID   = 1
	personsLock sync.RWMutex
	csvFile     = "persons.csv"
)

func main() {

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			getPerson(w, r)
		case http.MethodPost:
			createPerson(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deletePerson(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPatch {
			updatePerson(w, r)
		} else if r.Method == http.MethodGet {
			getPerson(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server on port 8080
	port := "8080" // Change this to your desired port
	fmt.Printf("Server listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	personsLock.RLock()
	defer personsLock.RUnlock()

	// Parse the URL to get the person's ID
	idStr := r.URL.Path[len("/api/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	person, found := persons[id]
	if !found {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	// Return the person's information as JSON
	jsonResponse, err := json.Marshal(person)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	personsLock.Lock()
	defer personsLock.Unlock()
	loadDataFromCSV()

	// Parse the JSON payload
	var newPerson Person
	err := json.NewDecoder(r.Body).Decode(&newPerson)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate the input (Name should only contain valid characters)
	if !isValidName(newPerson.Name) {
		http.Error(w, "Invalid Name", http.StatusBadRequest)
		return
	}

	// Assign a new ID to the person
	newPerson.ID = currentID
	currentID++

	// Save the person's data
	persons[newPerson.ID] = newPerson

	// Check if the CSV file exists, and create it if not
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		err := createCSVFile()
		if err != nil {
			http.Error(w, "Failed to create CSV file", http.StatusInternalServerError)
			return
		}
	}

	// Save the person to the CSV file
	err = savePersonToCSV(newPerson)
	if err != nil {
		http.Error(w, "Failed to save data to CSV", http.StatusInternalServerError)
		return
	}

	// Return the person's information as JSON
	jsonResponse, err := json.Marshal(newPerson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func createCSVFile() error {
	// Create a new CSV file and initialize it with headers
	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := []string{"ID", "Name"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	return nil
}

func updatePerson(w http.ResponseWriter, r *http.Request) {
	personsLock.Lock()
	defer personsLock.Unlock()
	loadDataFromCSV()

	// Parse the URL to get the person's ID
	idStr := r.URL.Path[len("/api/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Check if the person exists
	person, found := persons[id]
	if !found {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	// Parse the JSON payload
	var updatedPerson Person
	err = json.NewDecoder(r.Body).Decode(&updatedPerson)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate the input (Name should only contain valid characters)
	if !isValidName(updatedPerson.Name) {
		http.Error(w, "Invalid Name", http.StatusBadRequest)
		return
	}

	// Update the person's data
	person.Name = updatedPerson.Name
	persons[id] = person

	// Update the person in the CSV file
	err = updatePersonInCSV(person)
	if err != nil {
		http.Error(w, "Failed to update data in CSV", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Person updated successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	personsLock.Lock()
	defer personsLock.Unlock()

	// Parse the URL to get the person's ID
	idStr := r.URL.Path[len("/api/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Check if the person exists
	_, found := persons[id]
	if !found {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	// Delete the person
	delete(persons, id)

	// Delete the person from the CSV file
	err = deletePersonFromCSV(id)
	if err != nil {
		http.Error(w, "Failed to delete data from CSV", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Person deleted successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

func savePersonToCSV(person Person) error {
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{strconv.Itoa(person.ID), person.Name}
	err = writer.Write(record)
	if err != nil {
		return err
	}

	return nil
}

func updatePersonInCSV(person Person) error {
	lines, err := readCSVLines()
	if err != nil {
		return err
	}

	// Find the index of the person's record in the CSV data
	index := -1
	for i, line := range lines {
		if len(line) > 0 && line[0] == strconv.Itoa(person.ID) {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("person not found in CSV")
	}

	lines[index] = []string{strconv.Itoa(person.ID), person.Name}

	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		err = writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func deletePersonFromCSV(id int) error {
	lines, err := readCSVLines()
	if err != nil {
		return err
	}

	// Filter out the person's record from the CSV data
	filteredLines := make([][]string, 0)
	for _, line := range lines {
		if len(line) > 0 && line[0] != strconv.Itoa(id) {
			filteredLines = append(filteredLines, line)
		}
	}

	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, line := range filteredLines {
		err = writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func readCSVLines() ([][]string, error) {
	file, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func loadDataFromCSV() {
	lines, err := readCSVLines()
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
		return
	}

	for _, line := range lines {
		if len(line) >= 2 {
			id, _ := strconv.Atoi(line[0])
			persons[id] = Person{
				ID:   id,
				Name: line[1],
			}
			if id >= currentID {
				currentID = id + 1
			}
		}
	}
}

// isValidName checks if a name contains only valid characters (letters and spaces).
func isValidName(name string) bool {
	// Use regular expression to check if the name contains only letters and spaces.
	validNamePattern := "^[A-Za-z ]+$"
	match, err := regexp.MatchString(validNamePattern, name)
	if err != nil {
		return false
	}
	return match
}
