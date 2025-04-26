# Generate a private key for your CA
openssl genrsa -out ca.key 2048

# Generate a CA certificate
openssl req -x509 -new -nodes -key ca.key -subj "/CN=webhook-ca" -days 3650 -out ca.crt

# Create a private key for the server
openssl genrsa -out server.key 2048

# Create a certificate signing request (CSR)
openssl req -new -key server.key -subj "/CN=leader-election-webhook.default.svc" -out server.csr

# Create a config file for the certificate
cat > server.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = leader-election-webhook
DNS.2 = leader-election-webhook.default
DNS.3 = leader-election-webhook.default.svc
DNS.4 = leader-election-webhook.default.svc.cluster.local
EOF

# Sign the CSR with your CA
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -extfile server.ext

kubectl create secret tls webhook-certs \
  --cert=server.crt \
  --key=server.key \
  -n default
