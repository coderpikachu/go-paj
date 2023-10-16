//go:build e2e

package integration

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"test/12paj/my_orm"
)

type Suite struct {
	suite.Suite

	driver string
	dsn    string

	db *my_orm.DB
}

func (i *Suite) SetupSuite() {
	db, err := my_orm.Open(i.driver, i.dsn)
	require.NoError(i.T(), err)
	err = db.Wait()
	require.NoError(i.T(), err)
	i.db = db
}
