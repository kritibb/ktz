package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kritibb/ktz/tzdata"
	"io"
	"os"
	"strings"
)

// constants and variables for list view
const listHeight = 14

// style variables for list display
var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var locationData locationInfo

// variable for table view
var baseTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// decalre item type for listing
type item string

// declare viewState type which represents either table or list view
type viewState int

// declare constant listView and tableView (iota starts from 0 and increments by 1)
// in golang, when you declare const within a `const` block and only specify the
// type for the first const, the subsequent constants inherit from the first one.
const (
	listView viewState = iota
	tableView
)

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {

	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	table    table.Model
	choice   string
	quitting bool
	state    viewState
}

func initialModel(state viewState) model {
	//Initialize the list model
	const defaultWidth = 20
	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	l.Title = fmt.Sprintf("Select one timezone:")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	//initialize table model
	columns := []table.Column{}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(1),
	)

	s := table.DefaultStyles()

	s.Cell = s.Cell.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderRight(true).
		BorderLeft(true)

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderRight(true).
		BorderBottom(true).
		BorderLeft(true).
		Bold(false)

	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)
	m := model{list: l, table: t, state: state}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.state == listView {
			m.list.SetWidth(msg.Width)
			return m, nil
		}

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if m.state == listView {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
					//check if selected option is a country and it has more than one timezones
					if timezones, ok := tzdata.CountryToIanaTimezone[m.choice]; ok && len(timezones) > 1 {
						locationData.country = string(i)
						m.state = listView
						items := make([]list.Item, len(timezones))
						for idx, tz := range timezones {
							items[idx] = item(tz)
						}
						m.list.SetItems(items)
						return m, cmd
					} else {
						// remove the listview state by picking -1
						m.state = -1
						//if selected option is not a city or a country, but a country's tz in the form of "America/New_York"
						if tzCountry, ok1 := tzdata.CityToIanaTimezone[m.choice]; !ok1 {
							if tzList, ok2 := tzdata.CountryToIanaTimezone[m.choice]; !ok2 {
								locationData.timezone = m.choice
							} else {
								locationData.country = m.choice
								locationData.timezone = tzList[0]
							}
						} else {
							//if selected option is a city
							locationData.country = tzCountry["country"]
							locationData.city = m.choice
							locationData.timezone = tzCountry["tz"]
						}
						return m, tea.Quit
					}

				}
			}
		}
	}
	switch m.state {
	case listView:
		m.list, cmd = m.list.Update(msg)
	case tableView:
		m.table, cmd = m.table.Update(msg)
		cmd = tea.Quit
	}
	return m, cmd

}

func (m model) View() string {
	switch m.state {
	case listView:
		if m.quitting {
			return quitTextStyle.Render("Don't wanna check time? That’s cool.")
		}
		return "\n" + m.list.View()
	case tableView:
		m.quitting = true
		return "\n" + baseTableStyle.Render(m.table.View()) + "\n \n"
	default:
		return ""
	}

}

// listViewTz lists(bubbletea simple-list format) all possible timzones returned by prefix search for country/city
//
// Parameters:
//
//	-timezones: list of timezones
func listViewTz(timezones []string) {
	const defaultWidth = 20
	m := initialModel(listView)
	//Accumulate items in a slice
	items := []list.Item{}
	for _, tz := range timezones {
		items = append(items, item(tz))
	}
	m.list.SetItems(items)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// renderZoneInfoTable returns a table consisting timezone, datetime, zone abbreviation, if any.
//
// Parameters:
//
//	-zoneData: A ZoneInfo type containing formatted time, zone and abbreviation, if any.
func renderZoneInfoTable(zoneData zoneInfo) {
	if zoneData.abbreviation != "" {
		fmt.Printf("\n Timezone for %v:", zoneData.abbreviation)
	} else {
		fmt.Printf("\n Timezone for %v:", zoneData.timezoneName)
	}
	m := initialModel(tableView)
	rows := []table.Row{}
	columns := []table.Column{}
	columns = append(columns,
		table.Column{Title: "TimeZone", Width: 20},
		table.Column{Title: "Date/Time", Width: 30},
	)
	if zoneData.abbreviation != "" {
		columns = append(columns,
			table.Column{Title: "Zone Abbr.", Width: 20},
		)
		rows = append(rows, table.Row{zoneData.timezoneName, zoneData.formattedTime, zoneData.abbreviation})
	} else {
		rows = append(rows, table.Row{zoneData.timezoneName, zoneData.formattedTime})
	}
	if rows != nil && columns != nil {
		m.table.SetColumns(columns)
		m.table.SetRows(rows)
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}
}

// renderDateTimeTableFromLocation returns a table consisting timezone, country, datetime.
//
// Parameters:
//
//	-locationList: A list of country or a city
func renderDateTimeTableFromLocation(currentLocationData locationInfo) {
	if currentLocationData.city == "" {
		fmt.Printf("\n Timezone for %v:", currentLocationData.country)
	} else {
		fmt.Printf("\n Timezone for %v:", currentLocationData.city)
	}
	m := initialModel(tableView)
	row := []table.Row{}
	columns := []table.Column{}
	columns = append(columns,
		table.Column{Title: "TimeZone", Width: 20},
		table.Column{Title: "Country", Width: 25},
		table.Column{Title: "Date/Time", Width: 30},
	)
	row = append(row, table.Row{currentLocationData.timezone, currentLocationData.country, currentLocationData.formattedTime})
	if row != nil && columns != nil {
		m.table.SetColumns(columns)
		m.table.SetRows(row)
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}

}
