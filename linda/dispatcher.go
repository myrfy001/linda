package linda

import (
	"github.com/amlun/linda/linda/broker"
	_ "github.com/amlun/linda/linda/broker/redis"
	"github.com/amlun/linda/linda/core"
	"github.com/amlun/linda/linda/saver"
	_ "github.com/amlun/linda/linda/saver/cassandra"
)

type dispatcher struct {
	brokerURL string
	broker    broker.Broker
	saverURL  string
	saver     saver.Saver
}

func (d *dispatcher) Init() error {
	b, err := broker.NewBroker(d.brokerURL)
	if err != nil {
		return err
	}
	d.broker = b
	s, err := saver.NewSaver(d.saverURL)
	if err != nil {
		return err
	}
	d.saver = s
	return nil
}

func (d *dispatcher) Close() {
	d.broker.Close()
	d.saver.Close()
}

// push a [period] task to saver
func (d *dispatcher) PushTask(task core.Task) error {
	err := d.saver.PublishTask(&task)
	if err != nil {
		Logger.WithField("action", "PushTask").Errorf("push task error: [%s]", err)
		return err
	}
	Logger.WithField("action", "PushTask").WithField("task", task).Info("ok")
	return nil
}

// push a job to broker and saver
func (d *dispatcher) PushJob(job core.Job) error {
	err := d.broker.PushJob(&job)
	if err != nil {
		Logger.WithField("action", "PushJob").Errorf("push job to broker error: [%s]", err)
		return err
	}
	err = d.saver.PublishJob(&job)
	if err != nil {
		Logger.WithField("action", "PushJob").Errorf("push job to saver error: [%s]", err)
		return err
	}
	Logger.WithField("action", "PushJob").WithField("job", job).Info("ok")
	return nil
}

// get a job and delete it from the queue
func (d *dispatcher) GetJob(queue string) core.Job {
	var job core.Job
	err := d.broker.GetJob(queue, &job)
	if err != nil {
		Logger.WithField("action", "GetJob").Errorf("get job error: [%s]", err)
	} else {
		Logger.WithField("action", "GetJob").WithField("job", job).Info("ok")
	}
	return job

}
