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

type customerForm struct {
	name  *widget.Entry
	email *widget.Entry
	phone *widget.Entry
}

func CustomersTab(db *sql.DB) fyne.CanvasObject {
	var customers []store.Customer
	selectedIndex := -1

	form := customerForm{
		name:  widget.NewEntry(),
		email: widget.NewEntry(),
		phone: widget.NewEntry(),
	}

	refresh := func(list *widget.List) {
		items, err := store.ListCustomers(db)
		if err != nil {
			fmt.Println("Failed to load customers:", err)
			return
		}
		customers = items
		selectedIndex = -1
		form.name.SetText("")
		form.email.SetText("")
		form.phone.SetText("")
		list.Refresh()
	}

	list := widget.NewList(
		func() int { return len(customers) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, item fyne.CanvasObject) {
			customer := customers[id]
			item.(*widget.Label).SetText(fmt.Sprintf("%d - %s", customer.ID, customer.Name))
		},
	)

	addButton := widget.NewButton("Add", func() {
		if form.name.Text == "" || form.email.Text == "" || form.phone.Text == "" {
			return
		}
		_, err := store.CreateCustomer(db, store.Customer{
			Name:  form.name.Text,
			Email: form.email.Text,
			Phone: form.phone.Text,
		})
		if err != nil {
			fmt.Println("Failed to create customer:", err)
			return
		}
		refresh(list)
	})

	updateButton := widget.NewButton("Update", func() {
		if selectedIndex < 0 || selectedIndex >= len(customers) {
			return
		}
		customer := customers[selectedIndex]
		customer.Name = form.name.Text
		customer.Email = form.email.Text
		customer.Phone = form.phone.Text
		if err := store.UpdateCustomer(db, customer); err != nil {
			fmt.Println("Failed to update customer:", err)
			return
		}
		refresh(list)
	})

	deleteButton := widget.NewButton("Delete", func() {
		if selectedIndex < 0 || selectedIndex >= len(customers) {
			return
		}
		if err := store.DeleteCustomer(db, customers[selectedIndex].ID); err != nil {
			fmt.Println("Failed to delete customer:", err)
			return
		}
		refresh(list)
	})

	formItems := []*widget.FormItem{
		{Text: "Name", Widget: form.name},
		{Text: "Email", Widget: form.email},
		{Text: "Phone", Widget: form.phone},
	}
	formWidget := widget.NewForm(formItems...)

	buttons := container.NewHBox(addButton, updateButton, deleteButton)
	controls := container.NewVBox(formWidget, buttons)

	footer := widget.NewLabel("Selected: none")
	list.OnSelected = func(id widget.ListItemID) {
		selectedIndex = id
		customer := customers[id]
		form.name.SetText(customer.Name)
		form.email.SetText(customer.Email)
		form.phone.SetText(customer.Phone)
		footer.SetText("Selected ID: " + strconv.FormatInt(customer.ID, 10))
	}

	refresh(list)

	return container.NewBorder(nil, footer, nil, nil,
		container.NewHSplit(
			container.NewBorder(nil, nil, nil, nil, list),
			container.NewVBox(controls, layout.NewSpacer()),
		),
	)
}
