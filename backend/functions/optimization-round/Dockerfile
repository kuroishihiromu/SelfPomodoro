# AWS Lambda Python Runtime Interface Client for Container Images (x86_64)
FROM public.ecr.aws/lambda/python:3.9-x86_64

# Install system dependencies
RUN yum update -y && \
    yum install -y gcc g++ && \
    yum clean all

# Copy requirements first for better caching
COPY requirements.txt ${LAMBDA_TASK_ROOT}/

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy function code
COPY handler.py ${LAMBDA_TASK_ROOT}/

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "handler.lambda_handler" ]
