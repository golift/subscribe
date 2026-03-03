package subscribe

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)

	sub, err := GetDB("")
	require.NoError(t, err, "getting an empty db must produce no error")

	json, err := sub.StateGetJSON()
	assertions.JSONEq(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, json,
		"the initial state must be empty")
	require.NoError(t, err, "getting an empty state must produce no error")
}

func TestStateFileLoad(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	testFile := filepath.Join(t.TempDir(), "state.json")

	// test with good data.
	testJSON := `{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[{"id":0,"meta":null,"api":` +
		`"http","contact":"testUser","events":{"eventsMap":{}},"isAdmin":false,"ignored":false}]}`
	require.NoError(t, os.WriteFile(testFile, []byte(testJSON), 0o600), "problem writing test file")

	sub, err := GetDB(testFile)
	require.NoError(t, err, "there must be no error loading the state file")

	json, err := sub.StateGetJSON()
	require.NoError(t, err, "there must be no error getting the state data")
	assertions.JSONEq(testJSON, json)

	// Test missing file.
	require.NoError(t, os.RemoveAll(testFile), "problem removing test file")

	_, err = GetDB(testFile)
	require.NoError(t, err, "there must be no error when the state file is missing")

	// #nosec G304 -- test controls this temporary file path.
	data, err := os.ReadFile(testFile)
	require.NoError(t, err, "error reading test file")

	assertions.JSONEq(`{"enabledApis":[],"events":{"eventsMap":{}},"subscribers":[]}`, string(data),
		"the initial state file must be empty")

	// Test uncreatable file.
	_, err = GetDB("/tmp/xxx/yyy/zzz/aaa/bbb/this_file_dont_exist")
	require.Error(t, err, "there must be an error when the state cannot be created")

	// Test unreadable path (directory).
	_, err = GetDB(t.TempDir())
	require.Error(t, err, "there must be an error when the state path is not readable as a file")

	// Test bad data.
	err = os.WriteFile(testFile, []byte("this aint good json}}"), 0o600)
	require.NoError(t, err, "problem writing test file")

	_, err = GetDB(testFile)
	require.Error(t, err, "there must be an error when the state file is corrupt")
}

func TestStateFileSave(t *testing.T) {
	t.Parallel()

	testFile := filepath.Join(t.TempDir(), "state-save.json")
	sub, err := GetDB(testFile)
	require.NoError(t, err, "there must be no error creating the initial state file")
	require.NoError(t, sub.StateFileSave(), "there must be no error saving the state file")
	sub, err = GetDB("")
	require.NoError(t, err, "there must be no error when the state file does not exist")
	require.NoError(t, sub.StateFileSave(), "there must be no error when the state file does not exist")
}

func TestStateFileRelocate(t *testing.T) {
	t.Parallel()

	assertions := assert.New(t)
	testFile4 := filepath.Join(t.TempDir(), "state-relocate.json")

	sub, err := GetDB("")
	require.NoError(t, err, "there must be no error when the state file does not exist")
	require.NoError(t, sub.StateFileSave(), "there must be no error when the state file does not exist")

	err = sub.StateFileRelocate(testFile4)
	assertions.Equal(testFile4, sub.stateFile, "the path was not to the new value")
	require.NoError(t, err, "there must be no error creating the initial state file")

	err = sub.StateFileRelocate(t.TempDir())
	require.Error(t, err, "there must be an error trying to write a file as a folder")
	assertions.Equal(testFile4, sub.stateFile, "the path was not changed back to the previous value")
}

func TestStateGetJSONMarshalError(t *testing.T) {
	t.Parallel()

	sub := &Subscribe{
		Events:      &Events{Map: make(map[string]*Rules)},
		Subscribers: make([]*Subscriber, 0),
	}

	subscriber := sub.CreateSub("marshal", "api", false, false)
	subscriber.Meta = map[string]any{"bad": make(chan int)}

	_, err := sub.StateGetJSON()
	require.Error(t, err, "unsupported meta values should fail marshaling")
}

func TestStateFileSaveMarshalError(t *testing.T) {
	t.Parallel()

	stateFile := filepath.Join(t.TempDir(), "state-save-marshal.json")
	sub, err := GetDB(stateFile)
	require.NoError(t, err)

	subscriber := sub.CreateSub("marshal", "api", false, false)
	subscriber.Meta = map[string]any{"bad": make(chan int)}

	err = sub.StateFileSave()
	require.Error(t, err)
	assert.ErrorContains(t, err, "marshaling json")
}
