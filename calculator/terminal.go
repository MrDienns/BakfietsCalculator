package calculator

import (
	"fmt"
	"sort"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// view is a struct used to keep track of sub containers for the GUI and some of the data that needs to high level be
// available. It also has a reference to the catalogue data so that all views can visualise this data.
type view struct {
	data          *Data
	app           *tview.Application
	brandsView    *tview.List
	brands        []*Brand
	modelsView    *tview.List
	models        []*Model
	optionsView   *tview.List
	options       []*Option
	schedulerView *tview.Form
	priceView     *tview.Table

	startdate, enddate string
}

// NewView constructs a new *view and returns it. It initialises all GUI's and adds the provided *Data as reference on
// the struct.
func NewView(data *Data) *view {

	brands := tview.NewList().ShowSecondaryText(false)
	brands.SetTitle("Merken").SetBorder(true)

	models := tview.NewList().ShowSecondaryText(false)
	models.SetTitle("Modellen").SetBorder(true)

	options := tview.NewList().ShowSecondaryText(false)
	options.SetTitle("Opties").SetBorder(true)

	scheduler := tview.NewForm()
	scheduler.AddInputField("Start datum (dd-mm-yyyy)", "", 20, nil, nil)
	scheduler.AddInputField("Eind datum (dd-mm-yyyy)", "", 20, nil, nil)
	scheduler.AddButton("Verder", nil)
	scheduler.AddButton("Terug", nil)
	scheduler.SetTitle("Reservering informatie").SetBorder(true)

	price := tview.NewTable()
	price.SetBorders(true)
	price.SetCellSimple(0, 0, "Omschrijving")
	price.SetCellSimple(0, 1, "Prijs per dag")
	price.SetCellSimple(0, 2, "Totaal")
	price.SetCellSimple(0, 3, "Btw")
	price.SetBorderPadding(2, 2, 5, 5)
	price.SetTitle("Prijs").SetBorder(true)

	return &view{data, tview.NewApplication(), brands, make([]*Brand, 0), models,
		make([]*Model, 0), options, make([]*Option, 0), scheduler, price, "", ""}
}

// Open will start the terminal GUI and visualise it to the screen.
func (v *view) Open() {
	// Init brand view
	v.openBrands()

	// Build our flex panel
	coll1 := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(v.brandsView, 0, 1, true).
		AddItem(v.modelsView, 0, 1, false)

	coll2 := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(v.optionsView, 0, 1, false).
		AddItem(v.schedulerView, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(coll1, 0, 1, true).
		AddItem(coll2, 0, 1, false).
		AddItem(v.priceView, 0, 1, false)

	// Launch the view
	pages := tview.NewPages().AddPage("Fietsen", flex, true, true)
	v.app.SetRoot(pages, true)
	v.app.Run()
}

// openBrands opens the brand list panel. It takes all brands from the provided data and visualises it in the GUI.
func (v *view) openBrands() {
	i := 0
	v.brands = make([]*Brand, len(v.data.Brands))
	for _, brand := range v.data.Brands {
		v.brands[i] = brand
		v.brandsView.AddItem(brand.Name, "", rune(i+49), func() {
			v.openModels(v.brands[v.brandsView.GetCurrentItem()])
		})
		i++
	}
}

// openModels opens the model list panel for the provided brand. It takes all models from the provided brand and
// visualises it in the GUI.
func (v *view) openModels(b *Brand) {
	i := 0
	v.models = make([]*Model, len(b.Models))
	for _, model := range b.Models {
		v.models[i] = model
		v.modelsView.AddItem(fmt.Sprintf("%v - %v euro per dag", model.Name, model.DailyPrice), "", rune(i+49), func() {
			v.openOptions(b, v.models[v.modelsView.GetCurrentItem()], nil)
		})
		i++
	}
	v.app.SetFocus(v.modelsView)
}

// openOptions opens the option list panel for the provided model. It takes all options from the provided model and
// visualises it in the GUI. In order to allow for some better usability, the function makes sure to sort all option
// names alphabetically.
func (v *view) openOptions(b *Brand, m *Model, options map[*Option]bool) {

	v.optionsView.Clear()

	if options == nil {
		options = make(map[*Option]bool, len(m.Options))
		for _, option := range m.Options {
			options[option] = false
		}
	}

	i := 0
	v.options = make([]*Option, len(m.Options))
	for _, option := range m.Options {
		v.options[i] = option
		i++
	}
	sort.Slice(v.options, func(i, j int) bool {
		return v.options[i].Name < v.options[j].Name
	})

	i = 0
	for _, option := range v.options {
		v.optionsView.AddItem(optionText(option, options), "", rune(i+49), func() {
			// Can't seem to use option here due to mutability screw-up
			o := v.options[v.optionsView.GetCurrentItem()]
			if chosen, exists := options[o]; exists && chosen {
				options[o] = false
			} else {
				options[o] = true
			}
			v.openOptions(b, m, options)
		})
		i++
	}
	v.optionsView.AddItem("Naar planning", "", rune(i+49), func() {
		selectedOptions := make([]*Option, len(options))
		i := 0
		for option := range options {
			selectedOptions[i] = option
			i++
		}
		v.openScheduler(b, m, selectedOptions)
	})
	v.optionsView.AddItem("Geen opties", "", rune(i+50), func() {
		v.openScheduler(b, m, nil)
	})
	v.app.SetFocus(v.optionsView)
}

// optionText is used to pass an option as well as the selected option map. It will return a string depending on whether
// or not the provided option is part of the provided selected option map. If it is, a string is returned which can be
// visually represented to the user to indicate that the option was selected. If not, the base string is returned which
// contains the option name and the daily price.
func optionText(o *Option, options map[*Option]bool) string {
	if options == nil {
		return fmt.Sprintf("%v - %v euro per dag", o.Name, o.DailyPrice)
	}
	if chosen, exists := options[o]; exists && chosen {
		return fmt.Sprintf("%v - %v euro per dag (Geselecteerd)", o.Name, o.DailyPrice)
	}
	return fmt.Sprintf("%v - %v euro per dag", o.Name, o.DailyPrice)
}

// openScheduler opens the scheduler GUI. This GUI asks for two dates; a start and end date. Input is validated by
// trying to parse the input to a time struct. It then checks if the dates are logical. If either check fails, the input
// fields are marked as red to indicate a validation issue.
func (v *view) openScheduler(b *Brand, m *Model, o []*Option) {
	v.schedulerView.GetFormItem(0).(*tview.InputField).SetChangedFunc(func(text string) {
		v.startdate = text
	})
	v.schedulerView.GetFormItem(1).(*tview.InputField).SetChangedFunc(func(text string) {
		v.enddate = text
	})
	v.schedulerView.GetButton(0).SetSelectedFunc(func() {
		valid := true
		start, err := time.Parse("02-01-2006", v.startdate)
		if err != nil {
			v.schedulerView.GetFormItem(0).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			valid = false
		}
		end, err := time.Parse("02-01-2006", v.enddate)
		if err != nil {
			v.schedulerView.GetFormItem(1).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			valid = false
		}
		if start.After(end) {
			v.schedulerView.GetFormItem(0).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			v.schedulerView.GetFormItem(1).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			valid = false
		}
		if !valid {
			// Literally cannot set the color for individual fields, looks like a bug?
			//v.schedulerView.GetFormItem(0).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			//v.schedulerView.GetFormItem(1).(*tview.InputField).SetFieldBackgroundColor(tcell.ColorRed)
			v.schedulerView.SetFieldBackgroundColor(tcell.ColorRed)
			return
		}
		v.schedulerView.SetFieldBackgroundColor(tcell.ColorBlue)
		days := end.Sub(start).Hours()/24 + 1
		v.openPrice(b, m, o, int(days))
	})
	v.schedulerView.GetButton(1).SetSelectedFunc(func() {
		v.openOptions(b, m, nil)
	})

	v.app.SetFocus(v.schedulerView)
}

// openPrice will open the GUI element that builds a table which represents the price for the full order. It lists the
// brand, the model and the chosen options individually, per day and as a total.
func (v *view) openPrice(b *Brand, m *Model, o []*Option, days int) {
	v.priceView.SetCellSimple(1, 0, fmt.Sprintf("%v - %v", b.Name, m.Name))
	v.priceView.SetCellSimple(1, 1, fmt.Sprintf("%v", m.DailyPrice))
	v.priceView.SetCellSimple(1, 2, fmt.Sprintf("%v", m.DailyPrice*days))
	v.priceView.SetCellSimple(1, 3, "21%")

	dailyPrice := m.DailyPrice

	for i, option := range o {
		v.priceView.SetCellSimple(2+i, 0, option.Name)
		v.priceView.SetCellSimple(2+i, 1, fmt.Sprintf("%v", option.DailyPrice))
		v.priceView.SetCellSimple(2+i, 2, fmt.Sprintf("%v", option.DailyPrice*days))
		v.priceView.SetCellSimple(2+i, 3, "21%")
		dailyPrice = dailyPrice + option.DailyPrice
	}

	v.priceView.SetCellSimple(len(o)+3, 0, "Totaal")
	v.priceView.SetCellSimple(len(o)+3, 1, fmt.Sprintf("%v", dailyPrice))
	v.priceView.SetCellSimple(len(o)+3, 2, fmt.Sprintf("%v", dailyPrice*days))
	v.priceView.SetCellSimple(len(o)+3, 3, "21%")

	v.app.SetFocus(v.priceView)
}
