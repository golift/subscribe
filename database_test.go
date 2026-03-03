package subscribe

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testFile  = "/tmp/this_is_a_testfile_for_subtscribe_test.go.json"
	testFile2 = "/tmp/this_is_a_testfile_for_subtscribe_test2.go.json"
	testFile4 = "/tmp/this_is_a_testfile_for_subtscribe_test4.go.json"
)

func TestGetDB(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)

	sub, err := GetDB("")
	require.NoError(t, err, "getting an empty db must produce no error")

	json, err := sub.StateGetJSON()
	asert.JSONEq(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, json,
		"the initial state must be empty")
	require.NoError(t, err, "getting an empty state must produce no error")
}

func TestStateFileLoad(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)

	// test with good data.
	testJSON := `{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[{"id":0,"meta":null,"api":` +
		`"http","contact":"testUser","events":{"eventsMap":{}},"isAdmin":false,"ignored":false}]}`
	require.NoError(t, os.WriteFile(testFile, []byte(testJSON), 0o600), "problem writing test file")

	sub, err := GetDB(testFile)
	require.NoError(t, err, "there must be no error loading the state file")

	json, err := sub.StateGetJSON()
	require.NoError(t, err, "there must be no error getting the state data")
	asert.JSONEq(testJSON, json)

	// Test missing file.
	require.NoError(t, os.RemoveAll(testFile), "problem removing test file")

	_, err = GetDB(testFile)
	require.NoError(t, err, "there must be no error when the state file is missing")

	data, err := os.ReadFile(testFile)
	require.NoError(t, err, "error reading test file")

	asert.JSONEq(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, string(data),
		"the initial state file must be empty")

	// Test uncreatable file.
	_, err = GetDB("/tmp/xxx/yyy/zzz/aaa/bbb/this_file_dont_exist")
	require.Error(t, err, "there must be an error when the state cannot be created")

	// Test unreadable file.
	_, err = GetDB("/etc/sudoers")
	require.Error(t, err, "there must be an error when the state cannot be read")

	// Test bad data.
	err = os.WriteFile(testFile, []byte("this aint good json}}"), 0o600)
	require.NoError(t, err, "problem writing test file")

	_, err = GetDB(testFile)
	require.Error(t, err, "there must be an error when the state file is corrupt")
}

func TestStateFileSave(t *testing.T) {
	t.Parallel()
	require.NoError(t, os.RemoveAll(testFile2), "problem removing test file")
	sub, err := GetDB(testFile2)
	require.NoError(t, err, "there must be no error creating the initial state file")
	require.NoError(t, sub.StateFileSave(), "there must be no error saving the state file")
	sub, err = GetDB("")
	require.NoError(t, err, "there must be no error when the state file does not exist")
	require.NoError(t, sub.StateFileSave(), "there must be no error when the state file does not exist")
}

func TestStateFileRelocate(t *testing.T) {
	t.Parallel()

	asert := assert.New(t)
	require.NoError(t, os.RemoveAll(testFile4), "problem removing test file")

	sub, err := GetDB("")
	require.NoError(t, err, "there must be no error when the state file does not exist")
	require.NoError(t, sub.StateFileSave(), "there must be no error when the state file does not exist")

	err = sub.StateFileRelocate(testFile4)
	asert.Equal(testFile4, sub.stateFile, "the path was not to the new value")
	require.NoError(t, err, "there must be no error creating the initial state file")

	err = sub.StateFileRelocate("/tmp")
	require.Error(t, err, "there must be an error trying to write asert file as asert /tmp folder")
	asert.Equal(testFile4, sub.stateFile, "the path was not changed back to the previous value")
}
