# Pulling Alpine image
FROM alpine:latest

# Install required packages
RUN apk add --update --no-cache curl bash dos2unix

WORKDIR /usr/scheduler

ENV SCHEDULER_ENVIRONMENT=production

# Copy files
COPY . .

# Fix line endings && execute permissions
RUN dos2unix *.sh jobs/*.* \
    && \
    find . -type f -iname "*.sh" -exec chmod +x {} \;

# Run cron on container startup
CMD ["./start.sh"]