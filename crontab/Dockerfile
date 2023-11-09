# Pulling Alpine image
FROM alpine:latest

# Install required packages
RUN apk add --update --no-cache curl bash dos2unix tzdata

WORKDIR /usr/scheduler

ENV SCHEDULER_ENVIRONMENT=production
ENV TZ=Asia/Jakarta

COPY . .

# Fix line endings && execute permissions
RUN dos2unix *.sh jobs/*.* \
    && \
    find . -type f -iname "*.sh" -exec chmod +x {} \;

# execute permissions fro start.sh
RUN chmod +x start.sh

# Run cron on container startup
CMD ["start.sh"]