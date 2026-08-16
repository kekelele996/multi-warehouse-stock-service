package service

import (
	"log/slog"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/util"
)

// ProductService 商品业务逻辑。
type ProductService struct {
	repo         *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	logger       *slog.Logger
}

// NewProductService 构造商品服务。
func NewProductService(repo *repository.ProductRepository, categoryRepo *repository.CategoryRepository, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, categoryRepo: categoryRepo, logger: logger}
}

// Create 创建商品。
func (s *ProductService) Create(name, sku string, categoryID uint64, spec, unit string, weight, volume float64,
	barcode string, minStock, maxStock, shelfLifeDays int, imageURL string) (*model.Product, error) {
	if categoryID > 0 {
		if _, err := s.categoryRepo.FindByID(categoryID); err != nil {
			return nil, util.Wrap(err, "Product[category_id=%d] create: category not found", categoryID)
		}
	}
	if unit == "" {
		unit = "件"
	}
	p := &model.Product{Name: name, SKU: sku, CategoryID: categoryID, Spec: spec, Unit: unit,
		Weight: weight, Volume: volume, Barcode: barcode, MinStock: minStock, MaxStock: maxStock,
		ShelfLifeDays: shelfLifeDays, ImageURL: imageURL}
	if err := s.repo.Create(p); err != nil {
		return nil, util.Wrap(err, "Product[sku=%s] create failed", sku)
	}
	s.logger.Info(constants.LogProductCreateSuccess, "product_id", p.ID)
	return p, nil
}

// Update 更新商品。
func (s *ProductService) Update(id uint64, name, spec, unit string, weight, volume float64, minStock, maxStock int, imageURL string) (*model.Product, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Product[id=%d] update find failed", id)
	}
	if name != "" {
		p.Name = name
	}
	if spec != "" {
		p.Spec = spec
	}
	if unit != "" {
		p.Unit = unit
	}
	if weight > 0 {
		p.Weight = weight
	}
	if volume > 0 {
		p.Volume = volume
	}
	if minStock >= 0 {
		p.MinStock = minStock
	}
	if maxStock >= 0 {
		p.MaxStock = maxStock
	}
	if imageURL != "" {
		p.ImageURL = imageURL
	}
	if err := s.repo.Update(p); err != nil {
		return nil, util.Wrap(err, "Product[id=%d] update save failed", id)
	}
	s.logger.Info(constants.LogProductUpdateSuccess, "product_id", p.ID)
	return p, nil
}

// List 分页查询商品。
func (s *ProductService) List(page, pageSize int, categoryID uint64, keyword string) ([]model.Product, int64, error) {
	return s.repo.List(page, pageSize, categoryID, keyword)
}

// Get 商品详情。
func (s *ProductService) Get(id uint64) (*model.Product, error) {
	return s.repo.FindByID(id)
}

// CreateCategory 创建分类。
func (s *ProductService) CreateCategory(name string, parentID uint64, sortOrder int, icon string) (*model.Category, error) {
	c := &model.Category{Name: name, ParentID: parentID, SortOrder: sortOrder, Icon: icon}
	if err := s.categoryRepo.Create(c); err != nil {
		return nil, util.Wrap(err, "Category[name=%s] create failed", name)
	}
	s.logger.Info(constants.LogCategoryCreateSuccess, "category_id", c.ID)
	return c, nil
}

// ListCategories 分类列表。
func (s *ProductService) ListCategories() ([]model.Category, error) {
	return s.categoryRepo.List()
}
