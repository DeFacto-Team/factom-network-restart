# Factom Network Restart System
The app connects to remote Portainer server via Portainer API and asynchronously restarts all containers with label `name=factomd`.

## Configuration
Config file with Portainer credentials should be placed in `~/.factom-restart/config.yaml` or provided as a flag `-c /path/to/config.yaml` while running the application.

## Usage
### Dry-run
Run in dry-run mode (scans Portainer endpoints and containers, no restart happens):
```bash
go run main.go
```
Output:
```bash
INFO[0000] Using config: /Users/anton/.factom-restart/config.yaml 
INFO[0000] Starting restart system                      
INFO[0000] Portainer endpoint: https://test.com 
INFO[0000] DRY RUN mode: simulating restart             
INFO[0001] Successfully logged in as anton
------------------------------
DEBU[0001] Trying connecting to X (xxx.xxx.xxx.xxx), ID=2 
ERRO[0001] Can not connect to remote host
------------------------------
DEBU[0003] Trying connecting to Y (yyy.yyy.yyy.yyy), ID=29 
DEBU[0003] factomd container is running
307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126
v6.6.0-alpine
…
------------------------------
INFO[0033] SIMULATING RESTART (DRY-RUN)                 
WARN[0033] SKIP X                                       
INFO[0033] OK Y (307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126)
…
```

### Live mode
Run in live mode (restarts the network):
```bash
go run main.go --live
```
Output:
```bash
INFO[0000] Using config: /Users/anton/.factom-restart/config.yaml 
INFO[0000] Starting restart system                      
INFO[0000] Portainer endpoint: https://test.com 
WARN[0000] LIVE mode: network will be restarted!             
INFO[0001] Successfully logged in as anton
------------------------------
DEBU[0001] Trying connecting to X (xxx.xxx.xxx.xxx), ID=2 
ERRO[0001] Can not connect to remote host
------------------------------
DEBU[0003] Trying connecting to Y (yyy.yyy.yyy.yyy), ID=29 
DEBU[0003] factomd container is running
307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126
v6.6.0-alpine
…
------------------------------
INFO[0033] RESTARTING NOW                 
WARN[0033] SKIP X                                       
INFO[0033] OK Y (307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126)
…
```
