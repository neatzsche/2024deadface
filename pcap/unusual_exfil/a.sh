#!/bin/zsh

# Define the host and port
HOST="192.168.50.188"
PORT="6855"

# Path to the text file
FILE="QR.txt"

# Check if the file exists
if [[ ! -f "$FILE" ]]; then
  echo "File not found!"
  exit 1
fi

# Iterate over each line in the file
while IFS= read -r line; do
  echo "$line" | nc "$HOST" "$PORT"
  echo "Sent: $line"
  
done < "$FILE"