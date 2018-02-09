package faketmpfile

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

func ExampleFakeTemp() {
	rdr := strings.NewReader("Hello, world!\n")
	tmpfile, err := FakeTemp(rdr)
	if err != nil {
		log.Fatal(err)
	}
	defer tmpfile.Close()
	cat, _ := exec.LookPath("cat")
	cmd := exec.Command(cat, tmpfile.Name())
	tmpfile.ExtraFiles(cmd)
	data, err := cmd.CombinedOutput()
	println(string(data))
}
func TestExec(t *testing.T) {
	rdr := strings.NewReader("Hello, world!\n")
	tmpfile, err := FakeTemp(rdr)
	if err != nil {
		t.Fatal(err)
	}
	defer tmpfile.Close()
	cat, _ := exec.LookPath("cat")
	cmd := exec.Command(cat, tmpfile.Name())
	tmpfile.ExtraFiles(cmd)
	data, err := cmd.CombinedOutput()
	if err != nil {
		// this fails when the tests run from within emacs
		// no idea why, inherited flags?
		t.Fatal(string(data) + ":" + err.Error())
	}
	if string(data) != "Hello, world!\n" {
		t.Fatalf("got %q", string(data))
	}
}
