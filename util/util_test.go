package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/akutz/gotil"

	"github.com/emccode/polly/core/version"
)

var r10 string
var tmpPrefixDirs []string

func TestMain(m *testing.M) {
	r10 = gotil.RandomString(10)

	exitCode := m.Run()
	for _, d := range tmpPrefixDirs {
		os.RemoveAll(d)
	}
	os.Exit(exitCode)
}

func newPrefixDir(testName string, t *testing.T) string {
	tmpDir, err := ioutil.TempDir(
		"", fmt.Sprintf("polly-util_test-%s", testName))
	if err != nil {
		t.Fatal(err)
	}

	Prefix(tmpDir)
	os.MkdirAll(tmpDir, 0755)
	tmpPrefixDirs = append(tmpPrefixDirs, tmpDir)
	return tmpDir
}

func TestPrefix(t *testing.T) {
	if IsPrefixed() {
		t.Fatalf("is prefixed %s", GetPrefix())
	}
	Prefix("")
	if IsPrefixed() {
		t.Fatalf("is prefixed %s", GetPrefix())
	}
	Prefix("/")
	if IsPrefixed() {
		t.Fatalf("is prefixed %s", GetPrefix())
	}

	tmpDir := newPrefixDir("TestHomeDir", t)
	Prefix(tmpDir)
	if !IsPrefixed() {
		t.Fatalf("is not prefixed %s", GetPrefix())
	}

	p := GetPrefix()
	if p != tmpDir {
		t.Fatalf("prefix != %s, == %s", tmpDir, p)
	}
}

func TestPrefixAndDirs(t *testing.T) {
	tmpDir := newPrefixDir("TestPrefixAndDirs", t)

	etcDirPath := EtcDirPath()
	expEtcDirPath := fmt.Sprintf("%s/etc/polly", tmpDir)
	if etcDirPath != expEtcDirPath {
		t.Fatalf("EtcDirPath() == %s, != %s", etcDirPath, expEtcDirPath)
	}

	etcDirFilePath := EtcFilePath("etcFile")
	expEtcFilePath := fmt.Sprintf("%s/%s", etcDirPath, "etcFile")
	if expEtcFilePath != etcDirFilePath {
		t.Fatalf("EtcFilePath(\"etcFile\") == %s, != %s",
			etcDirFilePath, expEtcFilePath)
	}

	runDirPath := RunDirPath()
	expRunDirPath := fmt.Sprintf("%s/var/run/polly", tmpDir)
	if runDirPath != expRunDirPath {
		t.Fatalf("RunDirPath() == %s, != %s", runDirPath, expRunDirPath)
	}

	logDirPath := LogDirPath()
	expLogDirPath := fmt.Sprintf("%s/var/log/polly", tmpDir)
	if logDirPath != expLogDirPath {
		t.Fatalf("LogDirPath() == %s, != %s", logDirPath, expLogDirPath)
	}

	logDirFilePath := LogFilePath("logFile")
	expLogFilePath := fmt.Sprintf("%s/%s", logDirPath, "logFile")
	if expLogFilePath != logDirFilePath {
		t.Fatalf("LogFilePath(\"logFile\") == %s, != %s",
			logDirFilePath, expLogFilePath)
	}

	libDirPath := LibDirPath()
	expLibDirPath := fmt.Sprintf("%s/var/lib/polly", tmpDir)
	if libDirPath != expLibDirPath {
		t.Fatalf("LibDirPath() == %s, != %s", libDirPath, expLibDirPath)
	}

	libDirFilePath := LibFilePath("libFile")
	expLibFilePath := fmt.Sprintf("%s/%s", libDirPath, "libFile")
	if expLibFilePath != libDirFilePath {
		t.Fatalf("LibFilePath(\"libFile\") == %s, != %s",
			libDirFilePath, expLibFilePath)
	}

	binDirPath := BinDirPath()
	expBinDirPath := fmt.Sprintf("%s/usr/bin", tmpDir)
	if binDirPath != expBinDirPath {
		t.Fatalf("BinDirPath() == %s, != %s", binDirPath, expBinDirPath)
	}

	binDirFilePath := BinFilePath()
	expBinFilePath := fmt.Sprintf("%s/%s", binDirPath, "polly")
	if expBinFilePath != binDirFilePath {
		t.Fatalf("BinFilePath(\"polly\") == %s, != %s",
			binDirFilePath, expBinFilePath)
	}

	pidFilePath := PidFilePath()
	expPidFilePath := fmt.Sprintf("%s/var/run/polly/polly.pid", tmpDir)
	if expPidFilePath != pidFilePath {
		t.Fatalf("PidFilePath() == %s, != %s", pidFilePath, expPidFilePath)
	}
}

func TestStdOutAndLogFile(t *testing.T) {
	newPrefixDir("TestStdOutAndLogFile", t)

	if _, err := StdOutAndLogFile("BadFile/"); err == nil {
		t.Fatal("error expected in created BadFile")
	}

	out, err := StdOutAndLogFile("TestStdOutAndLogFile")

	if err != nil {
		t.Fatal(err)
	}

	if out == nil {
		t.Fatal("out == nil")
	}
}

func TestWriteReadCurrentPidFile(t *testing.T) {
	newPrefixDir("TestWriteReadPidFile", t)

	var err error
	var pidRead int

	pid := os.Getpid()

	if err = WritePidFile(-1); err != nil {
		t.Fatalf("error writing pidfile=%s", PidFilePath())
	}

	if pidRead, err = ReadPidFile(); err != nil {
		t.Fatalf("error reading pidfile=%s", PidFilePath())
	}

	if pidRead != pid {
		t.Fatalf("pidRead=%d != pid=%d", pidRead, pid)
	}
}

func TestWriteReadCustomPidFile(t *testing.T) {
	newPrefixDir("TestWriteReadPidFile", t)

	var err error
	if _, err = ReadPidFile(); err == nil {
		t.Fatal("error expected in reading pid file")
	}

	pidWritten := int(time.Now().Unix())
	if err = WritePidFile(pidWritten); err != nil {
		t.Fatalf("error writing pidfile=%s", PidFilePath())
	}

	var pidRead int
	if pidRead, err = ReadPidFile(); err != nil {
		t.Fatalf("error reading pidfile=%s", PidFilePath())
	}

	if pidRead != pidWritten {
		t.Fatalf("pidRead=%d != pidWritten=%d", pidRead, pidWritten)
	}
}

func TestReadPidFileWithErrors(t *testing.T) {
	newPrefixDir("TestWriteReadPidFile", t)

	var err error
	if _, err = ReadPidFile(); err == nil {
		t.Fatal("error expected in reading pid file")
	}

	gotil.WriteStringToFile("hello", PidFilePath())

	if _, err = ReadPidFile(); err == nil {
		t.Fatal("error expected in reading pid file")
	}
}

func TestPrintVersion(t *testing.T) {
	version.Arch = "Linux-x86_64"
	version.Branch = "master"
	version.ShaLong = gotil.RandomString(32)
	version.Epoch = fmt.Sprintf("%d", time.Now().Unix())
	version.SemVer = "1.0.0"
	_, _, thisAbsPath := gotil.GetThisPathParts()
	epochStr := version.EpochToRfc1123()

	t.Logf("thisAbsPath=%s", thisAbsPath)
	t.Logf("epochStr=%s", epochStr)

	var buff []byte
	b := bytes.NewBuffer(buff)

	PrintVersion(b)

	vs := b.String()

	evs := `Binary: ` + thisAbsPath + `
SemVer: ` + version.SemVer + `
OsArch: ` + version.Arch + `
Branch: ` + version.Branch + `
Commit: ` + version.ShaLong + `
Formed: ` + epochStr + `
`

	if vs != evs {
		t.Fatalf("nexpectedVersionString=%s\n\nversionString=%s\n", evs, vs)
	}
}

func TestInstall(t *testing.T) {
	Install()
}

func TestInstallChownRoot(t *testing.T) {
	InstallChownRoot()
}

func TestInstallDirChownRoot(t *testing.T) {
	InstallDirChownRoot("--help")
}
