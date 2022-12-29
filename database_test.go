package subscribe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFile  = filepath.Join(os.TempDir(), "this_is_a_testfile_for_subtscribe_test.go.json")
	testFile2 = filepath.Join(os.TempDir(), "this_is_a_testfile_for_subtscribe_test2.go.json")
	testFile4 = filepath.Join(os.TempDir(), "this_is_a_testfile_for_subtscribe_test4.go.json")
)

func TestGetDB(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	sub, err := GetDB("")
	assert.Nil(err, "getting an empty db must produce no error")

	json, err := sub.StateGetJSON()
	assert.EqualValues(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, json,
		"the initial state must be empty")
	assert.Nil(err, "getting an empty state must produce no error")
}

func TestStateFileLoad(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	// test with good data.
	testJSON := `{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[{"id":0,"meta":null,"api":` +
		`"http","contact":"testUser","events":{"eventsMap":{}},"isAdmin":false,"ignored":false}]}`
	assert.Nil(os.WriteFile(testFile, []byte(testJSON), 0o600), "problem writing test file")

	sub, err := GetDB(testFile)
	assert.Nil(err, "there must be no error loading the state file")

	json, err := sub.StateGetJSON()
	assert.Nil(err, "there must be no error getting the state data")
	assert.EqualValues(testJSON, json)

	// Test missing file.
	assert.Nil(os.RemoveAll(testFile), "problem removing test file")

	_, err = GetDB(testFile)
	assert.Nil(err, "there must be no error when the state file is missing")

	data, err := os.ReadFile(testFile)
	assert.Nil(err, "error reading test file")

	assert.EqualValues(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, data,
		"the initial state file must be empty")

	// Test uncreatable file.
	_, err = GetDB("/tmp/xxx/yyy/zzz/aaa/bbb/this_file_dont_exist")
	assert.NotNil(err, "there must be an error when the state cannot be created")

	// Test unreadable file.
	_, err = GetDB("/etc/sudoers")
	assert.NotNil(err, "there must be an error when the state cannot be read")

	// Test bad data.
	err = os.WriteFile(testFile, []byte("this aint good json}}"), 0o600)
	assert.Nil(err, "problem writing test file")

	_, err = GetDB(testFile)
	assert.NotNil(err, "there must be an error when the state file is corrupt")
}

func TestStateFileSave(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	assert.Nil(os.RemoveAll(testFile2), "problem removing test file")
	sub, err := GetDB(testFile2)
	assert.Nil(err, "there must be no error creating the initial state file")
	assert.Nil(sub.StateFileSave(), "there must be no error saving the state file")
	sub, err = GetDB("")
	assert.Nil(err, "there must be no error when the state file does not exist")
	assert.Nil(sub.StateFileSave(), "there must be no error when the state file does not exist")
}

func TestStateFileRelocate(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	assert.Nil(os.RemoveAll(testFile4), "problem removing test file")

	sub, err := GetDB("")
	assert.Nil(err, "there must be no error when the state file does not exist")
	assert.Nil(sub.StateFileSave(), "there must be no error when the state file does not exist")

	err = sub.StateFileRelocate(testFile4)
	assert.EqualValues(testFile4, sub.stateFile, "the path was not to the new value")
	assert.Nil(err, "there must be no error creating the initial state file")

	err = sub.StateFileRelocate(os.TempDir())
	assert.NotNil(err, "there must be an error trying to write a file as a tmp folder")
	assert.EqualValues(testFile4, sub.stateFile, "the path was not changed back to the previous value")
}
