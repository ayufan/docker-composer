# docker-composer

Git push-to-deploy for Docker Compose applications.

Every application is a git repository on the server holding a
`docker-compose.yml`. You push to it, and the push itself builds the images
and brings the stack up. There is no agent to install on the server and no
control plane: the whole thing is a single container, and the `compose` user's
login shell *is* that container, so every operation arrives over SSH.

```
you@laptop $ git push compose@server:blog master
...
remote: Applying new changes...
remote: Deploying application...
remote: Tagging release...
```

## How it works

Applications live in `APPS_DIR` (`/srv/apps` by default), one directory per
application, each a git repository with a checked-out working tree:

```
/srv/apps
├── blog/
│   ├── .git/
│   ├── docker-compose.yml
│   └── .env            # optional, untracked
└── demo/
    ├── .git/
    └── docker-compose.yml
```

Three pieces connect a `git push` to a running stack:

1. **The login shell.** `compose`'s shell is `run.sh`, which does nothing but
   `docker run` the image with the arguments it was given. SSH invokes a login
   shell as `shell -c "<command>"`, and `docker-composer` accepts a global `-c`
   flag that re-parses its value as its own command line. So `ssh
   compose@server ls` becomes the `ls` command, and git's
   `git-receive-pack 'blog'` becomes the `git-receive-pack` command scoped to
   the `blog` application. The container is the shell, so nothing but Docker
   and SSH is needed on the host.

2. **`receive.denyCurrentBranch=updateInstead`.** Pushing to the branch a
   repository currently has checked out is normally refused. Each application
   is configured to accept it instead, which makes git run the
   `push-to-checkout` hook to update the working tree.

3. **The `push-to-checkout` hook.** `init` installs it as a symlink back to
   the `docker-composer` binary. When git runs it, the binary notices it was
   invoked from `.git/hooks`, derives the application name from
   `GIT_WORK_TREE`, and dispatches to its own `git-push-to-checkout` command:
   it reads the pushed revision into the working tree, runs
   `docker-compose build` and `up -d`, and moves the `latest` tag to the
   revision just deployed. If any step fails the working tree is reset, so a
   broken push leaves the previous deployment running.

That `latest` tag is the rollback point. `switch blog latest` returns to the
last deployment that came up cleanly.

## Install

On the Docker host, as root:

```sh
curl -L https://raw.githubusercontent.com/ayufan/docker-composer/main/run.sh > /usr/local/bin/docker-composer
chmod +x /usr/local/bin/docker-composer
docker pull ayufan/docker-composer
docker-composer ls
```

The first run creates the `compose` user if it is missing: shell set to
`/usr/local/bin/docker-composer`, a member of the `docker` group, and
`~/.ssh/authorized_keys` copied from the invoking user so the keys that reach
the server already work. Add more keys to `~compose/.ssh/authorized_keys`;
anyone holding one can deploy.

`docker-composer --update` replaces the script with the current version from
this repository, and `docker-composer self-upgrade` pulls a newer image.

## Example

Create an application and give it a stack:

```sh
$ docker-composer init demo
Initialized empty Git repository in /srv/apps/demo/.git/

$ docker-composer edit demo        # $EDITOR on docker-compose.yml
```

```yaml
services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
```

`edit` validates the file with `docker-compose config -q` before it does
anything — on a syntax error it hands the editor back rather than deploying —
then deploys, commits and tags. The same validation guards `env edit`.

```sh
$ docker-composer ls --detail
APP   STATUS       REF
demo  running      master

$ docker-composer logs demo -f      # any docker-compose command, scoped to the app
```

Now deploy to it from a working copy:

```sh
$ git clone compose@server:demo && cd demo
$ vi docker-compose.yml
$ git commit -am "Expose web on 8080"
$ git push origin master
remote: Applying new changes...
remote: Deploying application...
remote: Tagging release...
```

> `init` creates the repository with `master` checked out, and a push only
> triggers a deployment when it targets the checked-out branch. If your client
> defaults to `main`, push with `git push origin main:master`.

Secrets that should not live in the repository go in `.env`, which
`docker-compose` reads automatically:

```sh
$ docker-composer env edit demo    # edit, validate, deploy, tag
$ docker-composer env view demo
```

`env` deliberately refuses to work once `.env` is tracked by git — an app that
version-controls its own `.env` owns it, and a server-side edit would be
overwritten by the next push.

If a deployment misbehaves, go back to the last good revision, or stop the
application from being deployed at all:

```sh
$ docker-composer switch demo latest
$ docker-composer disable demo     # `deploy --all` skips it; `enable` undoes this
```

Everything above works the same over SSH, without a shell on the server:

```sh
$ ssh compose@server ls --detail
$ ssh compose@server deploy --all
```

## Commands

Run `docker-composer help` for the authoritative list.

| | |
|---|---|
| `ls [APP...]` | list applications, `--detail` adds status and branch |
| `stats [APP...]` | resource usage per application |
| `deploy APP...` | build, pull and bring up; `--all` for every enabled application |
| `init APP` | create an application and install its hook |
| `release APP` | tear down and delete an application, with `--volumes` and `--rmi` |
| `edit APP` | edit `docker-compose.yml`, then validate, deploy, commit, tag |
| `env view\|edit APP` | manage the untracked `.env` |
| `switch APP REF` | deploy another revision, reverting if it fails to come up |
| `enable\|disable APP` | control whether `deploy --all` includes it |
| `services APP` | list the services of an application |
| `cleanup` | prune containers, volumes and images, each behind its own flag |
| `update-hooks [APP...]` | reinstall hooks, needed after upgrading the image |
| `self-upgrade` | pull a newer `docker-composer` image |

The `docker-compose` subcommands — `build`, `config`, `create`, `down`,
`events`, `exec`, `kill`, `logs`, `pause`, `port`, `ps`, `pull`, `restart`,
`rm`, `run`, `scale`, `start`, `stop`, `unpause`, `up` — are passed straight
through, taking the application name first: `docker-composer exec blog web sh`.

## Configuration

| | |
|---|---|
| `APPS_DIR` | where applications are stored, default `/srv/apps` |
| `DOCKER_COMPOSER_TAG` | image tag `run.sh` runs, default `latest` |
| `DOCKER_HOST` | bind-mounted when it is a socket, forwarded with `DOCKER_TLS_VERIFY` and `DOCKER_CERT_PATH` otherwise |
| `EDITOR`, `VISUAL`, `GIT_EDITOR` | editor for `edit` and `env edit` |
| `LOG_LEVEL`, `DEBUG` | `debug`, `info`, `warn`, `error`, `fatal`, `panic` |

## Building

```sh
make                 # build the images for every architecture
make version         # print the docker, docker-compose and git versions in the image
make dockerhub       # build and push, TAG defaults to beta
make tag-latest      # publish as latest
./test.sh ls         # build and run a command against the local image
```

The images are multi-architecture (`amd64`, `arm32v7`, `arm64v8`) and
published to [`ayufan/docker-composer`](https://hub.docker.com/r/ayufan/docker-composer).
Pushing to `main` publishes `beta`; the *Build* and *Make latest* workflows
publish versions on demand.

## License

[MIT](LICENSE), Kamil Trzciński.
