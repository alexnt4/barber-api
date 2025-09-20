package repository

import (
	"context"

	"github.com/alexnt4/barber-api/internal/domain"
	"gorm.io/gorm"
)

type GormAppoinmentRepo struct {
	db *gorm.DB
}

func NewGormAppoinmentRepo(db *gorm.DB) domain.AppointmentRepo {
	return &GormAppoinmentRepo{db}
}

func (r *GormAppoinmentRepo) Create(ctx context.Context, appt *domain.Appointment) error {
	return r.db.WithContext(ctx).Create(appt).Error
}

func (r *GormAppoinmentRepo) GetById(ctx context.Context, id uint) (*domain.Appointment, error) {
	var appt domain.Appointment

	if err := r.db.WithContext(ctx).Preload("Products").First(&appt, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrorNotFound
		}
		return nil, err
	}

	return &appt, nil
}

func (r *GormAppoinmentRepo) List(ctx context.Context) ([]domain.Appointment, error) {
	var appts []domain.Appointment

	if err := r.db.WithContext(ctx).Preload("Products").Find(&appts).Error; err != nil {
		return nil, err
	}

	return appts, nil
}

func (r *GormAppoinmentRepo) Update(ctx context.Context, appt *domain.Appointment) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(appt).Error
}

func (r *GormAppoinmentRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Appointment{}, id).Error
}
