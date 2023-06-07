# Pulling Alpine image
FROM alpine:latest

# Install required packages
RUN apk add --update --no-cache bash dos2unix

WORKDIR /usr/scheduler

# Copy files
COPY jobs/*.* ./jobs/
COPY start.sh .

# Fix line endings && execute permissions
RUN dos2unix *.sh jobs/*.* \
    && \
    find . -type f -iname "*.sh" -exec chmod +x {} \;

# Run cron on container startup
CMD ["./start.sh"]