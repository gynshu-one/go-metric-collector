#!/bin/bash

# Generate a private key
openssl genpkey -algorithm RSA -out private.pem

# Extract the public key
openssl rsa -pubout -in private.pem -out public.pem
