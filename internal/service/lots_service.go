package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"lots-service/internal/domain"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type LotsService struct {
	repo              domain.LotsRepository
	storageServiceURL string
	httpClient        http.Client
}

func NewLotsService(repo domain.LotsRepository, storageServiceURL string) *LotsService {
	return &LotsService{
		repo:              repo,
		storageServiceURL: storageServiceURL,
		httpClient:        http.Client{Timeout: 10 * time.Second},
	}
}

func (s *LotsService) StorageRequest(requestURL string, requestBody io.Reader, contentType string) error {
	req, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("помилка запиту до storage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage повернув помилку: %d, %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (s *LotsService) SaveImages(files []*multipart.FileHeader) ([]string, error) {
	var newReqBody bytes.Buffer
	writer := multipart.NewWriter(&newReqBody)

	var generatedNames []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			slog.Debug("Помилка відкриття файлу", "err", err.Error())
			return nil, err
		}

		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		generatedNames = append(generatedNames, newFilename)

		part, err := writer.CreateFormFile("files", newFilename)
		if err != nil {
			file.Close()
			slog.Debug("Помилка створення форми", "err", err.Error())
			return nil, err
		}

		if _, err := io.Copy(part, file); err != nil {
			file.Close()
			slog.Debug("Помилка копіювання файлу", "err", err.Error())
			return nil, err
		}

		file.Close()
	}

	err := writer.Close()
	if err != nil {
		slog.Debug("Помилка закриття writer", "err", err.Error())
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/api/storage/upload_images", s.storageServiceURL)
	s.StorageRequest(requestURL, &newReqBody, writer.FormDataContentType())

	return generatedNames, nil
}

func (s *LotsService) DeleteImages(filenames []string) error {
	if len(filenames) == 0 {
		return nil
	}

	// JSON: {"filenames": ["a.jpg", "b.jpg"]}
	payload, err := json.Marshal(map[string][]string{
		"filenames": filenames,
	})
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/api/storage/delete_images", s.storageServiceURL)
	s.StorageRequest(requestURL, bytes.NewBuffer(payload), "application/json")

	return nil
}

func (s *LotsService) GetLotsCount() (int, error) {
	return s.repo.GetLotsCount()
}

func (s *LotsService) GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear string) (int, error) {
	return s.repo.GetLotsByParamsCount(brand, model, minPrice, maxPrice, minYear, maxYear)
}

func (s *LotsService) GetLotByID(userID, lotID int) (*domain.Lot, error) {
	return s.repo.GetLotByID(userID, lotID)
}

func (s *LotsService) GetPageLots(userID, page, limit int) (*[]domain.Lot, error) {
	lots, _, err := s.repo.GetLotsByParams(userID, page, limit, "", "", "", "", "", "")
	return lots, err
}

func (s *LotsService) GetLotsByParams(userID int, page, limit int, brand, model, minPrice, maxPrice, minYear, maxYear string) (*[]domain.Lot, int, error) {
	return s.repo.GetLotsByParams(userID, page, limit, brand, model, minPrice, maxPrice, minYear, maxYear)
}

func (s *LotsService) GetBrands() (*[]domain.Brand, error) {
	return s.repo.GetBrands()
}

func (s *LotsService) GetModels(brandName string) (*[]domain.Model, error) {
	return s.repo.GetModels(brandName)
}

func (s *LotsService) GetUserPostedLots(userID int) (*[]domain.Lot, error) {
	return s.repo.GetUserPostedLots(userID)
}

func (s *LotsService) GetUserLikedLots(userID int) (*[]domain.Lot, error) {
	return s.repo.GetUserLikedLots(userID)
}

func (s *LotsService) CreateLot(ctx context.Context, lot *domain.Lot, files []*multipart.FileHeader) error {
	if len(files) > 0 {
		images, err := s.SaveImages(files)
		if err != nil {
			return fmt.Errorf("помилка збереження зображень: %w", err)
		}
		lot.Images = images
	}

	return s.repo.CreateLot(ctx, lot)
}

func (s *LotsService) UpdateLot(ctx context.Context, lot *domain.Lot, newFiles []*multipart.FileHeader, deleteImages []string, oldImages []string) error {
	existingLot, err := s.repo.GetLotByID(lot.SellerID, lot.LotID)
	if err != nil {
		return err
	}

	if existingLot.SellerID != lot.SellerID {
		return fmt.Errorf("sellerID не співпадає з userID")
	}

	if len(deleteImages) > 0 {
		if err := s.DeleteImages(deleteImages); err != nil {
			slog.Warn("Помилка видалення зображень", "err", err.Error())
		}
	}

	var newImageNames []string
	if len(newFiles) > 0 {
		newImageNames, err = s.SaveImages(newFiles)
		if err != nil {
			return err
		}
	}

	deletedSet := make(map[string]bool)
	for _, img := range deleteImages {
		deletedSet[img] = true
	}

	var finalImages []string
	for _, img := range oldImages {
		if !deletedSet[img] {
			finalImages = append(finalImages, img)
		}
	}
	finalImages = append(finalImages, newImageNames...)

	lot.Images = finalImages

	return s.repo.UpdateLot(ctx, lot)
}

func (s *LotsService) DeleteLot(ctx context.Context, lotID, userID int) error {
	lot, err := s.repo.GetLotByID(userID, lotID)
	if err != nil {
		return err
	}
	if lot.SellerID != userID {
		return fmt.Errorf("sellerID не співпадає з userID")
	}

	if err := s.repo.DeleteLot(ctx, lotID); err != nil {
		return err
	}

	if err := s.DeleteImages(lot.Images); err != nil {
		slog.Warn("Помилка очистки зображень", "lotID", lotID, "err", err.Error())
		return err
	}

	return nil
}

func (s *LotsService) LikeLot(userID, lotID int) error {
	return s.repo.LikeLot(userID, lotID)
}

func (s *LotsService) UnlikeLot(userID, lotID int) error {
	return s.repo.UnlikeLot(userID, lotID)
}

func (s *LotsService) BuyLot(userID, lotID int) error {
	lot, err := s.repo.GetLotByID(userID, lotID)
	if err != nil {
		return err
	}
	if lot.SaleStatus == "Продано" {
		return fmt.Errorf("лот вже продано")
	}

	return s.repo.MarkLotAsSold(lotID)
}
