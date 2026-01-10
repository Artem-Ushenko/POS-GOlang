package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"pos-system/internal/store"
	"pos-system/internal/ui"
)

func main() {
	db, err := store.OpenDB("pos.db")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := store.Migrate(db); err != nil {
		fmt.Println("Failed to migrate database:", err)
		os.Exit(1)
	}

	gui := app.New()
	window := gui.NewWindow("POS System")

	scanner := ui.NewScannerService()
	salesView := ui.NewSalesTab(db, window)
	salesActive := false
	scanner.OnScan(func(barcode string) {
		if !salesActive {
			return
		}
		salesView.HandleScan(barcode)
	})
	scanner.Start(window)
	defer scanner.Stop()

	tabs := container.NewAppTabs(
		container.NewTabItem("Customers", ui.CustomersTab(db)),
		container.NewTabItem("Products", ui.ProductsTab(db)),
		salesView.Tab,
	)
	tabs.OnSelected = func(item *container.TabItem) {
		salesActive = item == salesView.Tab
	}
	salesActive = tabs.Selected() == salesView.Tab

	window.SetContent(container.NewMax(tabs, scanner.Widget()))
	window.Resize(fyne.NewSize(900, 600))
	window.ShowAndRun()
}
