package ui

import (
	"database/sql"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"pos-system/internal/store"
)

type CheckoutView struct {
	Tab        *container.TabItem
	active     bool
	HandleScan func(string)
	cartByCode map[string]*CartLine
	cartLines  []*CartLine
	totalLabel *widget.Label
	receipt    *widget.List
	status     *widget.Label
	results    []store.Product
	resultsUI  *widget.List
	search     *widget.Entry
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
				view.withUI(func() {
					view.applyQuantity(line, line.Qty+1)
				})
			}
			subtract.OnTapped = func() {
				view.withUI(func() {
					view.applyQuantity(line, line.Qty-1)
				})
			}
			remove.OnTapped = func() {
				view.withUI(func() {
					view.removeLine(line.Barcode)
					view.setStatus("Removed " + line.Name)
				})
			}
		},
	)

	view.totalLabel = widget.NewLabel("Total: 0.00")

	leftPane := container.NewBorder(nil, view.totalLabel, nil, nil, view.receipt)
	view.status = widget.NewLabel("Ready to scan.")

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by name or barcode")
	view.search = searchEntry

	searchMode := widget.NewCheck("Search mode (pause scanner focus)", func(checked bool) {
		scanner.SetFocusLockEnabled(!checked)
		view.withUI(func() {
			if checked {
				window.Canvas().Focus(searchEntry)
			} else {
				window.Canvas().Focus(scanner.Widget())
			}
		})
	})

	view.resultsUI = widget.NewList(
		func() int { return len(view.results) },
		func() fyne.CanvasObject {
			name := widget.NewLabel("")
			meta := widget.NewLabel("")
			add := widget.NewButton("Add", nil)
			text := container.NewVBox(name, meta)
			return container.NewBorder(nil, nil, nil, add, text)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			product := view.results[id]
			border := item.(*fyne.Container)
			add := border.Objects[1].(*widget.Button)
			text := border.Objects[0].(*fyne.Container)
			name := text.Objects[0].(*widget.Label)
			meta := text.Objects[1].(*widget.Label)

			name.SetText(product.Name)
			meta.SetText(fmt.Sprintf("Barcode: %s | Price: %.2f | Stock: %d", product.Barcode, product.Price, product.Stock))
			add.OnTapped = func() {
				view.withUI(func() {
					line := &CartLine{
						ProductID: int(product.ID),
						Name:      product.Name,
						Barcode:   product.Barcode,
						UnitPrice: product.Price,
						Qty:       1,
						Stock:     int(product.Stock),
					}
					view.addOrIncrement(line)
					searchEntry.SetText("")
					focusInput()
				})
			}
		},
	)

	focusInput := func() {
		window.Canvas().Focus(searchEntry)
	}

	searchExact := func(barcode string) (store.Product, bool, error) {
		var product store.Product
		err := db.QueryRow(
			`SELECT id, name, barcode, price, stock FROM products WHERE barcode = ? LIMIT 1`,
			barcode,
		).Scan(&product.ID, &product.Name, &product.Barcode, &product.Price, &product.Stock)
		if err == sql.ErrNoRows {
			return store.Product{}, false, nil
		}
		if err != nil {
			return store.Product{}, false, err
		}
		return product, true, nil
	}

	searchResults := func(query string) ([]store.Product, error) {
		likeQuery := "%" + query + "%"
		rows, err := db.Query(
			`SELECT id, name, barcode, price, stock FROM products WHERE (name LIKE ? OR barcode LIKE ?) AND stock > 0 ORDER BY name LIMIT ?`,
			likeQuery,
			likeQuery,
			50,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var results []store.Product
		for rows.Next() {
			var product store.Product
			if err := rows.Scan(&product.ID, &product.Name, &product.Barcode, &product.Price, &product.Stock); err != nil {
				return nil, err
			}
			results = append(results, product)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return results, nil
	}

	handleInput := func(value string) {
		query := strings.TrimSpace(value)
		if query == "" {
			return
		}

		go func() {
			product, found, err := searchExact(query)
			if err != nil {
				view.withUI(func() {
					view.setStatus(fmt.Sprintf("Lookup failed: %v", err))
				})
				return
			}
			if found {
				view.withUI(func() {
					if product.Stock <= 0 {
						view.setStatus("Out of stock: " + product.Name)
						return
					}
					line := &CartLine{
						ProductID: int(product.ID),
						Name:      product.Name,
						Barcode:   product.Barcode,
						UnitPrice: product.Price,
						Qty:       1,
						Stock:     int(product.Stock),
					}
					view.addOrIncrement(line)
					searchEntry.SetText("")
					focusInput()
				})
				return
			}

			results, err := searchResults(query)
			view.withUI(func() {
				if err != nil {
					view.setStatus(fmt.Sprintf("Search failed: %v", err))
					return
				}
				view.results = results
				view.resultsUI.Refresh()
				if len(results) == 0 {
					view.setStatus("No matches for: " + query)
				}
				focusInput()
			})
		}()
	}

	searchEntry.OnSubmitted = func(value string) {
		handleInput(value)
	}
	searchEntry.OnChanged = func(value string) {
		if strings.TrimSpace(value) == "" {
			view.results = nil
			view.resultsUI.Refresh()
			view.setStatus("Ready to scan.")
		}
	}

	rightPane := container.NewVBox(searchEntry, searchMode, view.resultsUI, view.status, layout.NewSpacer())

	content := container.NewHSplit(leftPane, rightPane)
	content.SetOffset(0.55)

	view.Tab = container.NewTabItem("Checkout", content)
	view.HandleScan = func(barcode string) {
		if !view.active {
			return
		}
		handleInput(barcode)
	}

	return view
}

func (c *CheckoutView) SetActive(active bool) {
	c.active = active
}

func (c *CheckoutView) addOrIncrement(line *CartLine) {
	if existing, ok := c.cartByCode[line.Barcode]; ok {
		c.applyQuantity(existing, existing.Qty+1)
		return
	}
	if line.Stock == 0 {
		c.setStatus("Out of stock: " + line.Name)
		return
	}
	c.cartByCode[line.Barcode] = line
	c.cartLines = append(c.cartLines, line)
	c.setStatus("Added " + line.Name)
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

func (c *CheckoutView) setStatus(message string) {
	if c.status != nil {
		c.status.SetText(message)
	}
}

func (c *CheckoutView) applyQuantity(line *CartLine, desired int) {
	if line.Stock == 0 {
		c.setStatus("Out of stock: " + line.Name)
		return
	}
	target := desired
	if target < 1 {
		target = 1
		if line.Qty == 1 {
			c.setStatus("Min qty is 1: " + line.Name)
			return
		}
	}
	if target > line.Stock {
		target = line.Stock
		if line.Qty == line.Stock {
			c.setStatus("Max stock reached: " + line.Name)
			return
		}
	}
	if target == line.Qty {
		return
	}
	if target > line.Qty {
		c.setStatus("Added " + line.Name)
	}
	line.Qty = target
	c.refreshReceiptUI()
}

func (c *CheckoutView) withUI(action func()) {
	app := fyne.CurrentApp()
	if app == nil {
		action()
		return
	}
	app.Driver().RunOnMain(action)
}

func (c *CheckoutView) FocusSearch() {
	if c.search == nil {
		return
	}
	c.withUI(func() {
		c.search.Focus()
	})
}
