# Digital Diary
A cloud native digital diary that allows users to add journal entries in their calendar built using Golang for API, backed by Postgres database for storage. 

The application consists of Golang backend API which interacts with the frontend C# Blazor WASM user interface.

## What is Digital Diary and why use it?
Nowadays to write down notes or information you would use a book. With the increased usage of internet, everything around us is getting digitalized. Digital Diary offers the same. A pocket diary you can use from any device and anywhere in the world with an internet connection.

With Digital Diary:

* Keep your thoughts organized.
* Set and achieve your goals.
* Record your ideas.
* Allow yourself to self reflect.
* In corporate settings, use it to take meeting notes in meetings.

And many other different usecases!

## How to run the application?
1. To run the application, make sure docker and docker-compose is installed.
2. git clone this repository
3. To run both the frontend and backend applications, in root of the app directory, run the command: `docker-compose up -d`

4. In browser, go to `http://localhost:7001` to view diary note application.

## Demo
https://app.diarynote.net
