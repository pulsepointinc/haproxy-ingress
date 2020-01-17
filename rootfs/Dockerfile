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

FROM registry.pulsepoint.com/haproxy:v1.8_pp

# dumb-init kindly manages SIGCHLD from forked HAProxy processes
ARG DUMB_INIT_SHA256=81231da1cd074fdc81af62789fead8641ef3f24b6b07366a1c34e5b059faf363
RUN savedAptMark="$(apt-mark showmanual)" \
 && apt-get update && apt-get install -y --no-install-recommends \
 && apt install -y wget \
 && wget -O/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 \
 && echo "$DUMB_INIT_SHA256  /dumb-init" | sha256sum -c -\
 && chmod +x /dumb-init \
 && mkdir -p /ingress-controller /etc/haproxy/maps \
 && apt-mark auto '.*' > /dev/null \
 && { [ -z "$savedAptMark" ] || apt-mark manual $savedAptMark; } \
 && find /usr/local -type f -executable -exec ldd '{}' ';' \
      | awk '/=>/ { print $(NF-1) }' \
      | sort -u \
      | xargs -r dpkg-query --search \
      | cut -d: -f1 \
      | sort -u \
      | xargs -r apt-mark manual \
 && apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false \
 && apt install -y socat lua5.3 lua-socket \
 && rm -rf /var/lib/apt/lists/*

COPY . /

STOPSIGNAL SIGTERM
ENTRYPOINT ["/dumb-init", "--", "/start.sh"]
