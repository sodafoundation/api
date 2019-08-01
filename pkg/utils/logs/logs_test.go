package logs

import "testing"

func TestLogContent(t *testing.T) {
	logPrefix := "hotpot"
	pid := "11111"
	logLevel := "ERROR"
	filename := "test.go"
	line := "20"
	funcname := "testFunc"
	logMsg := "Hello, World"
	error := "404"
	logContent := GetLogContent(logPrefix, pid, logLevel, filename, line, funcname, logMsg, error)
	expected := "hotpot [ pid: 11111 ] [ test.go : 20 testFunc ] Hello, World"
	t.Logf(logContent)
	if logContent[56:] != expected {
		t.Errorf("Expected %s\n got %s\n", expected, logContent[56:])
	}
}
