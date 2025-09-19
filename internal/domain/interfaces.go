package domain

import "context"

type AppointmentRepo interface {
	Create(ctx context.Context, appt *Appointment) error
	GetById(ctx context.Context, id uint) (*Appointment, error)
	List(ctx context.Context) ([]Appointment, error)
	Update(ctx context.Context, appt *Appointment) error
	Delete(ctx context.Context, id uint) error
}
