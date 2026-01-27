package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chris-alexander-pop/system-design-library/pkg/workflow"
	"github.com/chris-alexander-pop/system-design-library/pkg/workflow/adapters/memory"
	"github.com/chris-alexander-pop/system-design-library/pkg/workflow/saga"
	"github.com/chris-alexander-pop/system-design-library/pkg/workflow/scheduler"
	"github.com/stretchr/testify/suite"
)

// WorkflowEngineSuite tests WorkflowEngine implementations.
type WorkflowEngineSuite struct {
	suite.Suite
	engine workflow.WorkflowEngine
	ctx    context.Context
}

func (s *WorkflowEngineSuite) SetupTest() {
	s.engine = memory.New()
	s.ctx = context.Background()
}

func (s *WorkflowEngineSuite) TestRegisterAndGetWorkflow() {
	def := workflow.WorkflowDefinition{
		ID:   "order-workflow",
		Name: "Order Processing",
		States: []workflow.State{
			{Name: "validate", Type: "Task", Next: "process"},
			{Name: "process", Type: "Task", End: true},
		},
		StartAt: "validate",
	}

	err := s.engine.RegisterWorkflow(s.ctx, def)
	s.Require().NoError(err)

	got, err := s.engine.GetWorkflow(s.ctx, "order-workflow")
	s.Require().NoError(err)
	s.Equal("Order Processing", got.Name)
}

func (s *WorkflowEngineSuite) TestGetWorkflowNotFound() {
	_, err := s.engine.GetWorkflow(s.ctx, "nonexistent")
	s.Error(err)
}

func (s *WorkflowEngineSuite) TestStartExecution() {
	s.engine.RegisterWorkflow(s.ctx, workflow.WorkflowDefinition{ID: "test-wf"})

	exec, err := s.engine.Start(s.ctx, workflow.StartOptions{
		WorkflowID: "test-wf",
		Input:      map[string]string{"order": "123"},
	})
	s.Require().NoError(err)
	s.NotEmpty(exec.ID)
	s.Equal("test-wf", exec.WorkflowID)
}

func (s *WorkflowEngineSuite) TestStartExecutionWorkflowNotFound() {
	_, err := s.engine.Start(s.ctx, workflow.StartOptions{WorkflowID: "nonexistent"})
	s.Error(err)
}

func (s *WorkflowEngineSuite) TestGetExecution() {
	s.engine.RegisterWorkflow(s.ctx, workflow.WorkflowDefinition{ID: "test-wf"})
	exec, _ := s.engine.Start(s.ctx, workflow.StartOptions{WorkflowID: "test-wf"})

	got, err := s.engine.GetExecution(s.ctx, exec.ID)
	s.Require().NoError(err)
	s.Equal(exec.ID, got.ID)
}

func (s *WorkflowEngineSuite) TestListExecutions() {
	s.engine.RegisterWorkflow(s.ctx, workflow.WorkflowDefinition{ID: "test-wf"})

	for i := 0; i < 3; i++ {
		s.engine.Start(s.ctx, workflow.StartOptions{WorkflowID: "test-wf"})
	}

	result, err := s.engine.ListExecutions(s.ctx, workflow.ListOptions{})
	s.Require().NoError(err)
	s.Len(result.Executions, 3)
}

func (s *WorkflowEngineSuite) TestCancelExecution() {
	s.engine.RegisterWorkflow(s.ctx, workflow.WorkflowDefinition{ID: "test-wf"})
	exec, _ := s.engine.Start(s.ctx, workflow.StartOptions{WorkflowID: "test-wf"})

	err := s.engine.Cancel(s.ctx, exec.ID)
	s.Require().NoError(err)

	exec, _ = s.engine.GetExecution(s.ctx, exec.ID)
	s.Equal(workflow.StatusCancelled, exec.Status)
}

func (s *WorkflowEngineSuite) TestWaitForCompletion() {
	s.engine.RegisterWorkflow(s.ctx, workflow.WorkflowDefinition{ID: "test-wf"})
	exec, _ := s.engine.Start(s.ctx, workflow.StartOptions{WorkflowID: "test-wf"})

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	result, err := s.engine.Wait(ctx, exec.ID)
	s.Require().NoError(err)
	s.Equal(workflow.StatusCompleted, result.Status)
}

// SagaSuite tests the Saga pattern.
type SagaSuite struct {
	suite.Suite
	ctx context.Context
}

func (s *SagaSuite) SetupTest() {
	s.ctx = context.Background()
}

func (s *SagaSuite) TestSagaSuccess() {
	var steps []string

	orderSaga := saga.New("order-saga").
		AddStep(saga.Step{
			Name: "reserve-inventory",
			Action: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "reserve")
				return data, nil
			},
			Compensate: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "release")
				return nil, nil
			},
		}).
		AddStep(saga.Step{
			Name: "charge-payment",
			Action: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "charge")
				return data, nil
			},
			Compensate: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "refund")
				return nil, nil
			},
		})

	exec, err := orderSaga.Execute(s.ctx, "order-123")
	s.Require().NoError(err)
	s.Equal(saga.StatusCompleted, exec.Status)
	s.Equal([]string{"reserve", "charge"}, steps)
}

func (s *SagaSuite) TestSagaCompensation() {
	var steps []string

	orderSaga := saga.New("order-saga").
		AddStep(saga.Step{
			Name: "reserve-inventory",
			Action: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "reserve")
				return data, nil
			},
			Compensate: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "release")
				return nil, nil
			},
		}).
		AddStep(saga.Step{
			Name: "charge-payment",
			Action: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "charge")
				return nil, errors.New("payment failed")
			},
			Compensate: func(ctx context.Context, data interface{}) (interface{}, error) {
				steps = append(steps, "refund")
				return nil, nil
			},
		})

	exec, err := orderSaga.Execute(s.ctx, "order-123")
	s.Error(err)
	s.Equal(saga.StatusCompensated, exec.Status)
	s.Equal([]string{"reserve", "charge", "release"}, steps)
}

// SchedulerSuite tests the job scheduler.
type SchedulerSuite struct {
	suite.Suite
	sched *scheduler.Scheduler
	ctx   context.Context
}

func (s *SchedulerSuite) SetupTest() {
	s.sched = scheduler.New()
	s.ctx = context.Background()
}

func (s *SchedulerSuite) TestScheduleJob() {
	err := s.sched.Schedule("test-job", "@hourly", func(ctx context.Context) error {
		return nil
	})
	s.Require().NoError(err)

	job, err := s.sched.GetJob("test-job")
	s.Require().NoError(err)
	s.Equal("test-job", job.Name)
	s.True(job.Enabled)
}

func (s *SchedulerSuite) TestScheduleOnce() {
	runAt := time.Now().Add(time.Hour)
	err := s.sched.ScheduleOnce("once-job", runAt, func(ctx context.Context) error {
		return nil
	})
	s.Require().NoError(err)

	job, err := s.sched.GetJob("once-job")
	s.Require().NoError(err)
	s.Equal("once", job.Schedule)
}

func (s *SchedulerSuite) TestListJobs() {
	for i := 0; i < 3; i++ {
		s.sched.Schedule("job-"+string(rune('a'+i)), "@daily", func(ctx context.Context) error { return nil })
	}

	jobs := s.sched.ListJobs()
	s.Len(jobs, 3)
}

func (s *SchedulerSuite) TestEnableDisableJob() {
	s.sched.Schedule("toggle-job", "@hourly", func(ctx context.Context) error { return nil })

	err := s.sched.DisableJob("toggle-job")
	s.Require().NoError(err)

	job, _ := s.sched.GetJob("toggle-job")
	s.False(job.Enabled)

	err = s.sched.EnableJob("toggle-job")
	s.Require().NoError(err)

	job, _ = s.sched.GetJob("toggle-job")
	s.True(job.Enabled)
}

func (s *SchedulerSuite) TestRunNow() {
	executed := false
	s.sched.Schedule("run-now-job", "@daily", func(ctx context.Context) error {
		executed = true
		return nil
	})

	exec, err := s.sched.RunNow(s.ctx, "run-now-job")
	s.Require().NoError(err)
	s.True(executed)
	s.Equal(scheduler.JobStatusCompleted, exec.Status)
}

func (s *SchedulerSuite) TestRunNowWithError() {
	s.sched.Schedule("fail-job", "@daily", func(ctx context.Context) error {
		return errors.New("job failed")
	})

	exec, err := s.sched.RunNow(s.ctx, "fail-job")
	s.Error(err)
	s.Equal(scheduler.JobStatusFailed, exec.Status)
}

func TestWorkflowEngineSuite(t *testing.T) {
	suite.Run(t, new(WorkflowEngineSuite))
}

func TestSagaSuite(t *testing.T) {
	suite.Run(t, new(SagaSuite))
}

func TestSchedulerSuite(t *testing.T) {
	suite.Run(t, new(SchedulerSuite))
}
