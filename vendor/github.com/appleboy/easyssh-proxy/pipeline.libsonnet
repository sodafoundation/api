{
  test:: {
    kind: 'pipeline',
    name: 'testing',
    platform: {
      os: 'linux',
      arch: 'amd64',
    },
    steps: [
      {
        name: 'vet',
        image: 'golang:1.12',
        pull: 'always',
        environment: {
          GO111MODULE: 'on',
        },
        commands: [
          'make vet',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'lint',
        image: 'golang:1.12',
        pull: 'always',
        environment: {
          GO111MODULE: 'on',
        },
        commands: [
          'make lint',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'misspell',
        image: 'golang:1.12',
        pull: 'always',
        environment: {
          GO111MODULE: 'on',
        },
        commands: [
          'make misspell-check',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'embedmd',
        image: 'golang:1.12',
        pull: 'always',
        environment: {
          GO111MODULE: 'on',
        },
        commands: [
          'make embedmd',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'test',
        image: 'golang:1.12-alpine',
        pull: 'always',
        environment: {
          GO111MODULE: 'on',
        },
        commands: [
          'apk add git make curl perl bash build-base zlib-dev ucl-dev',
          'make ssh-server',
          'make test',
        ],
        volumes: [
          {
            name: 'gopath',
            path: '/go',
          },
        ],
      },
      {
        name: 'codecov',
        image: 'robertstettner/drone-codecov',
        pull: 'always',
        settings: {
          token: { 'from_secret': 'codecov_token' },
        },
      },
    ],
    volumes: [
      {
        name: 'gopath',
        temp: {},
      },
    ],
  }
}
