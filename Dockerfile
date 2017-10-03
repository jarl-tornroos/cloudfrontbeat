# Use an official Alpine 3.6 as a parent image
FROM alpine:3.6

# Set the working directory to /cloudfrontbeat
WORKDIR /cloudfrontbeat

# Install the CA certificates
RUN apk add --update ca-certificates

# Copy cloudfrontbeat into the container at /usr/local/bin
ADD cloudfrontbeat /usr/local/bin

# Copy the configuration files into the container at /cloudfrontbeat
ADD *.json /cloudfrontbeat/
ADD *.yml /cloudfrontbeat/

# Run cloudfrontbeat when the container launches
CMD ["cloudfrontbeat", "-e"]
