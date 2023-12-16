package services

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/templates"
	utils "attendance.com/src/util"
)

type AdminPageVariables struct {
	User User
	Tab  string
}
type AdminService struct {
	Variables AdminPageVariables
}

var (
	Admin AdminService
)

func (p *AdminService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	p.Variables.User = currUser
	p.Variables.Tab = strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]

	err := templates.Tpl.ExecuteTemplate(w, "adminPage", p.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

// UploadStudentsList function is used to upload a CSV file containing a list of students.
// Errors on non-CSV files, saves a copy of the uploaded CSV file to /db/uploads,
// and updates the MapUsers state. DB users.json is then updated with the new MapUsers state.
func (p *AdminService) UploadStudentsList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("error updating users.json::" + err.(error).Error())
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}()

	file, fileInfo, err := r.FormFile("csvFile")
	if err != nil {
		logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Check if the file has a ".csv" extension.
	if !strings.HasSuffix(fileInfo.Filename, ".csv") {
		http.Error(w, "Invalid file format. Please upload a .csv file", http.StatusBadRequest)
		return
	}

	logger.Println(fmt.Sprint("\nFile:", file, "\nFile-HeaderProperties:", fileInfo, "\nerr: ", err))

	csvData, err := utils.ReadCSV(file)
	if err != nil {
		logger.Println(err)
		http.Error(w, "Error processing CSV file", http.StatusInternalServerError)
		return
	}

	// Create a new CSV file for saving the uploaded data.
	// The new file will be created in the uploads folder with the name
	// studentList_<timestamp>.csv
	// e.g. studentList_2021-08-01_12:00:00.csv
	saveCSV := utils.WriteCSV(fmt.Sprintf("%s/db/uploads/studentList_%s.csv", os.Getenv("APP_BASE_PATH"), time.Now().Format("2006-01-02_15:04:05")), csvData)

	// Update MapUsers with the uploaded student list
	headlessCSVData := csvData[1:]
	for _, line := range headlessCSVData {
		student := User{
			ID:    line[0],
			First: line[1],
			Last:  line[2],
		}

		// if the student already exists in MapUsers, update their name
		if user, ok := MapUsers[student.ID]; ok {
			user.First, user.Last = student.First, student.Last
			MapUsers[student.ID] = user
			continue
		}

		MapUsers[student.ID] = student
	}

	// Write MapUsers to database
	// Can potentially panic here if the database is not writable
	err = db.Write(MapUsers, "users.json")
	if err != nil {
		logger.Println(err)
	}

	// Wait for the CSV file to be saved
	<-saveCSV

	http.Redirect(w, r, "/admin/success", http.StatusFound)
}
