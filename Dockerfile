# Use an official Python runtime as a parent image
FROM python:3.12-slim

# Set the working directory in the container
WORKDIR /app

# Copy the dependency definition file
COPY pyproject.toml .

# Install the dependencies
# We use --no-cache-dir to keep the image size down
RUN pip install --no-cache-dir .

# Copy the rest of the application source code
COPY src/ ./src
COPY configs/ ./configs

# Expose the port the app runs on
EXPOSE 8080

# Define the command to run the application
CMD ["python3", "-m", "src"]