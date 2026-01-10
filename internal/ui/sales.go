package ui

import (
	"database/sql"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"pos-system/internal/store"
)

type SalesView struct {
	Tab        *container.TabItem
	HandleScan func(string)
}

type cartItem struct {
	Product  store.Product
	Quantity int64
}

func NewSalesTab(db *sql.DB, window fyne.Window) *SalesView {
	var (
		items      []cartItem
		itemByCode = make(map[string]int)
	)

	list := widget.NewList(
		func() int { return len(items) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			cart := items[id]
			lineTotal := float64(cart.Quantity) * cart.Product.Price
			item.(*widget.Label).SetText(
				fmt.Sprintf("%s x%d @ %.2f = %.2f", cart.Product.Name, cart.Quantity, cart.Product.Price, lineTotal),
			)
		},
	)

	status := widget.NewLabel("Scan a barcode to add items.")

	handleScan := func(barcode string) {
		if barcode == "" {
			return
		}
		product, err := store.GetProductByBarcode(db, barcode)
		if err == sql.ErrNoRows {
			dialog.NewInformation("Unknown Barcode", "No product found for barcode: "+barcode, window).Show()
			return
		}
		if err != nil {
			dialog.NewError(err, window).Show()
			return
		}

		if index, ok := itemByCode[barcode]; ok {
			items[index].Quantity++
		} else {
			items = append(items, cartItem{Product: product, Quantity: 1})
			itemByCode[barcode] = len(items) - 1
		}
		status.SetText(fmt.Sprintf("Added %s", product.Name))
		list.Refresh()
	}

	content := container.NewBorder(nil, status, nil, nil, list)

	return &SalesView{
		Tab:        container.NewTabItem("Sales", content),
		HandleScan: handleScan,
	}
}
