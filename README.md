# Reminder-GO Application

## Description
The Reminder-GO application is a handy task management tool. Using the application, you can easily add, view, delete, and change the status of reminds, customize and change them, and manage the sending of email notifications by profile and for each remind separately. The application use PostgreSQL to store all reminds

## Who can benefit
The application can be useful for anyone who needs a convenient task management tool, including entrepreneurs, students, housewives, etc.

## Features: 
- Adding and removing reminds
- Remind status change: when the status changes to "done", the "when completed" date is automatically updated
- View all reminds
- View completed tasks: the ability to view a list of completed reminds in a certain date range with the ability to sort by the date when they should be completed 
- View current tasks: give a complete list of current tasks with the ability to sort by the date when they should be completed 
- Ability to edit remind - change its title, description, deadline, and notification time
- Ability to add your profile notification config. By agreeing to receive notifications by email, you can set how many days before the remind deadline you will receive a letter notification in your mail. This works across your profile and all your reminds
- Ability to configure notifications for each remind at your discretion. You can set the notification time for the number of minutes/hours/days you want before the reminder deadline. It works on a similar principle as in Google Calendar

## **Reminder-GO local installation**
_Pre-requisites: GIT, Docker, Docker Compose, Golang v1.19 and newer, MinGW or Cygwin(for Windows users)_

Steps:
1. Install Docker and Docker Compose manually
2. Install Golang v1.19 and newer
3. Install git and clone the Reminder-GO repository
4. If you are using Windows install MinGW or Cygwin. You need to add the path to the "make" utility to the PATH environment variable to the ability to run the "make" utility from the command line.
5. Install all dependencies using the command `go mod download`

 
> **Note**: Before running the application locally you need an `.env` file which you can get from the admins of Reminder-GO repository. You need to add this file to the root folder. Also, you need serviceAccountKey.json with credentials for Firebase authentication.

All commands and their descriptions are described in `Makefile`

## Reminder-GO local start
You can launch all applications using one single command `docker-compose up`

Or run the next commands step-by-step
#### 1. Create and run the PostgreSQL DB

`make createdb`

`make migrateup`

`make db-run`

It was launched DB `reminder` on port `:5432` with two tables `todo` and `users_config`

#### 2. Launch Reminder and Notification worker

`make run`

`make run-worker`

Now Application is running on port `:8000` and ready to accept requests

## Running tests
- To launch application tests you need to run commands:

`make mocks`

`make create_testdb`

`make migrateup_test`

`make test`

`make int_test`

- To get test coverage of the application you need to run command - `make coverage` or in HTML - `make coverage-html`

## Swagger documentation

To generate actual swagger documentation run `make swag_gen`

## Reminder app structure

It's a Restful API CRUD application with routes"

- `/remind` - [method GET] - get a list of reminds by the query("all", "current", "completed"). Also required params for pagination and date range

- `/remind` - [method POST] - create new remind

- `/remind/${id}` - [method GET] - get remind by ID

- `/remind/${id}` - [method DELETE] - delete remind by ID

- `/remind/${id}` - [method PUT] - update remind by ID

- `/status/${id}` - [method PUT] - change remind status

- `/configs/${id}