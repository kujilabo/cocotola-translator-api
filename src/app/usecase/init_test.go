package usecase_test

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
