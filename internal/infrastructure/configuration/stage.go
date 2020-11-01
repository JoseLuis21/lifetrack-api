package configuration

const (
	// ProductionStage used by clients, full prod data
	ProductionStage = "prod"
	// StagingStage used by QA engineers and clients for UAT, limited prod data
	StagingStage = "staging"
	// TestingStage used by QA engineers, mocked data
	TestingStage = "test"
	// DevelopmentStage used by developers, mocked data
	DevelopmentStage = "dev"
)
