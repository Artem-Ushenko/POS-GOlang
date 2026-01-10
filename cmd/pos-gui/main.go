package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

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

	tabs := container.NewAppTabs(
		container.NewTabItem("Customers", ui.CustomersTab(db)),
		container.NewTabItem("Products", ui.ProductsTab(db)),
		container.NewTabItem("Sales", widget.NewLabel("Sales view coming soon.")),
	)

	window.SetContent(tabs)
	window.Resize(fyne.NewSize(900, 600))
	window.ShowAndRun()
}
