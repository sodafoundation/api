# Copyright (c) 2018 Huawei Technologies Co., Ltd. All Rights Reserved.
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

IMAGE = opensdsio/dashboard
VERSION := latest

all:build
.PHONY: all

build:dashboard
.PHONY: build

dashboard:package
	chmod +x ./image_builder.sh \
	  && ./image_builder.sh
.PHONY: dashboard

package:
	apt-get update && apt-get install -y --no-install-recommends \
	  wget \
	  make \
	  g++ \
	  nginx \
	  && rm -rf /var/lib/apt/lists/* \
	  && apt-get clean
	wget --no-check-certificate https://deb.nodesource.com/setup_8.x \
	  && chmod +x setup_8.x && ./setup_8.x \
	  && apt-get install -y nodejs
.PHONY: package

docker:
	docker build . -t $(IMAGE):$(VERSION)
.PHONY: docker

clean:
	service nginx stop
	rm -rf /etc/nginx/sites-available/default /var/www/html/* ./dist warn=False
	npm uninstall --unsafe-perm
	npm uninstall --unsafe-perm -g @angular/cli@1.7.4
	apt-get --purge remove -y nodejs nginx \
	  && rm -rf ./setup_8.x warn=False
.PHONY: clean
