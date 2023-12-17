# Attendance Tracking System

## Overview

The Attendance Tracking System is a Go-based web application that allows users to log in and check in to timestamp their attendance.
Admins can upload a list of users through a .csv file, and users can register and log in using their assigned user ID, which is recorded in the .csv file.

## Features

- **User Authentication:** Users can register then log in using their unique user ID.
- **Admin Functionality:** Admins can upload a list of users through a .csv file.
- **Attendance Logging:** Users can check in to timestamp their attendance.
- **Attendance Reports:** Admins can view attendance records filtered by dates and export to a .csv file.

## Setup

Run the following commands to setup

```bash
git clone https://github.com/mahe993/attendance.git
```

```bash
go mod download
```

```bash
cp .env.example .env
```

## Usage

Replace values in `.env` with actual configuration values

To run the application, execute the following command:

```bash
cd <project-root-directory>/src
```

```bash
./attendance.exe
```

## Tech Spec

- User registration is limited to the IDs provided by the admin in the .csv file
- Data validation is enforced throughout the app:
  - User cannot register more than once
  - User cannot check in attendance more than once
  - User can only check in if on the appropriate WIFI
  - Admin can only upload .csv files with proper headers and data
  - If there are ID repeats in .csv uploads, the first/last names are modified only
- HTML injection is not possible through use of html/template package
- Passwords are handled with encryption
- .env files used for hiding sensitive data
- local database is maintained through JSON encoding/decoding
- Errors properly panics when needed and are logged (no outfile)
- Authenticated sessions are sent to the client through cookies
- Nested templates are used together with template functions and variables to provide a seamless browsing experience
- Codebase is divided mainly into three sections:
  - Router -- provides URL routing to specific controllers
  - Controllers -- breaks down URL by CRUD operations and routes to specific service
  - Services -- Handles the business logics of the endpoint
- This ensures modularity for easy scaling. Entire routes can also be easily protected at the router level
- Notable subsections includes:
  - Utility -- for util functions
  - Database functions -- for Read/Write operations to JSON

Go doc available at:

```bash
godoc -http=:6060
```
