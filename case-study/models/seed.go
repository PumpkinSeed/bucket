package models

import (
	fmt "fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/rs/xid"
)

// GenerateProfile ...
func GenerateProfile() *Profile {
	return &Profile{
		AboutMe:           gofakeit.HipsterSentence(10),
		Activities:        []string{gofakeit.HackerAdjective()},
		AffiliationCount:  uint64(gofakeit.Number(10, 99)),
		Birthday:          time.Now().Add(-20000 * time.Hour).String(),
		FavoriteBooks:     []string{gofakeit.HackerAdjective()},
		FavoriteMovies:    []string{gofakeit.HackerAdjective()},
		FavoriteMusic:     []string{gofakeit.HackerAdjective()},
		FavoriteQuotes:    []string{gofakeit.HackerAdjective()},
		FavoriteTvShoes:   []string{gofakeit.HackerAdjective()},
		FirstName:         gofakeit.FirstName(),
		Interests:         []string{gofakeit.HackerAdjective()},
		IsApplicationUser: false,
		LastName:          gofakeit.LastName(),
		NotesCount:        uint64(gofakeit.Number(10, 99)),
		PictureBigUrl:     gofakeit.ImageURL(500, 500),
		PictureSmallUrl:   gofakeit.ImageURL(500, 500),
		PictureUrl:        gofakeit.ImageURL(500, 500),
		PoliticalViews:    []string{gofakeit.HackerAdjective()},
		Religion:          gofakeit.HackerAdjective(),
		SchoolCount:       uint64(gofakeit.Number(10, 99)),
		WallCount:         uint64(gofakeit.Number(10, 99)),
		Status:            GenerateStatus(),
		PrimarySchool:     GenerateSchool(),
		HighSchool:        GenerateSchool(),
		Location:          GenerateLocation(),
		Photo:             GeneratePhoto(),
	}
}

// GenerateEvent ...
func GenerateEvent() *Event {
	return &Event{
		Description: gofakeit.HipsterSentence(10),
		EndTime:     time.Now().Add(100 * time.Hour).String(),
		SubType:     gofakeit.HackerAdjective(),
		Type:        gofakeit.HackerAdjective(),
		Host:        gofakeit.HackerAdjective(),
		Photo:       GeneratePhoto(),
		Location:    GenerateLocation(),
		Name:        gofakeit.HackerAdjective(),
		StartTime:   time.Now().Add(98 * time.Hour).String(),
	}
}

// GenerateStatus ...
func GenerateStatus() *Status {
	return &Status{
		Status:     gofakeit.HackerAdjective(),
		UpdateTime: time.Now().Add(-1999 * time.Hour).String(),
	}
}

// GenerateSchool ...
func GenerateSchool() *School {
	return &School{
		Concentrations: []string{gofakeit.HackerAdjective()},
		GraduationYear: 2022,
		Name:           gofakeit.HackerAdjective(),
	}
}

// GenerateLocation ...
func GenerateLocation() *Location {
	addr := gofakeit.Address()
	return &Location{
		City:    addr.City,
		Country: addr.Country,
		State:   addr.State,
		Street:  addr.Street,
		ZipCode: addr.Zip,
	}
}

// GeneratePhoto ...
func GeneratePhoto() *Photo {
	return &Photo{
		Caption:      gofakeit.HackerAdjective(),
		CreatedAt:    time.Now().Add(-19 * time.Hour).String(),
		LargeSource:  gofakeit.ImageURL(500, 500),
		Link:         gofakeit.ImageURL(500, 500),
		MediumSource: gofakeit.ImageURL(500, 500),
		SmallSource:  gofakeit.ImageURL(500, 500),
	}
}

// GenerateOrder ...
func GenerateOrder() *Order {
	name := gofakeit.Name()
	return &Order{
		Token:                 xid.New().String(),
		CreationDate:          time.Now().UTC().String(),
		Status:                "processed",
		PaymentMethod:         "card",
		InvoiceNumber:         fmt.Sprintf("inv-00%d", gofakeit.Number(10000, 99999)),
		Email:                 gofakeit.Email(),
		CardholderName:        name,
		CreditCardLast4Digits: fmt.Sprintf("%d", gofakeit.Number(1000, 9999)),
		BillingAddressName:    name,
		BillingAddress:        GenerateLocation(),
		BillingAddressPhone:   gofakeit.Phone(),
		Notes:                 gofakeit.HipsterSentence(10),
		ShippingAddressName:   name,
		ShippingAddress:       GenerateLocation(),
		ShippingAddressPhone:  gofakeit.Phone(),
		FinalGrandTotal:       443,
		ShippingFees:          0,
		ShippingMethod:        "Free shipping",
		PaymentTransactionId:  xid.New().String(),
		Product:               GenerateProduct(),
		Store:                 GenerateStore(),
	}
}

// GenerateProduct ...
func GenerateProduct() *Product {
	beer := gofakeit.BeerName()
	return &Product{
		Id:          xid.New().String(),
		UserId:      xid.New().String(),
		Name:        beer,
		Description: gofakeit.HipsterSentence(10),
		Slug:        strings.ToLower(beer),
		Price:       1233,
		SalePrice:   1400,
		CurrencyId:  2,
		OnSale:      123,
		Status:      "active",
	}
}

// GenerateStore ...
func GenerateStore() *Store {
	return &Store{
		Id:          xid.New().String(),
		Name:        "productshop",
		Description: "Product shop",
	}
}
