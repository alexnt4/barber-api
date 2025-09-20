package domain

import "context"

type AppointmentRepo interface {
	Create(ctx context.Context, appt *Appointment) error
	GetById(ctx context.Context, id uint) (*Appointment, error)
	List(ctx context.Context) ([]Appointment, error)
	Update(ctx context.Context, appt *Appointment) error
	Delete(ctx context.Context, id uint) error
}

type ProductRepo interface {
	Create(ctx context.Context, prod *Product) error
	GetById(ctx context.Context, id uint) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, prod *Product) error
	Delete(ctx context.Context, id uint) error
}
