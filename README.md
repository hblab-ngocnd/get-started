[![tests](https://github.com/hblab-ngocnd/get-started/actions/workflows/tests.yml/badge.svg)](https://github.com/hblab-ngocnd/get-started/actions/workflows/tests.yml)
[![Deploy to IBM Cloud Foundry](https://github.com/hblab-ngocnd/get-started/actions/workflows/dev.yml/badge.svg)](https://github.com/hblab-ngocnd/get-started/actions/workflows/dev.yml)
## Prerequisites

You'll need the following:
* [Git](https://git-scm.com/downloads)
* [Go](https://golang.org/dl/)

## 1. Clone the app

Now you're ready to start working with the simple Go *hello world* app. Clone the repository and change to the directory where the sample app is located.
  ```
git clone https://github.com/hblab-ngocnd/get-started
cd get-started
  ```

Peruse the files in the *get-started-go* directory to familiarize yourself with the contents.

## 2. Run the app locally use [air](https://github.com/cosmtrek/air)
Create .env file
```cmd
CLOUDANT_URL=dasdsa
GOOGLE_APPLICATION_API_KEY=sss
DEBUG=true
SYNC_PASS=123456
```
Build and run the app.
  ```
air
  ```

View your app at: http://localhost:8080