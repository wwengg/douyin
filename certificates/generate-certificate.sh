#!/bin/bash
basepath=$(dirname "$0")

openssl genrsa -out "$basepath/proxy-ca.key" 2048
openssl req -x509 -new -nodes -key "$basepath/proxy-ca.key" -sha256 -days 1825 -out "$basepath/proxy-ca.crt"

