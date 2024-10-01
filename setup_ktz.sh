#!/bin/sh

# Set the project name
PROJECT_NAME="ktz"

# Build the Go application
echo "Building the Go application..."
go build -o $PROJECT_NAME

# Check if the build was successful
if [ $? -ne 0 ]; then
    echo "Build failed. Please check for errors."
    exit 1
fi

# Move the project to /usr/local/bin
echo "Moving $PROJECT_NAME to /usr/local/bin..."
sudo mv $PROJECT_NAME /usr/local/bin/

# Ensure the project is executable
echo "Setting execute permissions on /usr/local/bin/$PROJECT_NAME..."
sudo chmod +x /usr/local/bin/$PROJECT_NAME

# Verify if ktz is available in the PATH
if command -v $PROJECT_NAME >/dev/null 2>&1; then
    echo "$PROJECT_NAME is now available globally. You can use it as: ktz"
else
    echo "Something went wrong. Please check your PATH or permissions."
fi

