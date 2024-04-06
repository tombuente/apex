package sql

import _ "embed"

//go:embed logistics-schema.sql
var LogisticsSchema string

//go:embed logistics-fixture.sql
var LogisticsFixture string
