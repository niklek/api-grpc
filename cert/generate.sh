DURATION_DAYS=365

rm *.pem

# For server side TLS
# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days $DURATION_DAYS -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=AT/ST=Vienna/L=Vienna/O=Test CA/OU=Education/CN=myapi.local/emailAddress=testca@example.com"

# CA's self-signed certificate
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=AT/ST=Vienna/L=Vienna/O=Test Server/OU=Computer/CN=myapi.local/emailAddress=test@example.com"

# 3. Use CA's private key to sign server's CSR and get back the signed certificate
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
# Server's signed certificate
openssl x509 -in server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=AT/ST=Vienna/L=Vienna/O=Test Client/OU=Computer/CN=myapi.local/emailAddress=test@example.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf
# Client's signed certificate
openssl x509 -in client-cert.pem -noout -text
