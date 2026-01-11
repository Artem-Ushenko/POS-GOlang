package main

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"pos-system/internal/backup"
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
	window.SetOnClosed(func() {
		go func() {
			if _, err := backup.BackupDatabase("pos.db", "backups"); err != nil {
				fmt.Println("Auto-backup failed:", err)
			}
		}()
	})

	scanner := ui.NewScannerService()
	checkoutView := ui.NewCheckoutTab(db, window, scanner)
	salesView := ui.NewSalesTab(db, window)
	productsView := ui.NewProductsTab(db)
	checkoutActive := false
	scanner.OnScan(func(barcode string) {
		if checkoutActive {
			checkoutView.HandleScan(barcode)
		}
	})
	scanner.Start(window)
	defer scanner.Stop()

	tabs := container.NewAppTabs(
		checkoutView.Tab,
		container.NewTabItem("Customers", ui.CustomersTab(db)),
		productsView.Tab,
		salesView.Tab,
	)
	tabs.OnSelected = func(item *container.TabItem) {
		checkoutActive = item == checkoutView.Tab
		checkoutView.SetActive(checkoutActive)
	}
	checkoutActive = tabs.Selected() == checkoutView.Tab
	checkoutView.SetActive(checkoutActive)

	salesView.OnSaleCreated = func() {
		productsView.Refresh()
	}

	backupButton := widget.NewButton("Backup Now", func() {
		go func() {
			path, err := backup.BackupDatabase("pos.db", "backups")
			if err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			dialog.NewInformation("Backup Complete", "Backup saved to "+path, window).Show()
		}()
	})

	content := container.NewBorder(backupButton, nil, nil, nil, tabs)
	window.SetContent(container.NewMax(content, scanner.Widget()))
	window.Resize(fyne.NewSize(900, 600))
	window.ShowAndRun()
}
