package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	queue []Task
}

func NewScheduler() Scheduler {
	return Scheduler{
		queue: make([]Task, 0, 10),
	}
}

func (s *Scheduler) max(i, j int) int {
	if s.queue[i].Priority > s.queue[j].Priority {
		return i
	}

	return j
}

func (s *Scheduler) siftDown(i int) {
	l, r := i*2+1, i*2+2
	if r < len(s.queue) {
		if max := s.max(l, r); s.max(max, i) != i {
			s.queue[i], s.queue[max] = s.queue[max], s.queue[i]
			s.siftDown(max)
		}
	} else if l < len(s.queue) && s.queue[l].Priority > s.queue[i].Priority {
		s.queue[i], s.queue[l] = s.queue[l], s.queue[i]
	}
}

func (s *Scheduler) siftUp(i int) {
	if s.queue[i].Priority > s.queue[(i-1)/2].Priority {
		s.queue[i], s.queue[(i-1)/2] = s.queue[(i-1)/2], s.queue[i]
		s.siftUp((i - 1) / 2)
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.queue = append(s.queue, task)
	s.siftUp(len(s.queue) - 1)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	for i, task := range s.queue {
		if task.Identifier != taskID {
			continue
		}

		s.queue[i].Priority = newPriority
		if s.queue[i].Priority > s.queue[(i-1)/2].Priority {
			s.siftUp(i)
		} else {
			s.siftDown(i)
		}

		break
	}
}

func (s *Scheduler) GetTask() Task {
	if len(s.queue) == 0 {
		return Task{}
	}

	task := s.queue[0]
	s.queue[0], s.queue[len(s.queue)-1] = s.queue[len(s.queue)-1], s.queue[0]
	s.queue = s.queue[:len(s.queue)-1]

	s.siftDown(0)

	return task
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	task1.Priority = 100
	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
