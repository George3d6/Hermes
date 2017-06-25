#!/bin/bash

printf "\e[1;33m"

printf "Setting up environment !\n"
echo 'This is a test upload file' > test_upload.txt;

cd ..;
rm -rf ser;
rm -rf storage;
mkdir -p ser;
mkdir -p storage;

printf "Environment set up successfully !\n"

printf "\e[0m"
printf "\e[1;32m"

printf "Building and launching the hermes file server !\n\n"

printf "\e[0m"

cd server && go build && mv server ../hermes && cd ..;
./hermes 'config.json' 'admin' 'admin' &
PROC_ID=$!;
sleep 2s;

printf "\e[1;32m"

printf "\nThe server \"should\" have been launched \"successfully\" !\n (if you encounter problems modify the sleep or re write the script to wait until stdout gets the server's launch message)\n"

printf "Running regression testing using phantomjs !\n\n"

printf "\e[0m"

cd tests;
node regression.js;

printf "\e[1;35m"

printf "\nRegression ran !\n";

printf "Cleaning everything up !\n";

kill -9 $PROC_ID;
cd ..;
rm -rf ser;
rm -rf storage;
cd tests;
rm -rf test_upload.txt;

printf "Cleanup completed, congratulations, you have run the hermes test suite !\n";

printf "\e[0m"
