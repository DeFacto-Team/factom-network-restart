# Factom Network Restart System
The app connects to remote Portainer server via Portainer API and asynchronously restarts all containers with label `name=factomd`.

## Configuration
Config file with Portainer credentials should be placed in `~/.factom-restart/config.yaml` or provided as a flag `-c /path/to/config.yaml` while running the application.

## Usage
### Dry-run
Run in dry-run mode (scans Portainer endpoints and containers, no restart happens):
```bash
go run .
```
Output:
```bash
2020/09/15 13:49:01 Using config: /Users/anton/.factom-restart/config.yaml 
2020/09/15 13:49:01 Starting restart system                      
2020/09/15 13:49:01 Portainer endpoint: https://test.com 
2020/09/15 13:49:01 DRY RUN mode: simulating restart             
2020/09/15 13:49:01 Successfully logged in as anton
------------------------------
2020/09/15 13:49:01 Trying connecting to X (xxx.xxx.xxx.xxx), ID=2 
2020/09/15 13:49:01 Empty response received from Portainer API
------------------------------
2020/09/15 13:49:02 Trying connecting to Y (yyy.yyy.yyy.yyy), ID=29 
2020/09/15 13:49:02 factomd container is running
307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126
v6.6.0-alpine
…
------------------------------
2020/09/15 13:49:05 SIMULATING RESTART (DRY-RUN)                 
2020/09/15 13:49:05 SKIP X                                       
2020/09/15 13:49:05 OK Y (307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126)
…
```

### Live mode
Run in live mode (restarts the network):
```bash
go run . --live
```
Output:
```bash
2020/09/15 13:49:01 Using config: /Users/anton/.factom-restart/config.yaml 
2020/09/15 13:49:01 Starting restart system                      
2020/09/15 13:49:01 Portainer endpoint: https://test.com 
2020/09/15 13:49:01 LIVE mode: network will be restarted             
2020/09/15 13:49:01 Successfully logged in as anton
------------------------------
2020/09/15 13:49:01 Trying connecting to X (xxx.xxx.xxx.xxx), ID=2 
2020/09/15 13:49:01 Empty response received from Portainer API
------------------------------
2020/09/15 13:49:02 Trying connecting to Y (yyy.yyy.yyy.yyy), ID=29 
2020/09/15 13:49:02 factomd container is running
307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126
v6.6.0-alpine
…
------------------------------
2020/09/15 13:49:05 RESTARTING NOW               
2020/09/15 13:49:05 SKIP X                                       
2020/09/15 13:49:05 OK Y (307f871a97e1dd91851c6298cf3a183f1709d908810b0059ce9094622bb78126)
…
```
