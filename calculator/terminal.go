package calculator

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sort"
)

type view struct {
	data          *Data
	app           *tview.Application
	brandsView    *tview.List
	brands        []*Brand
	modelsView    *tview.List
	models        []*Model
	optionsView   *tview.List
	options       []*Option
	schedulerView *tview.InputField
	priceView     *tview.Table
}

func NewView(data *Data) *view {

	brands := tview.NewList()
	brands.ShowSecondaryText(false).SetBorder(true)

	models := tview.NewList().ShowSecondaryText(false)
	models.SetTitle("Modellen").SetBorder(true)

	options := tview.NewList().ShowSecondaryText(false)
	options.SetTitle("Opties").SetBorder(true)

	scheduler := tview.NewInputField()
	scheduler.SetFieldBackgroundColor(tcell.ColorBlack)
	scheduler.SetTitle("Planning").SetBorder(true)

	price := tview.NewTable()
	price.SetTitle("Prijs").SetBorder(true)

	return &view{data, tview.NewApplication(), brands, make([]*Brand, 0), models,
		make([]*Model, 0), options, make([]*Option, 0), scheduler, price}
}

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

func (v *view) openModels(b *Brand) {
	i := 0
	v.models = make([]*Model, len(b.Models))
	for _, model := range b.Models {
		v.models[i] = model
		v.modelsView.AddItem(model.Name, "", rune(i+49), func() {
			v.openOptions(b, v.models[v.modelsView.GetCurrentItem()], nil)
		})
		i++
	}
	v.app.SetFocus(v.modelsView)
}

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

func optionText(o *Option, options map[*Option]bool) string {
	if options == nil {
		return o.Name
	}
	if chosen, exists := options[o]; exists && chosen {
		return o.Name + " (Geselecteerd)"
	}
	return o.Name
}

func (v *view) openScheduler(b *Brand, m *Model, o []*Option) {
	v.app.SetFocus(v.schedulerView)
}
