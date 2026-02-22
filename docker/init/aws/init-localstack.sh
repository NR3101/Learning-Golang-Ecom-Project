#!/bin/bash

# Create bucket
awslocal s3 mb s3://ecom-uploads


# Create a sqs queue
awslocal sqs create-queue --queue-name ecom-events


echo "LocalStack initialization complete"
