package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"attendance.com/src/db"
	"attendance.com/src/logger"
	"attendance.com/src/states"
	"attendance.com/src/templates"
	utils "attendance.com/src/util"
)

// OverviewFilters struct represents the filters used in the overview page
type OverviewFilters struct {
	DateFrom string
	DateTo   string
}

// AdminPageVariables struct represents the variables that are passed to the admin page template
type AdminPageVariables struct {
	User    states.User
	Tab     string
	Filters OverviewFilters
}

// AdminService struct provides methods for handling business logics for requests to the /admin endpoint
type AdminService struct {
	Variables   AdminPageVariables
	VariablesMu sync.Mutex
	ExportMu    sync.Mutex
}

var (
	// Admin is a global variable that provides access to the AdminService methods
	Admin AdminService
)

// Index handles the HTTP request to the admin index page.
// It retrieves the current user, date filters, and tab information from the request.
// If the tab is "overview" and the date filters are not provided, it redirects to the overview page with today's date.
// It renders the admin page template with the provided variables.
func (p *AdminService) Index(w http.ResponseWriter, r *http.Request) {
	currUser := Auth.GetUser(r)
	dateFrom, dateTo := r.FormValue("dateFrom"), r.FormValue("dateTo")

	// Mutex lock to ensure thread-safe access to shared Variables field
	p.VariablesMu.Lock()
	defer p.VariablesMu.Unlock()
	p.Variables.User = currUser
	p.Variables.Tab = strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]

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

// UploadStudentsList handles the HTTP request to upload a CSV file containing a list of students.
// It checks if the uploaded file is a CSV file, saves a copy of the file, and updates the user database.
// The CSV file is also data validated to ensure it has the correct format.
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

	if len(csvData[0]) != 3 {
		http.Error(w, "Invalid CSV file format. Please ensure the CSV file has only 3 columns (ID, First Name, Last Name)", http.StatusBadRequest)
		return
	}

	if csvData[0][0] != "ID" || csvData[0][1] != "First" || csvData[0][2] != "Last" {
		http.Error(w, "Invalid CSV file format. Please ensure the CSV file has header of 3 columns (ID, First, Last)", http.StatusBadRequest)
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
		if user, ok := states.GetMapUser(student.ID); ok {
			user.First, user.Last = student.First, student.Last
			states.SetMapUser(student.ID, user)
			continue
		}

		states.SetMapUser(student.ID, student)
	}

	// Wait for the CSV file to be saved
	<-saveCSV

	// Write MapUsers state to database
	// Can potentially panic here if the database is not writable
	err = db.Write(states.GetAllMapUsers(), "users.json")
	if err != nil {
		logger.Println(err)
	}

	http.Redirect(w, r, "/admin/success", http.StatusFound)
}

// ExportAttendanceCSV handles the HTTP request to export attendance data as a CSV file.
// It retrieves the date filters from the request and generates the CSV data.
// It locks the export process to ensure thread-safe access to the CSV file.
// It writes the CSV data to a temporary file, reads the file, and copies it to the response writer.
// Finally, it deletes the temporary file.
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

	// Mutex lock to ensure thread-safe access to returning CSV file to user.
	// Without lock, it is potentially unsafe since fileName could be the same for multiple requests.
	// The CSV file could be deleted before the response is returned to the user
	p.ExportMu.Lock()
	defer p.ExportMu.Unlock()
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
