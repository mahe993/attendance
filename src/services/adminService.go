package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	utils "attendance.com/src/util"
)

type OverviewFilters struct {
	DateFrom string
	DateTo   string
}

type AdminPageVariables struct {
	User    states.User
	Tab     string
	Filters OverviewFilters
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
	dateFrom, dateTo := r.FormValue("dateFrom"), r.FormValue("dateTo")

	if p.Variables.Tab == "overview" &&
		(dateFrom == "" || dateTo == "") {
		today := time.Now().Format("2006-01-02")
		http.Redirect(w, r, fmt.Sprintf("/admin/overview?dateFrom=%s&dateTo=%s", today, today), http.StatusFound)
		return
	}

	p.Variables.Filters = OverviewFilters{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	err := templates.Tpl.ExecuteTemplate(w, "adminPage", p.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

// UploadStudentsList function is used to upload a CSV file containing a list of students.
// Errors on non-CSV files, saves a copy of the uploaded CSV file to /db/uploads,
// and updates the states.MapUsers state. DB users.json is then updated with the new states.MapUsers state.
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

	// Update states.MapUsers with the uploaded student list
	headlessCSVData := csvData[1:]
	for _, line := range headlessCSVData {
		student := states.User{
			ID:    line[0],
			First: line[1],
			Last:  line[2],
		}

		// if the student already exists in states.MapUsers, update their name
		if user, ok := states.MapUsers[student.ID]; ok {
			user.First, user.Last = student.First, student.Last
			states.MapUsers[student.ID] = user
			continue
		}

		states.MapUsers[student.ID] = student
	}

	// Write MapUsers state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.MapUsers, "users.json")
	if err != nil {
		logger.Println(err)
	}

	// Wait for the CSV file to be saved
	<-saveCSV

	http.Redirect(w, r, "/admin/success", http.StatusFound)
}

func (p *AdminService) ExportAttendanceCSV(w http.ResponseWriter, r *http.Request) {
	dateFrom, dateTo := r.FormValue("dateFrom"), r.FormValue("dateTo")
	checkedInUsers := templates.GetCheckedInUsers(dateFrom, dateTo)

	if dateFrom == "" || dateTo == "" {
		http.Error(w, "Error exporting CSV, check to ensure a valid date range is selected.", http.StatusBadRequest)
		return
	}
	dateFromTime, err := time.ParseInLocation("2006-01-02", dateFrom, time.Now().Location())
	if err != nil {
		http.Error(w, "Error exporting CSV, check to ensure a valid date range is selected.", http.StatusBadRequest)
		return
	}
	dateToTime, err := time.ParseInLocation("2006-01-02", dateTo, time.Now().Location())
	if err != nil {
		http.Error(w, "Error exporting CSV, check to ensure a valid date range is selected.", http.StatusBadRequest)
		return
	}
	if dateFromTime.After(dateToTime) {
		http.Error(w, "Error exporting CSV, check to ensure a valid date range is selected.", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("Attendance_%s_TO_%s.csv", dateFrom, dateTo)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "text/csv")

	// parse checkedInUsers into a [][]string csv data
	csvData := [][]string{{"Date", "ID", "Name", "Check-In Time"}}
	for dateFromTime.Before(dateToTime) || dateFromTime.Equal(dateToTime) {
		k := dateFromTime.Format("2006-01-02")
		if users, ok := checkedInUsers[k]; ok {
			if users != nil {
				for id, details := range users {
					csvData = append(csvData, []string{k, id, details.Name, details.CheckInTime})
				}
			} else {
				csvData = append(csvData, []string{k, "-", "-", "-"})
			}
		}
		dateFromTime = dateFromTime.AddDate(0, 0, 1)
	}

	writeCsv := utils.WriteCSV("../temp/"+fileName, csvData)
	<-writeCsv

	// Open and read the CSV file
	file, err := os.Open("../temp/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the file to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file.Close()

	// Delete the CSV file from the temp folder
	err = os.Remove("../temp/" + fileName)
	if err != nil {
		logger.Println(err)
	}

}
