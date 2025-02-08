# IP CIDR Checker Service

A RESTful service that checks if an IP address belongs to a configured list of CIDR ranges.

## Configuration

The service is configured via `config.yaml`:

```yaml
server:
  port: 8080    # Port the service will listen on
cidr:
  file: "cidrs.txt"    # File containing CIDR ranges, one per line
```

## CIDR List Format

The CIDR list file should contain one CIDR range per line, for example:

```
192.168.0.0/24
10.0.0.0/8
172.16.0.0/12
```

## API Endpoint

### Check IP

```
GET /check/:ip
```

Parameters:
- `:ip` - The IP address to check

Response:
```json
{
    "ip": "192.168.1.1",
    "matches": true
}
```

## Running the Service

```bash
go run main.go
```

```shell
curl "http://localhost:8080/check/128.106.183.190"
curl "http://localhost:8080/check/120.244.38.93"
```
