package ui

import (
	"database/sql"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type CheckoutView struct {
	Tab        *container.TabItem
	active     bool
	HandleScan func(string)
	cartByCode map[string]*CartLine
	cartLines  []*CartLine
	totalLabel *widget.Label
	receipt    *widget.List
}

type CartLine struct {
	ProductID int
	Name      string
	Barcode   string
	UnitPrice float64
	Qty       int
	Stock     int
}

func NewCheckoutTab(db *sql.DB, window fyne.Window, scanner *ScannerService) *CheckoutView {
	_ = db
	_ = window
	_ = scanner

	view := &CheckoutView{
		cartByCode: make(map[string]*CartLine),
	}

	view.receipt = widget.NewList(
		func() int { return len(view.cartLines) },
		func() fyne.CanvasObject {
			name := widget.NewLabel("")
			unit := widget.NewLabel("")
			qty := widget.NewLabel("")
			total := widget.NewLabel("")
			add := widget.NewButton("+", nil)
			subtract := widget.NewButton("−", nil)
			remove := widget.NewButton("Remove", nil)

			info := container.NewGridWithColumns(4, name, unit, qty, total)
			buttons := container.NewHBox(add, subtract, remove)
			return container.NewBorder(nil, nil, nil, buttons, info)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			line := view.cartLines[id]
			border := item.(*fyne.Container)
			buttons := border.Objects[1].(*fyne.Container)
			info := border.Objects[0].(*fyne.Container)

			name := info.Objects[0].(*widget.Label)
			unit := info.Objects[1].(*widget.Label)
			qty := info.Objects[2].(*widget.Label)
			total := info.Objects[3].(*widget.Label)

			add := buttons.Objects[0].(*widget.Button)
			subtract := buttons.Objects[1].(*widget.Button)
			remove := buttons.Objects[2].(*widget.Button)

			name.SetText(line.Name)
			unit.SetText(fmt.Sprintf("Unit: %.2f", line.UnitPrice))
			qty.SetText(fmt.Sprintf("Qty: %d", line.Qty))
			total.SetText(fmt.Sprintf("Line: %.2f", line.UnitPrice*float64(line.Qty)))

			add.OnTapped = func() {
				if line.Qty < line.Stock {
					line.Qty++
					view.refreshReceiptUI()
				}
			}
			subtract.OnTapped = func() {
				if line.Qty > 1 {
					line.Qty--
					view.refreshReceiptUI()
				}
			}
			remove.OnTapped = func() {
				view.removeLine(line.Barcode)
			}
		},
	)

	view.totalLabel = widget.NewLabel("Total: 0.00")

	leftPane := container.NewBorder(nil, view.totalLabel, nil, nil, view.receipt)

	searchPlaceholder := widget.NewLabel("Scan or search products here.")
	resultsPlaceholder := widget.NewLabel("Search results will appear here.")

	rightPane := container.NewVBox(
		searchPlaceholder,
		resultsPlaceholder,
		layout.NewSpacer(),
	)

	content := container.NewHSplit(leftPane, rightPane)
	content.SetOffset(0.55)

	view.Tab = container.NewTabItem("Checkout", content)
	view.HandleScan = func(barcode string) {
		_ = barcode
	}

	return view
}

func (c *CheckoutView) SetActive(active bool) {
	c.active = active
}

func (c *CheckoutView) addOrIncrement(line *CartLine) {
	if existing, ok := c.cartByCode[line.Barcode]; ok {
		if existing.Qty < existing.Stock {
			existing.Qty++
		}
		c.refreshReceiptUI()
		return
	}
	c.cartByCode[line.Barcode] = line
	c.cartLines = append(c.cartLines, line)
	c.refreshReceiptUI()
}

func (c *CheckoutView) removeLine(barcode string) {
	line, ok := c.cartByCode[barcode]
	if !ok {
		return
	}
	delete(c.cartByCode, barcode)
	for i, existing := range c.cartLines {
		if existing == line {
			c.cartLines = append(c.cartLines[:i], c.cartLines[i+1:]...)
			break
		}
	}
	c.refreshReceiptUI()
}

func (c *CheckoutView) recomputeTotal() float64 {
	var total float64
	for _, line := range c.cartLines {
		total += float64(line.Qty) * line.UnitPrice
	}
	return total
}

func (c *CheckoutView) refreshReceiptUI() {
	if c.totalLabel != nil {
		c.totalLabel.SetText(fmt.Sprintf("Total: %.2f", c.recomputeTotal()))
	}
	if c.receipt != nil {
		c.receipt.Refresh()
	}
}
