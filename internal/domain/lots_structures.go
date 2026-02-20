package domain

import "context"

type Brand struct {
	BrandID   int
	BrandName string
}

type Model struct {
	ModelID   int
	BrandID   int
	ModelName string
}

type Car struct {
	CarID        int
	BrandID      int
	ModelID      int
	Brand        string
	Model        string
	Engine       string
	Transmission string
	WheelDrive   string
	MadeYear     int
	VinCode      string
	Color        string
	Mileage      int
}

type Lot struct {
	LotID       int
	SellerID    int
	Car         Car
	PostDate    string
	SalePrice   int
	SaleStatus  string
	Description string
	IsLiked     bool
	Images      []string
}

type LotsRepository interface {
	GetLotsCount() (int, error)
	GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear string) (int, error)
	GetLotByID(userID, lotID int) (*Lot, error)
	GetLotsByParams(userID int, page, limit int, brand, model, minPrice, maxPrice, minYear, maxYear string) (*[]Lot, int, error)

	GetBrands() (*[]Brand, error)
	GetModels(brandName string) (*[]Model, error)

	GetUserPostedLots(userID int) (*[]Lot, error)
	GetUserLikedLots(userID int) (*[]Lot, error)

	CreateLot(ctx context.Context, lot *Lot) error
	UpdateLot(ctx context.Context, lot *Lot) error
	DeleteLot(ctx context.Context, lotID int) error

	LikeLot(userID, lotID int) error
	UnlikeLot(userID, lotID int) error

	MarkLotAsSold(lotID int) error
}
