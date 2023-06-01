
# Reminder-GO Application

## Description
The Reminder-GO application is a handy task management tool. Using the application, you can easily add, view, delete and change the status of reminds, customize and change them, manage the sending of email notifications by profile and for each remind separately. Application use PostgreSQL to store all reminds

## Who can benefit
The application can be useful for anyone who needs a convenient task management tool, including entrepreneurs, students, housewives, etc.

## Features: 
- Adding and removing reminds
- Remind status change: when the status changes to "done", the "when completed" date is automatically updated
- View all reminds
- View completed tasks: the ability to view a list of completed reminds in a certain date range with the ability to sort by the date when they should be completed 
- View current tasks: give a complete list of current tasks with the ability to sort by the date when they should be completed 
- Ability to edit remind - change its title, description, deadline and notification time
- Ability to add your profile notification config. By agreeing to receive notifications by email, you can set how many days before the remind deadline you will receive a letter-notification in your mail. This works across your profile and across all your reminds
- Ability to configure notifications for each remind at your discretion. You can set the notification time for the number of minutes/hours/days you want before the reminder deadline. It works on a similar principle as in Google Calendar

## **Reminder-GO local installation**
_Pre-requisites: GIT, Docker, Docker Compose, Golang v1.19 and newer, MinGW or Cygwin(for Windows users)_

Steps:
1. Install Docker and Docker Compose manually
2. Install Golang v1.19 and newer
3. Install git and clone the Reminder-GO repository
4. If you are using Windows install MinGW or Cygwin. You need to add the path to the "make" utility to the PATH environment variable to ability to run the "make" utility from the command line.
5. Install all dependencies using command `go mod download`

 
> **Note**: Before running the application locally you need an `.env` file which you can get from the admins of Reminder-GO repository. You need add this file in the root folder. Also you need serviceAccountKey.json with credentials for firebase authentication.

All commands and their descriptions are described in `Makefile`

## Reminder-GO local start
You can launch all aplication using one single command `docker-compose up`

Or run next commands step-by-step
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
- To launch application tests you need run commands:

`make mocks`

`make create_testdb`

`make migrateup_test`

`make test`

`make int_test`

- To get test coverage of the application you need run command - `make coverage` or in HTML - `make coverage-html`

## Swagger documentation

To generate actual swagger documentation run `make swag_gen`

## Reminder app structure

It's Restful API CRUD application with routes"

- `/remind` - [method GET] - get list of reminds by query("all", "current", "completed"). Also required params for pagination and date range

- `/remind` - [method POST] - create new remind

- `/remind/${id}` - [method GET] - get remind by ID

- `/remind/${id}` - [method DELETE] - delete remind by ID

- `/remind/${id}` - [method PUT] - update remind by ID

- `/status/${id}` - [method PUT] - change remind status

- `/configs/${id}` - [method GET] - get user notification configs 

- `/configs/${id}` - [method PUT] - update user notification configs

Remidner use Firebase for authentication

You need to pass the verification token in each request. This token is checked in the `AuthMiddleware` which verifies it via Firebase Auth Client which is initialized with credentials from `serviceAccountKey.json` in the root folder 

## Notification worker  structure
This is a service that starts and runs in a goroutine. Every 5 seconds, the service goes through the database and looks for a reminder to send a notification via the SMTP protocol

## CI/CD
**Note:** currently CI/CD is implemented for stage environment only and steps for coa infrastructure (with frontend and backend services) are identical, except coa repo doesn’t have tests

![image](https://github.com/red-rocket-software/reminder-go/assets/73254444/f2667561-4cb7-4cf8-bfc3-a5168a1a099b)


Pipeline is triggered either by push or merge (pull request) but certain jobs work with dev branch only. If make feature or fix branch from dev branch, pipeline will only run lint and tests. For dev branch after tests and lint finish successfully build job will start. Build job builds Docker images and pushes them to GCR. After succesfull build job manual approval for deployment is required (again, dev branch only). Pipeline creates issue for specified approvers. One of the approvers has to use key word specified in issue in order to deploy both services into GKE cluster or deny deployment.

**Note:** notice two GitHub Actions secrets: GCP_SA_KEY and INFRASTRUCTURE_PRIVATE_KEY. The first one contains json key to access GCP. In case you want to update json key, make sure you remove all new lines and pass value as a single string since it may cause issues while reading the secret. INFRASTRUCTURE_PRIVATE_KEY is used to access coa-infrastructure repository. Pipeline requires Helm charts in order to run *helm upgrade* command so you have to access coa-infrastructure repository somehow. If for some reason you want to update INFRASTRUCTURE_PRIVATE_KEY, generate new public/private key pair and set private key as value for INFRASTRUCTURE_PRIVATE_KEY. Public key has to be added in coa-infrastructure repository.
