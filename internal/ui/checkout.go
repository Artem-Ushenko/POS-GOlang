package ui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CheckoutView struct {
	Tab        *container.TabItem
	active     bool
	HandleScan func(string)
}

func NewCheckoutTab(db *sql.DB, window fyne.Window, scanner *ScannerService) *CheckoutView {
	receiptPlaceholder := widget.NewLabel("Receipt items will appear here.")
	totalPlaceholder := widget.NewLabel("Total: 0.00")
	leftPane := container.NewBorder(nil, totalPlaceholder, nil, nil, receiptPlaceholder)

	searchPlaceholder := widget.NewLabel("Scan or search products here.")
	resultsPlaceholder := widget.NewLabel("Search results will appear here.")
	rightPane := container.NewVBox(searchPlaceholder, resultsPlaceholder, layout.NewSpacer())

	content := container.NewHSplit(leftPane, rightPane)
	content.SetOffset(0.55)

	view := &CheckoutView{}
	view.Tab = container.NewTabItem("Checkout", content)
	view.HandleScan = func(barcode string) {
		_ = barcode
	}

	return view
}

func (c *CheckoutView) SetActive(active bool) {
	c.active = active
}
