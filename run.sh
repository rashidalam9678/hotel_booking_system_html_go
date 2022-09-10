#!/bin/bash

go build -o bookings cmd/web/*.go && ./bookings -dbname=hotel_bookings -dbuser=mr.mra -cache=true -production=true
