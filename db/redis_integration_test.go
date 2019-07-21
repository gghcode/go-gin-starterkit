package db_test

import (
	"testing"

	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type redisIntegration struct {
	suite.Suite

	conf config.Configuration
}

func (suite *redisIntegration) SetupSuite() {
	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	require.NoError(suite.T(), err)

	suite.conf = conf
}

func TestRedisConnIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(redisIntegration))
}

func (suite *redisIntegration) TestNewRedisConn() {
	conn := db.NewRedisConn(suite.conf)
	pong, err := conn.Client().Ping().Result()

	assert.Equal(suite.T(), pong, "PONG")
	assert.Equal(suite.T(), err, nil)
}
