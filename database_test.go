package subscribe

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFile = "/tmp/this_is_a_testfile_for_subtscribe_test.go.json"
var testFile2 = "/tmp/this_is_a_testfile_for_subtscribe_test2.go.json"

func TestGetDB(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	sub, err := GetDB("")
	a.Nil(err, "getting an empty db must produce no error")
	json, err := sub.StateGetJSON()
	a.EqualValues(`{"enabled_apis":[],"events":{"events_map":{}},"subscribers":[]}`, string(json), "the initial state must be empty")
	a.Nil(err, "getting an empty state must produce no error")
}

func TestLoadStateFile(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	// test with good data.
	testJSON := `{"enabled_apis":[],"events":{"events_map":{}},"subscribers":[{"api":"http","contact":"testUser","events":{"events_map":{}},"is_admin":false,"ignored":false}]}`
	a.Nil(ioutil.WriteFile(testFile, []byte(testJSON), 0644), "problem writing test file")
	sub, err := GetDB(testFile)
	a.Nil(err, "there must be no error loading the state file")
	json, err := sub.StateGetJSON()
	a.Nil(err, "there must be no error getting the state data")
	a.EqualValues(testJSON, string(json))
	// Test missing file.
	a.Nil(os.RemoveAll(testFile), "problem removing test file")
	_, err = GetDB(testFile)
	a.Nil(err, "there must be no error when the state file is missing")
	data, err := ioutil.ReadFile(testFile)
	a.Nil(err, "error reading test file")
	a.EqualValues(`{"enabled_apis":[],"events":{"events_map":{}},"subscribers":[]}`, string(data), "the initial state file must be empty")
	// Test uncreatable file.
	_, err = GetDB("/tmp/xxx/yyy/zzz/aaa/bbb/this_file_dont_exist")
	a.NotNil(err, "there must be an error when the state cannot be created")
	// Test unreadable file.
	_, err = GetDB("/etc/sudoers")
	a.NotNil(err, "there must be an error when the state cannot be read")
	// Test bad data.
	err = ioutil.WriteFile(testFile, []byte("this aint good json}}"), 0644)
	a.Nil(err, "problem writing test file")
	_, err = GetDB(testFile)
	a.NotNil(err, "there must be an error when the state file is corrupt")
}

func TestSaveStateFile(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.Nil(os.RemoveAll(testFile2), "problem removing test file")
	sub, err := GetDB(testFile2)
	a.Nil(err, "there must be no error creating the initial state file")
	a.Nil(sub.StateFileSave(), "there must be no error saving the state file")
	sub, err = GetDB("")
	a.Nil(err, "there must be no error when the state file does not exist")
	a.Nil(sub.StateFileSave(), "there must be no error when the state file does not exist")
}
