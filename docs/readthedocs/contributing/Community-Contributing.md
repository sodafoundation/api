# Community-Contributing
## OpenSDS

[![Go Report Card](https://goreportcard.com/badge/github.com/sodafoundation/api?branch=master)](https://goreportcard.com/report/github.com/sodafoundation/api)
[![Build Status](https://travis-ci.org/sodafoundation/api.svg?branch=master)](https://travis-ci.org/sodafoundation/api)
[![Coverage Status](https://coveralls.io/repos/github/sodafoundation/api/badge.svg?branch=master)](https://coveralls.io/github/sodafoundation/api?branch=master)

<img src="https://www.opensds.io/wp-content/uploads/sites/18/2016/11/logo_opensds.png" width="100">


## How to contribute

opensds is Apache 2.0 licensed and accepts contributions via GitHub pull requests. This document outlines some of the conventions on commit message formatting, contact points for developers and other resources to make getting your contribution into opensds easier.

## Email and chat

- Email: [opensds-dev](https://groups.google.com/forum/?hl=en#!forum/opensds-dev)
- Slack: #[opensds](https://opensds.slack.com) 

Before you start, NOTICE that ```master``` branch is the relatively stable version
provided for customers and users. So all code modifications SHOULD be submitted to
```development``` branch.

## Getting started

- Fork the repository on GitHub.
- Read the README.md and INSTALL.md for project information and build instructions.

For those who just get in touch with this project recently, here is a proposed contributing [tutorial](https://github.com/leonwanghui/installation-note/blob/master/opensds_fork_contribute_tutorial.md).

## Contribution Workflow

### Code style

The coding style suggested by the Golang community is used in opensds. See the [doc](https://github.com/golang/go/wiki/CodeReviewComments) for more details.

Please follow this style to make opensds easy to review, maintain and develop.

### Report issues

A great way to contribute to the project is to send a detailed report when you encounter an issue. We always appreciate a well-written, thorough bug report, and will thank you for it!

When reporting issues, refer to this format:

- What version of env (opensds, os, golang etc) are you using?
- Is this a BUG REPORT or FEATURE REQUEST?
- What happened?
- What you expected to happen?
- How to reproduce it?(as minimally and precisely as possible)

### Propose PRs

- Raise your idea as an [issue](https://github.com/sodafoundation/api/issues)
- If it is a new feature that needs lots of design details, a design proposal should also be submitted [here](https://github.com/opensds/design-specs/pulls).
- After reaching consensus in the issue discussions and design proposal reviews, complete the development on the forked repo and submit a PR. 
  Here are the [PRs](https://github.com/sodafoundation/api/pulls?q=is%3Apr+is%3Aclosed) that are already closed.
- If a PR is submitted by one of the core members, it has to be merged by a different core member.
- After PR is sufficiently discussed, it will get merged, abondoned or rejected depending on the outcome of the discussion.

Thank you for your contribution !
