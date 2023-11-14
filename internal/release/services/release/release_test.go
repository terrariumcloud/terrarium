package release

import (
	"errors"
	"testing"

	"github.com/terrariumcloud/terrarium/internal/storage/mocks"

	"google.golang.org/grpc"
)

// Test_RegisterReleaseWithServer checks:
// - if there was no error with table init
// - if error is returned when Table initialization fails
func Test_RegisterReleaseWithServer(t *testing.T) {
	t.Parallel()

	t.Run("when there is no error with table init", func(t *testing.T) {
		db := &mocks.DynamoDB{}

		rs := &ReleaseService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != nil {
			t.Errorf("Expected no error, got %v.", err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected no calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})

	t.Run("when Table initialization fails", func(t *testing.T) {

		db := &mocks.DynamoDB{
			DescribeTableErrors: []error{errors.New("some error")},
		}

		rs := &ReleaseService{
			Db: db,
		}

		s := grpc.NewServer(*new([]grpc.ServerOption)...)

		err := rs.RegisterWithServer(s)

		if err != ReleaseTableInitializationError {
			t.Errorf("Expected %v, got %v.", ReleaseTableInitializationError, err)
		}

		if db.DescribeTableInvocations != 1 {
			t.Errorf("Expected 1 call to DescribeTable, got %v.", db.DescribeTableInvocations)
		}

		if db.CreateTableInvocations != 0 {
			t.Errorf("Expected 0 calls to CreateTable, got %v.", db.CreateTableInvocations)
		}
	})
}
