# Copyright 2017 The Kubernetes Authors. All rights reserved.
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

FROM registry.pulsepoint.com/haproxy:v2.4_ppfp_0.5.0

USER root

RUN apt-get update \
     && apt-get install -y --no-install-recommends \
     patch \
     socat openssl lua5.3 lua-json lua-socket dumb-init

COPY . /

RUN addgroup --gid 1001 haproxy \
 && adduser --uid 1001 --gid 1001 --disabled-password --no-create-home --shell /bin/false haproxy \
 && mkdir -p /var/empty /etc/haproxy/maps /var/lib/haproxy /var/run/haproxy \
 && chown -R haproxy:haproxy /run /etc/haproxy /var/lib/haproxy /var/run/haproxy \
 && chmod 0 /var/empty

# TODO it's not a great way to do it as we need to keep haproxy-pp in sync, rework this into using configuration snippets
# See warning here: https://haproxy-ingress.github.io/docs/configuration/template/
RUN patch /etc/templates/haproxy/haproxy.tmpl /etc/templates/haproxy/haproxy-pp.patch

STOPSIGNAL SIGTERM

# dumb-init reaps the old haproxy process in the embedded, non master-worker mode,
# after receiving SIGUSR1, avoiding it to become a zombie.
ENTRYPOINT ["/usr/bin/dumb-init", "--", "/start.sh"]
