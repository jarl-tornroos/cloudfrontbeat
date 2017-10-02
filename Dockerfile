# Use an official Ubuntu 16.04 as a parent image
FROM ubuntu:16.04

# Set the working directory to /cloudfrontbeat
WORKDIR /cloudfrontbeat

# Copy cloudfrontbeat into the container at /usr/local/bin
ADD cloudfrontbeat /usr/local/bin

# Copy the configuration files into the container at /cloudfrontbeat
ADD *.json /cloudfrontbeat/
ADD *.yml /cloudfrontbeat/

# Run cloudfrontbeat when the container launches
CMD ["cloudfrontbeat", "-e"]
