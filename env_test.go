package forge_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/fuzzingbits/forge"
)

type TestEnvStructChild struct {
	Bar       bool    `env:"FORGE_CONFIG_TEST_BAR"`
	PtrString *string `env:"FORGE_CONFIG_TEST_PTRSTRING"`
}

type TestEnvStruct struct {
	Name   string   `env:"FORGE_CONFIG_TEST_NAME" json:"name"`
	Skills []string `env:"FORGE_CONFIG_TEST_SKILLS"`
	Age    int      `env:"FORGE_CONFIG_TEST_AGE" json:"age"`
	PtrInt *int     `env:"FORGE_CONFIG_TEST_PTRINT"`
	Foo    TestEnvStructChild
}

type TestEnvUnexported struct {
	name string `env:"FORGE_CONFIG_TEST_NAME"`
}

func resetTest() {
	os.Unsetenv("FORGE_CONFIG_TEST_BAR")
	os.Unsetenv("FORGE_CONFIG_TEST_PTRSTRING")
	os.Unsetenv("FORGE_CONFIG_TEST_NAME")
	os.Unsetenv("FORGE_CONFIG_TEST_SKILLS")
	os.Unsetenv("FORGE_CONFIG_TEST_AGE")
	os.Unsetenv("FORGE_CONFIG_TEST_PTRINT")
}

func envTestHelper(t *testing.T, startingConfig interface{}, targetConfig interface{}, targetErr error, justLookForAnyError bool) {
	defer resetTest()

	err := forge.ParseEnvironment(startingConfig)
	if justLookForAnyError {
		if err == nil {
			t.Errorf("No error was found but one was expected")
			return
		}
	} else {
		if err != targetErr {
			t.Errorf("error was not correct, got: \"%v\", want: \"%v\"", err, targetErr)
			return
		}
	}

	if !reflect.DeepEqual(startingConfig, targetConfig) {
		t.Errorf("target config did not match the starting config. got: %+v, want: %+v", startingConfig, targetConfig)
		return
	}
}

func TestProviderEnvironment(t *testing.T) {
	var stringPointerExample = new(string)
	*stringPointerExample = "foobar2"

	startingConfig := TestEnvStruct{
		Name: "Aaron",
		Age:  22,
	}

	targetConfig := TestEnvStruct{
		Name: "George",
		Age:  42,
		Foo: TestEnvStructChild{
			Bar:       true,
			PtrString: stringPointerExample,
		},
	}

	os.Setenv("FORGE_CONFIG_TEST_NAME", "George")
	os.Setenv("FORGE_CONFIG_TEST_AGE", "42")
	os.Setenv("FORGE_CONFIG_TEST_BAR", "true")
	os.Setenv("FORGE_CONFIG_TEST_PTRSTRING", *stringPointerExample)

	envTestHelper(t, &startingConfig, &targetConfig, nil, false)
}

func TestProviderEnvironmentInvalidInt(t *testing.T) {
	startingConfig := TestEnvStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_BAR", "not a valid bool")

	envTestHelper(t, &startingConfig, &startingConfig, nil, true)
}

func TestProviderEnvironmentInvalidBool(t *testing.T) {
	startingConfig := TestEnvStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_AGE", "not a valid int")

	envTestHelper(t, &startingConfig, &startingConfig, nil, true)
}

func TestProviderEnvironmentErrUnexportedField(t *testing.T) {
	startingConfig := TestEnvUnexported{
		name: "Aaron",
	}

	os.Setenv("FORGE_CONFIG_TEST_NAME", "George")

	envTestHelper(t, &startingConfig, &startingConfig, forge.ErrUnexportedField, false)
}

func TestProviderEnvironmentErrUnsupportedType(t *testing.T) {
	startingConfig := TestEnvStruct{
		Skills: []string{"go"},
	}

	os.Setenv("FORGE_CONFIG_TEST_SKILLS", "go")

	envTestHelper(t, &startingConfig, &startingConfig, forge.ErrUnsupportedType, false)
}

func TestProviderEnvironmentPointerSetError(t *testing.T) {
	var intPointerExample = new(int)
	*intPointerExample = 22

	startingConfig := TestEnvStruct{
		Name: "Aaron",
		Age:  22,
	}

	targetConfig := TestEnvStruct{
		Name: "Aaron",
		Age:  22,
	}

	os.Setenv("FORGE_CONFIG_TEST_PTRINT", "not an int")

	envTestHelper(t, &startingConfig, &targetConfig, nil, true)
}

func TestProviderEnvironmentNotPointer(t *testing.T) {
	envTestHelper(t, TestEnvStruct{}, TestEnvStruct{}, nil, true)
}

func TestReadDotEnv(t *testing.T) {
	targetValues := map[string]string{
		"FOOBAR1": "SYSTEM",
		"FOOBAR2": "DEFAULT2",
		"FOOBAR3": "LOCAL3",
	}

	dotEnvTestHelper(t, "basic", targetValues, func() {
		os.Setenv("FOOBAR1", "SYSTEM")
	})
}

func TestReadDotEnvNoFiles(t *testing.T) {
	targetValues := map[string]string{
		"FOOBAR1": "SYSTEM",
		"FOOBAR2": "",
		"FOOBAR3": "",
	}

	dotEnvTestHelper(t, "fake_directory", targetValues, func() {
		os.Setenv("FOOBAR1", "SYSTEM")
	})
}

func dotEnvTestHelper(t *testing.T, directory string, targetValues map[string]string, setup func()) {
	// Fix directory
	wd, _ := os.Getwd()
	os.Chdir("./test_files/dotenv_tests/" + directory)
	defer os.Chdir(wd)

	// Reset environment
	for key := range targetValues {
		os.Unsetenv(key)
	}

	// Set environment state
	setup()

	// Read in the values
	forge.ReadDotEnv()

	// Confirm results
	for key, targetValue := range targetValues {
		actualValue := os.Getenv(key)
		if targetValue != actualValue {
			t.Errorf("ReadDotEnv() key: %s = %s ; want: %s", key, actualValue, targetValue)
		}
	}
}
