# How to Use
1. Clone the repo.
```
git clone https://github.com/mdanialr/webhook.git
```
2. Get dependencies.
```
go mod tidy
```
3. Create new config file.
```
cp app-config.yaml.example app-config.yaml
```
4. Fill in the app-config.yaml file as needed.
5. Build the project.
```
go build -o bin/webhook main.go
```
or use makefile instead
```
make
```
6. Run the binary file.
```
./bin/webhook
```
## Optional
7. Integrate with systemd.

_systemd script_
```
[Unit]
Description=instance to serve webhook service
After=network.target

[Service]
User=root
Group=your-username
ExecStart=/bin/sh -c "cd /path/to/cloned/repo && ./bin/webhook"

[Install]
WantedBy=multi-user.target
```
8. Save _systemd script_ in `/etc/systemd/system/` something like `webhook.service`.
9. Run and enable systemd, so it will run even after reboot.
```
sudo systemctl enable webhook.service --now
```