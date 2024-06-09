#!/bin/bash
sleep 5
bee migrate -conn="webapp:qwe123@tcp(db:3306)/koocbook"
go run main.go | tee logs/app.log
