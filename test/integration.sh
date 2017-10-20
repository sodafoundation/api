sudo docker build integration/osdsdock/ -t opensds/osdsdock:integration-test
sudo docker run -d --net=host opensds/osdsdock:integration-test

runTests() {
  go test -v github.com/opensds/opensds/test/integration
}

runTests
