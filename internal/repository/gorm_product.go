package repository

import (
	"context"

	"github.com/alexnt4/barber-api/internal/domain"
	"gorm.io/gorm"
)

type GormProductRepo struct {
	db *gorm.DB
}

func NewGormProducttRepo(db *gorm.DB) domain.ProductRepo {
	return &GormProductRepo{db}
}

func (r *GormProductRepo) Create(ctx context.Context, prod *domain.Product) error {
	return r.db.WithContext(ctx).Create(prod).Error
}

func (r *GormProductRepo) GetById(ctx context.Context, id uint) (*domain.Product, error) {
	var prod domain.Product

	if err := r.db.WithContext(ctx).First(&prod, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrorNotFound
		}
		return nil, err
	}

	return &prod, nil
}

func (r *GormProductRepo) List(ctx context.Context) ([]domain.Product, error) {
	var prod []domain.Product

	// if err := r.db.WithContext(ctx).Preload("Products").Find(&prod).Error; err != nil {
	//	return nil, err
	//}
	if err := r.db.WithContext(ctx).Find(&prod).Error; err != nil {
		return nil, err
	}

	return prod, nil
}

func (r *GormProductRepo) Update(ctx context.Context, prod *domain.Product) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(prod).Error
}

func (r *GormProductRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}
