# Use an official Node.js image as the base image
FROM node:18-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the package.json and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the application code
COPY . .

# Set environment variable for VITE_BACKEND_URL
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=${VITE_BACKEND_URL}

# Build the frontend application
RUN npm run build

# Use a minimal base image for the final stage
# FROM nginx:alpine
FROM bitnami/nginx:latest

# Copy the built frontend files to the Nginx HTML directory
# COPY --from=builder /app/dist /usr/share/nginx/html
COPY --from=builder /app/dist /opt/bitnami/nginx/html

# Expose the port Nginx will serve on
EXPOSE 8080

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]