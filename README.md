# Celcat to ICS Converter


NU London students use Celcat as the primary timetabling solution, but its built-in calendar has terrible support for essential functions like event notifications and sorting. Many use Google Calendar, iCalendar or similar calendar solutions for managing their calendars. This tool takes the Celcat event response and converts it to a subscribable ICS weblink. 

This tool is only intended for Northeastern University students, and *should* only work with their northeastern.edu email addresses. 


# The Basics

Currently I'm hosting this version on [aplb.xyz/celcat/*](https://aplb.xyz/celcat/*) in an ECS free-tier container. No data is collected from responses besides rudimentary error logging in cloudwatch. 

There is also a basic UI for intereacting with the service at [celcat.aplb.xyz](https://celcat.aplb.xyz/)

## Stack
This is a Golang api using Docker for containerized hosting. 

## Locally Testing
You can fork or download this repository and host the API. It runs on port 8000 and the primary endpoint is */getCalendar. There are no other dependencies. 
