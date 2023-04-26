# Registering and installing certificate

```sh
./generate-certificate.sh

# install the cerificate (on macOS)
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain proxy-ca.crt
```
