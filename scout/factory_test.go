package scout

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testScoutFactory struct{}

func (f *testScoutFactory) Model() string {
	return ""
}

func (f *testScoutFactory) Register(ScoutDelegate) error {
	return nil
}

func (f *testScoutFactory) UnRegister(string) {}

func (f *testScoutFactory) Done() {}

func testScoutFactoryCreator(conf *ModelConfig) (CheckModel, error) {
	return &testScoutFactory{}, nil
}

func TestCreateTestFactory(t *testing.T) {
	name := "test"
	Register(name, testScoutFactoryCreator)

	if m, err := CheckModelCreate(name, &ModelConfig{
		Data: nil,
	}); err != nil {
		t.Error(err)
		require.Equal(t, m.Model(), name, "")
	}
}
