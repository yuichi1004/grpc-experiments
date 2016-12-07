# Generate Certs

## Generate Root CA

```bash
$ openssl req -new -x509 -newkey rsa:2048 -config ./ssl.conf -out ./ca/cacert.pem -keyout ./ca/private/cakey.pem -days 3650
Country Name (2 letter code) [JP]:
State or Province Name (full name) [Tokyo]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Example Company]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:example.com
Email Address []:
```

## Generate Server Cert Requests

```bash
$ openssl req -config ./ssl.conf -new -keyout tokenkey.pem -out tokencsr.pem
Country Name (2 letter code) [JP]:
State or Province Name (full name) [Tokyo]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Example Company]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:token.example.com
Email Address []:

$ openssl rsa -in ./tokenkey.pem -out ./tokenkey.pem

$ openssl req -config ./ssl.conf -new -keyout fibokey.pem -out fibocsr.pem
Country Name (2 letter code) [JP]:
State or Province Name (full name) [Tokyo]:
Locality Name (eg, city) []:
Organization Name (eg, company) [Example Company]:
Organizational Unit Name (eg, section) []:
Common Name (e.g. server FQDN or YOUR name) []:fibo.example.com
Email Address []:

$ openssl rsa -in ./fibokey.pem -out ./fibokey.pem
```

## Sign Cert

```bash
$ openssl ca -in tokencsr.pem -out tokencert.pem -config ./ssl.conf
$ openssl ca -in fibocsr.pem -out tokencert.pem -config ./ssl.conf
```

