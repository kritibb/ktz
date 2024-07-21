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

//declare constant listView and tableView (iota starts from 0 and increments by 1)
//in golang, when you declare const within a `const` block and only specify the
//type for the first const, the subsequent constants inherit from the first one.

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

func initialListModel() model {
	//Initialize the list model
	const defaultWidth = 20
	l := list.New([]list.Item{}, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select one timezone:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	m := model{list: l, state: listView}
	return m
}

func initialTableModel() model {

	//Initalize the table model
	columns := []table.Column{
		{Title: "City", Width: 10},
		{Title: "TimeZone", Width: 20},
		{Title: "Country", Width: 20},
		{Title: "Date/Time", Width: 30},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = lipgloss.NewStyle()
	t.SetStyles(s)

	m := model{table: t, state: tableView}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case listView:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m, nil

		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
				}
				return m, tea.Quit
			}
		}

		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case tableView:
		if m.quitting {
			return m, tea.Quit
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				return m, tea.Quit
			}
		}
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	var viewResponse string
	switch m.state {
	case listView:
		if m.choice != "" {
			tableViewDateTime([]string{m.choice})
			// m.quitting = true
			// s, _ := formatTime(cityToTimezone[m.choice]["tz"])
			// return "\n"+s+"\n"
			return "\n"
		}
		if m.quitting {
			return quitTextStyle.Render("Don't wanna check time? That’s cool.")
		}
		viewResponse = "\n" + m.list.View()
	case tableView:
		viewResponse = baseTableStyle.Render(m.table.View()) + "\n"
	}
	return viewResponse

}

// listViewTz lists(bubbletea simple-list format) all possible timzones returned by prefix search (eg: "Asia/Kathmandu")
//
// Parameters:
//
//	-timezones: list of timezones
func listViewTz(timezones []string) {
	const defaultWidth = 20
	m := initialListModel()
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

	m := initialTableModel()
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
