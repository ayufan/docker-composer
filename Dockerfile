FROM docker/compose:1.7.0
RUN ["apk", "add", "-U", "git"]
RUN ["git", "config", "--global", "init.templatedir", "/git-templates"]
RUN ["git", "config", "--global", "receive.denyCurrentBranch", "updateInstead"]

ENV APPS_DIR=/srv/apps
VOLUME ["/srv/apps"]
ADD / /

ENTRYPOINT ["/entrypoint.sh"]
