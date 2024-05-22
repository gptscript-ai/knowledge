package flows

type FlowType string

const (
	FLowTypeIngest   FlowType = "ingest"
	FlowTypeRetrieve FlowType = "retrieve"
)

type FlowStep struct {
	StepType string
	StepName string
	NextStep *FlowStep
}

type Flow struct {
	Type  FlowType
	Steps *FlowStep
}
