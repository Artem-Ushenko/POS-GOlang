package ui

import (
	"database/sql"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"pos-system/internal/store"
)

type productForm struct {
	name    *widget.Entry
	barcode *widget.Entry
	price   *widget.Entry
	stock   *widget.Entry
}

func ProductsTab(db *sql.DB) fyne.CanvasObject {
	var products []store.Product
	selectedIndex := -1

	form := productForm{
		name:    widget.NewEntry(),
		barcode: widget.NewEntry(),
		price:   widget.NewEntry(),
		stock:   widget.NewEntry(),
	}

	refresh := func(list *widget.List) {
		items, err := store.ListProducts(db)
		if err != nil {
			fmt.Println("Failed to load products:", err)
			return
		}
		products = items
		selectedIndex = -1
		form.name.SetText("")
		form.barcode.SetText("")
		form.price.SetText("")
		form.stock.SetText("")
		list.Refresh()
	}

	list := widget.NewList(
		func() int { return len(products) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			product := products[id]
			item.(*widget.Label).SetText(fmt.Sprintf("%d - %s", product.ID, product.Name))
		},
	)

	addButton := widget.NewButton("Add", func() {
		barcode := form.barcode.Text
		if barcode != "" {
			existing, err := store.GetProductByBarcode(db, barcode)
			if err == nil && existing.ID != 0 {
				fmt.Println("Barcode already exists.")
				return
			}
			if err != nil && err != sql.ErrNoRows {
				fmt.Println("Failed to validate barcode:", err)
				return
			}
		}
		price, err := strconv.ParseFloat(form.price.Text, 64)
		if err != nil {
			return
		}
		stock, err := strconv.ParseInt(form.stock.Text, 10, 64)
		if err != nil {
			return
		}
		_, err = store.CreateProduct(db, store.Product{
			Name:    form.name.Text,
			Barcode: barcode,
			Price:   price,
			Stock:   stock,
		})
		if err != nil {
			fmt.Println("Failed to create product:", err)
			return
		}
		refresh(list)
	})

	updateButton := widget.NewButton("Update", func() {
		if selectedIndex < 0 || selectedIndex >= len(products) {
			return
		}
		barcode := form.barcode.Text
		if barcode != "" {
			existing, err := store.GetProductByBarcode(db, barcode)
			if err == nil && existing.ID != 0 && existing.ID != products[selectedIndex].ID {
				fmt.Println("Barcode already exists.")
				return
			}
			if err != nil && err != sql.ErrNoRows {
				fmt.Println("Failed to validate barcode:", err)
				return
			}
		}
		price, err := strconv.ParseFloat(form.price.Text, 64)
		if err != nil {
			return
		}
		stock, err := strconv.ParseInt(form.stock.Text, 10, 64)
		if err != nil {
			return
		}
		product := products[selectedIndex]
		product.Name = form.name.Text
		product.Barcode = barcode
		product.Price = price
		product.Stock = stock
		if err := store.UpdateProduct(db, product); err != nil {
			fmt.Println("Failed to update product:", err)
			return
		}
		refresh(list)
	})

	deleteButton := widget.NewButton("Delete", func() {
		if selectedIndex < 0 || selectedIndex >= len(products) {
			return
		}
		if err := store.DeleteProduct(db, products[selectedIndex].ID); err != nil {
			fmt.Println("Failed to delete product:", err)
			return
		}
		refresh(list)
	})

	formItems := []*widget.FormItem{
		{Text: "Name", Widget: form.name},
		{Text: "Barcode", Widget: form.barcode},
		{Text: "Price", Widget: form.price},
		{Text: "Stock", Widget: form.stock},
	}
	formWidget := widget.NewForm(formItems...)

	buttons := container.NewHBox(addButton, updateButton, deleteButton)
	controls := container.NewVBox(formWidget, buttons)

	footer := widget.NewLabel("Selected: none")
	list.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		product := products[id]
		form.name.SetText(product.Name)
		form.barcode.SetText(product.Barcode)
		form.price.SetText(fmt.Sprintf("%.2f", product.Price))
		form.stock.SetText(strconv.FormatInt(product.Stock, 10))
		footer.SetText("Selected ID: " + strconv.FormatInt(product.ID, 10))
	}

	refresh(list)

	return container.NewBorder(nil, footer, nil, nil,
		container.NewHSplit(
			container.NewBorder(nil, nil, nil, nil, list),
			container.NewVBox(controls, layout.NewSpacer()),
		),
	)
}
