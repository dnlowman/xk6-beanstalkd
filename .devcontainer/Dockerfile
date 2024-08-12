FROM golang:1.22.6

# Install xk6
RUN go install go.k6.io/xk6/cmd/xk6@latest

# Install Beanstalkd
RUN apt-get update && apt-get install -y beanstalkd

# Set the working directory
WORKDIR /workspace

# Start Beanstalkd
CMD ["beanstalkd", "-l", "0.0.0.0", "-p", "11300"]