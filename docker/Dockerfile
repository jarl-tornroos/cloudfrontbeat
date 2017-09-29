# Use an official ubuntu 16.04 as a parent image
FROM ubuntu:16.04

# Set the working directory to /app
WORKDIR /app

# Copy compiled version and configuration file into the container at /app
ADD ../cloudfrontbeat /app/
ADD ../cloudfrontbeat.yml /app/

# Run cloudfrontbeat when the container launches
CMD ["cloudfrontbeat", "-c", "cloudfrontbeat.yml", "-e"]
