FROM alpine:3.13

ARG BUILD_DATE
ARG CI_PROJECT_TITLE
ARG CI_PROJECT_URL
ARG CI_COMMIT_SHORT_SHA
ARG CI_COMMIT_REF_NAME
ARG GITLAB_USER_LOGIN

LABEL edu.psu.base-image.created=$BUILD_DATE \
      edu.psu.base-image.title=$CI_PROJECT_TITLE \
      edu.psu.base-image.source=$CI_PROJECT_URL \
      edu.psu.base-image.revision=$CI_COMMIT_SHORT_SHA \
      edu.psu.base-image.vendor="PSU-SWE" \
      edu.psu.base-image.version=$CI_COMMIT_REF_NAME \
      edu.psu.base-image.authors=$GITLAB_USER_LOGIN

RUN apk update \
 && apk upgrade \
 && apk add --no-cache ca-certificates tzdata \
 && rm -rf /var/cache/apk/*

# work-around for ipv6 localhost
COPY ./nsswitch.conf /etc/nsswitch.conf

COPY ./demoservice /usr/local/bin/demoservice

ENTRYPOINT ["demoservice"]