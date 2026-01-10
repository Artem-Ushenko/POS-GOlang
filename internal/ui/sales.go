package ui

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"pos-system/internal/store"
)

type SalesView struct {
	Tab           *container.TabItem
	HandleScan    func(string)
	ClearCart     func()
	OnSaleCreated func()
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
	customerIDEntry := widget.NewEntry()
	customerIDEntry.SetPlaceHolder("Customer ID")

	totalLabel := widget.NewLabel("Total: 0.00")

	clearCart := func() {
		items = nil
		itemByCode = make(map[string]int)
		status.SetText("Scan a barcode to add items.")
		totalLabel.SetText("Total: 0.00")
		list.Refresh()
	}

	updateTotal := func() {
		var total float64
		for _, item := range items {
			total += float64(item.Quantity) * item.Product.Price
		}
		totalLabel.SetText(fmt.Sprintf("Total: %.2f", total))
	}

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
		updateTotal()
	}

	view := &SalesView{}

	createSaleButton := widget.NewButton("Create Sale", func() {
		if len(items) == 0 {
			dialog.NewInformation("Empty Cart", "Scan products before creating a sale.", window).Show()
			return
		}
		if customerIDEntry.Text == "" {
			dialog.NewInformation("Customer Required", "Enter a customer ID to create a sale.", window).Show()
			return
		}
		customerID, err := strconv.ParseInt(customerIDEntry.Text, 10, 64)
		if err != nil {
			dialog.NewInformation("Invalid Customer", "Customer ID must be a number.", window).Show()
			return
		}

		saleItems := make([]store.SaleItem, 0, len(items))
		for _, item := range items {
			saleItems = append(saleItems, store.SaleItem{
				ProductID: item.Product.ID,
				Quantity:  item.Quantity,
			})
		}

		if _, err := store.CreateSale(db, &customerID, saleItems); err != nil {
			if errors.Is(err, store.ErrInsufficientStock) {
				dialog.NewInformation("Insufficient Stock", err.Error(), window).Show()
				return
			}
			dialog.NewError(err, window).Show()
			return
		}

		clearCart()
		if view.OnSaleCreated != nil {
			view.OnSaleCreated()
		}
	})

	controls := container.NewHBox(customerIDEntry, createSaleButton, totalLabel)

	content := container.NewBorder(controls, status, nil, nil, list)

	view.Tab = container.NewTabItem("Sales", content)
	view.HandleScan = handleScan
	view.ClearCart = clearCart
	return view
}
