# Client Auth

## CA

1. **CA** Generate CA
	- `ca.key`
	- `ca.crt`
	
## New Server

1. **Server** Generate Private key and CSR
	- `server.key`
	- `server.csr`
1. **Server** Transfer `server.csr` to **CA**
1. **CA** Signs Server's CSR using `x509.ExtKeyUsageServerAuth`
	- input `ca.key`, `ca.crt`, `server.csr`
	- output `server.crt`
1. **CA** Transfer `server.crt`, `ca.crt` to **Server**

## New Client

1. **Client** Generate Private key and CSR
	- `client.key`
	- `client.csr`
1. **Client** Transfer `client.csr` to **CA**
1. **CA** Signs Client's CSR using `x509.ExtKeyUsageClientAuth`
	- input: `ca.key`, `ca.crt`, `client.csr`
	- output: `client.crt`
1. **CA** Transfer `client.crt`, `ca.crt` to **Client**

## Handshake

1. **Client** verify **Server**
1. **Server** verify **Client**

### Server verify client concept

1. **Client** Connect to **Server**
1. **Client** Send `client.crt` to **Server**
1. **Server** Verify `client.crt` with `ca.crt`
1. **Server** Send `xxx` to **Client**
1. **Client** use `client.key` to signs `xxx`
1. **Client** Send `signed xxx` to **Server**
1. **Server** Verify `signed xxx` with `client.crt`

### Example

- Postgres: `postgres://client@server-ip/db?sslmode=verify-full&sslkey=client.key&sslcert=client.crt&sslrootcert=ca.crt`
