FROM docker/compose:1.7.0

RUN ["apk", "add", "-U", "git", "bash", "docker", "nano", "perl", "sed"]
RUN ["git", "config", "--global", "init.templatedir", "/git-templates"]
RUN ["git", "config", "--global", "receive.denyCurrentBranch", "updateInstead"]
RUN ["git", "config", "--global", "user.name", "Composer"]
RUN ["git", "config", "--global", "user.email", "you@example.com"]

ENV APPS_DIR=/srv/apps
ADD /demo /srv/apps/demo/
VOLUME ["/srv/apps"]
ADD / /

ENTRYPOINT ["/composer"]
