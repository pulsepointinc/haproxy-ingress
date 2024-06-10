# Copyright 2021 The HAProxy Ingress Controller Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.19.13-bullseye AS builder

ARG GIT_TAG
WORKDIR /src
RUN apt install git

COPY . .

RUN GIT_TAG=${GIT_TAG:-latest}\
 && GIT_COMMIT=git-$(git rev-parse --short HEAD)\
 && GIT_REPO=$(git config --get remote.origin.url)\
 && VERSION_PKG=github.com/jcmoraisjr/haproxy-ingress/pkg/version\
 && CGO_ENABLED=0 go build\
    -ldflags "-s -w -X ${VERSION_PKG}.RELEASE=${GIT_TAG} -X ${VERSION_PKG}.COMMIT=${GIT_COMMIT} -X ${VERSION_PKG}.REPO=${GIT_REPO}"\
    -o haproxy-ingress pkg/main.go

FROM registry.pulsepoint.com/haproxy:v2.4_ppfp_0.6.0

USER root

RUN apt-get update \
     && apt-get install -y --no-install-recommends \
     patch \
     socat lua5.3 lua-json lua-socket dumb-init \
     && rm -rf /var/lib/apt/lists/*

COPY rootfs/ /
COPY --from=builder /src/haproxy-ingress /haproxy-ingress-controller

RUN addgroup -gid 1001 haproxy\
 && adduser --uid 1001 --gid 1001 --disabled-password --no-create-home --shell /bin/false haproxy \
 && mkdir -p /var/empty /etc/haproxy /var/lib/haproxy /var/run/haproxy\
 && chown -R haproxy:haproxy /etc/haproxy /var/lib/haproxy /var/run/haproxy\
 && chmod 0 /var/empty

# TODO it's not a great way to do it as we need to keep haproxy-pp in sync, rework this into using configuration snippets
# See warning here: https://haproxy-ingress.github.io/docs/configuration/template/
RUN patch /etc/templates/haproxy/haproxy.tmpl /etc/templates/haproxy/haproxy-pp.patch

STOPSIGNAL SIGTERM
ENTRYPOINT ["/usr/bin/dumb-init", "--", "/start.sh"]
