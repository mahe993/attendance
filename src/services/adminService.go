package services

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	if currUser != nil {
		p.Variables.User = *currUser
	} else {
		p.Variables.User = User{}
	}

	p.Variables.Tab = strings.Split(r.URL.Path, "/")[len(strings.Split(r.URL.Path, "/"))-1]

	err := templates.Tpl.ExecuteTemplate(w, "adminPage", p.Variables)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal("Template execution error:", err)
		return
	}
}

func (p *AdminService) UploadStudentsList(w http.ResponseWriter, r *http.Request) {
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

	destFile, err := os.Create(fmt.Sprintf("./db/uploads/studentList_%s.csv", time.Now().Format("2006-01-02_15:04:05")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer destFile.Close()

	// Create a CSV writer.
	csvWriter := csv.NewWriter(destFile)
	defer csvWriter.Flush()

	// Write each row of the CSV data to the output file.
	for _, line := range csvData {
		if err := csvWriter.Write(line); err != nil {
			logger.Println(fmt.Sprint("Error writing to CSV file:", err))
			http.Error(w, "Error processing CSV file", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/admin/success", http.StatusFound)
}
