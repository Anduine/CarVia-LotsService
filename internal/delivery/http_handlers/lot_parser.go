package http_handlers

import (
	"fmt"
	"lots-service/internal/domain"
	"net/http"
	"strconv"
)

func ParseLotFromRequest(r *http.Request) (domain.Lot, error) {
	parseInt := func(key string) (int, error) {
		val := r.FormValue(key)
		if val == "" || val == "null" || val == "undefined" {
			return 0, nil
		}
		return strconv.Atoi(val)
	}

	var err error
	var lot domain.Lot

	lot.Car.Brand = r.FormValue("Brand")
	lot.Car.Model = r.FormValue("Model")
	lot.Car.Engine = r.FormValue("Engine")
	lot.Car.Transmission = r.FormValue("Transmission")
	lot.Car.WheelDrive = r.FormValue("WheelDrive")
	lot.Car.Color = r.FormValue("Color")
	lot.Car.VinCode = r.FormValue("VinCode")
	lot.Description = r.FormValue("Description")
	lot.SaleStatus = r.FormValue("SaleStatus")

	if lot.Car.MadeYear, err = parseInt("MadeYear"); err != nil {
		return lot, fmt.Errorf("bad MadeYear: %w", err)
	}
	if lot.Car.Mileage, err = parseInt("Mileage"); err != nil {
		return lot, fmt.Errorf("bad Mileage: %w", err)
	}
	if lot.SalePrice, err = parseInt("SalePrice"); err != nil {
		return lot, fmt.Errorf("bad SalePrice: %w", err)
	}
	// if lot.Car.CarID, err = parseInt("CarID"); err != nil {
	// 	return lot, fmt.Errorf("bad CarID: %w", err)
	// }
	// if lot.Car.BrandID, err = parseInt("BrandID"); err != nil {
	// 	return lot, fmt.Errorf("bad BrandID: %w", err)
	// }
	// if lot.Car.ModelID, err = parseInt("ModelID"); err != nil {
	// 	return lot, fmt.Errorf("bad ModelID: %w", err)
	// }

	return lot, nil
}
