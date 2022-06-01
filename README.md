# Webhook
A service that will listen to incoming request and do a CD job based on the incoming request and which routes were used

## Features
* `GitHub Webhook` `/hook/{repo-name}`: listening to GitHub webhook.
* `GitHub Actions` `/github/webhook`: listening to GitHub actions specifically [this](https://github.com/distributhor/workflow-webhook) action.
* `Docker Webhook` `/docker/webhook`: listening to Docker Hub webhook.

## Configuration Table in `app-config.yaml`
| Variable     | Required | Default   | Choices              | Description                                                                                                                                                |
|--------------|----------|-----------|----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| env          | &#9744;  | dev       | `dev` &#124; `prod`  | Which stage this service is running. `dev` would print log in terminal & `prod` in log file                                                                |
| host         | &#9744;  | localhost | -                    | Hostname where this service is running                                                                                                                     |
| port         | &#9744;  | 5050      | -                    | HTTP port where this service is attached to                                                                                                                |
| secret       | &#9745;  | -         | -                    | Shared secret string from GitHub webhook                                                                                                                   |
| log          | &#9744;  | ./log/    | -                    | Where log files would be written if `env`: `prod`                                                                                                          |
| max_worker   | &#9744;  | 1         | any unsigned integer | Number of workers for each type of service to do the CD job. Service type: `GitHub Webhook`, `GitHub Actions` & `Docker Webhook`                           |
| service      | &#9744;  | -         | -                    | Contain list of `repo` that would be used by `Github Webhook` & `Github Actions`                                                                           |
| repo         | &#9744;  | -         | -                    | Contain each `repo` that would be used by `Github Webhook` & `Github Actions`                                                                              |
| repo.user    | &#9744;  | user      | -                    | Username that own this repository                                                                                                                          |
| repo.name    | &#9745;  | -         | -                    | The name of the repo. For `GitHub Webhook` this must match with `{repo-name}` in the target url of GitHub webhook                                          |
| repo.branch  | &#9744;  | master    | `tags` &#124; `*`    | Branch name of the repo. Could be any valid branch name but if __listening__ to tags (__ref/tags/*__) this should be `tags`. Only used by `GitHub Actions` |
| repo.path    | &#9745;  | -         | -                    | Where this `repo` is located in local and the CD job is executed                                                                                           |
| repo.opt_cmd | &#9744;  | -         | -                    | Additional bash command that would be executed right after the CD job                                                                                      |
| dockers      | &#9744;  | -         | -                    | Contain list of `dcoker` that would be used by `Docker Webhook`                                                                                            |
| docker       | &#9744;  | -         | -                    | Contain each `docker` that would be used by `Docker Webhook`                                                                                               |
| docker.user  | &#9745;  | -         | -                    | Username to authenticate to Docker Hub                                                                                                                     |
| docker.pass  | &#9745;  | -         | -                    | Password to authenticating this docker.`user` to Docker Hub                                                                                                |
| docker.repo  | &#9745;  | -         | -                    | Name of the repository that would be monitored                                                                                                             |
| docker.tag   | &#9744;  | latest    | -                    | Tag name of this docker.`repo` that would be monitored                                                                                                     |
| docker.args  | &#9744;  | -         | -                    | Additional arguments that would be pass in along when exectuing __docker run *args*__. Ideally is used to pass in port mapping                             |

## How to Use
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
### Optional
7. Integrate with systemd.

    _systemd script_
    ```
    [Unit]
    Description=instance to serve webhook service After=network.target

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

# License
This project is licensed under the **MIT License** - see the [LICENSE](LICENSE "LICENSE") file for details.