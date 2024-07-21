package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	l.Title = "Select one timezone:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	//initialize table model
	columns := []table.Column{
		{Title: "City", Width: 15},
		{Title: "TimeZone", Width: 20},
		{Title: "Country", Width: 15},
		{Title: "Date/Time", Width: 30},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
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
					m.state = tableView
					rows := []table.Row{}
					timeResponse, err := formatTime(cityToTimezone[m.choice]["tz"])
					if err != nil {
						fmt.Printf("Error loading time zone: %v\n", err)
					} else {
						rows = append(rows, table.Row{m.choice, cityToTimezone[m.choice]["tz"], cityToTimezone[m.choice]["country"], timeResponse})
					}

					m.table.SetRows(rows)
				}
				return m, tea.Quit

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

// listViewTz lists(bubbletea simple-list format) all possible timzones returned by prefix search (eg: "Asia/Kathmandu")
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

// tableViewTime returns a table consisting city, timezone, country, datetime.
//
// Parameters:
//
//	-timezones: list of timezones
func tableViewDateTime(cities []string) {

	m := initialModel(tableView)
	//Accumulate items in a slice
	rows := []table.Row{}
	for _, city := range cities {
		timeResponse, err := formatTime(cityToTimezone[city]["tz"])
		if err != nil {
			fmt.Printf("Error loading time zone: %v\n", err)
		} else {
			rows = append(rows, table.Row{city, cityToTimezone[city]["tz"], cityToTimezone[city]["country"], timeResponse})
		}
	}

	m.table.SetRows(rows)
	m.table.SetHeight(len(cities))
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
