# Webhook
A service that will listen to incoming request and do a CD job based on the incoming request payload that sent by
[this](https://github.com/distributhor/workflow-webhook) GitHub Actions in the url `/github/webhook`.

## How to Use
1. Download the binary from [GitHub Releases](https://github.com/mdanialr/webhook/releases)
2. Create new config file with a filename `app.yml`. __The file name `app.yml`__ is mandatory otherwise
[Viper](https://github.com/spf13/viper) will not find it
3. Extract then run to check if there is any error in config file
 ```bash
tar -xzf webhook....tar.gz
./webhook
```

### Example
Create config file with the filename `app.yml`
```yaml
env: prod # default is dev when not set. prod will write all logs to the files in the designated folder in log keyword below.
host: 127.0.0.1 # default is 127.0.0.1
port: 7575 # default is 7575
secret: secret # this secret must be the same as in the GitHub Actions workflow
log: /path/to/where/log/file/is/stored # default is /tmp
max_worker: 1 # default is 1. this is worker count that will run the 'commands'. 1 is already good enough.
github:
  - name: # required
    user: # required
    branch: # required
    tags: # optional, default is false. this will listen only to tags ref heads.
    event: # optional, default is push. this will listen only to push event.
    path: # required. this the path where the 'commands' will be executed.
    commands: | # optional. this is the commands that will be executed whenever the event is triggered.
      echo hello
      echo world
```

### Optional (_Integrate with systemd_)
```bash
[Unit]
Description=instance to serve webhook service
After=network.target

[Service]
User=root
Group=your-username
ExecStart=/bin/sh -c "cd /path/to/binary/file && ./webhook"

[Install]
WantedBy=multi-user.target
```
1. Save above _systemd script_ in `/etc/systemd/system/` with a filename maybe something like `webhook.service`.
2. Run and enable systemd, so it will run even after reboot.
 ```bash
sudo systemctl enable webhook.service --now
```

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.
