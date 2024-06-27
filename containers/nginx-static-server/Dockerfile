FROM nginx:latest

# Set working directory to the Nginx document root
WORKDIR /usr/share/nginx/html

# Download the Terraform binary zip file, clean up and reduce the size of the Docker image by removing unnecessary files after installing packages.
RUN apt-get update && apt-get install -y curl && \
    mkdir -p terraform/1.4.7 && \
    curl -o terraform/1.4.7/terraform_1.4.7_linux_amd64.zip https://releases.hashicorp.com/terraform/1.4.7/terraform_1.4.7_linux_amd64.zip && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Start the NGINX server
CMD ["nginx", "-g", "daemon off;"]
