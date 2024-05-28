package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/hhhapz/doc"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("105")).Render
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type app struct {
	debug bool
	pkg   doc.Package
	w, h  int

	sidebar   textarea.Model
	main      viewport.Model
	debugMenu viewport.Model

	inp textinput.Model

	debugContent   string
	sidebarContent []string
	content        []string
	mainActive     bool
}

func (a *app) Debug(s string, args ...any) {
	a.debugContent = fmt.Sprintf(s+"\n", args...) + a.debugContent
	a.debugMenu.SetContent(a.debugContent)
}

func newModel(pkg doc.Package, debug bool) (*app, error) {
	w, h, err := term.GetSize(os.Stdin.Fd())
	if err != nil {
		return nil, err
	}

	sidebar := textarea.New()
	sidebar.SetWidth(30)
	sidebar.SetHeight(h - 4)
	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("81")).
		Margin(0).
		Padding(0)
	sidebar.FocusedStyle = textarea.Style{
		Base:       style,
		CursorLine: lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("7")),
	}
	sidebar.BlurredStyle = textarea.Style{
		Base:       style.BorderForeground(lipgloss.Color("60")),
		CursorLine: lipgloss.NewStyle().Background(lipgloss.Color("60")),
	}
	sidebar.Prompt = ""
	sidebar.ShowLineNumbers = false
	sidebar.CharLimit = 0

	main := viewport.New(w-34, h-4)
	main.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("60"))

	glamour.DarkStyleConfig.CodeBlock.Margin = new(uint)
	glamour.LightStyleConfig.CodeBlock.Margin = new(uint)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(w-36),
	)
	if err != nil {
		return nil, err
	}

	sidebarTpl := template.New("sidebar")
	sidebarTpl.Funcs(template.FuncMap{
		"sortedRange": sortedRange,
	})
	_, err = sidebarTpl.Parse(sdTpl)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := sidebarTpl.Execute(buf, pkg); err != nil {
		return nil, err
	}
	sidebar.SetValue(buf.String())
	sidebar.CursorStart()
	sidebar.SetCursor(0)
	sidebar.Focus()
	// why
	for range sidebar.LineCount() * 2 {
		sidebar.CursorUp()
	}

	mainTpl := template.New("doc")
	mainTpl.Funcs(template.FuncMap{
		"sortedRange": sortedRange,
	})

	_, err = mainTpl.Parse(pkgTpl)
	if err != nil {
		return nil, err
	}

	if err := mainTpl.Execute(renderer, pkg); err != nil {
		return nil, err
	}

	if err := renderer.Close(); err != nil {
		return nil, err
	}

	b, _ := io.ReadAll(renderer)
	main.SetContent(string(b))

	inp := textinput.New()
	inp.CharLimit = 32

	var debugMenu viewport.Model
	if debug {
		main.Height = h - 12
		debugMenu = viewport.New(w-34, 8)
	}
	return &app{
		debug:          debug,
		pkg:            pkg,
		w:              w,
		h:              h,
		sidebar:        sidebar,
		main:           main,
		debugMenu:      debugMenu,
		inp:            inp,
		sidebarContent: sidebarItems(buf.String()),
		content:        stripTrimAnsi(string(b)),
	}, nil
}

func (a *app) Init() tea.Cmd {
	return nil
}

func (a *app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.inp.Focused() {
			switch msg.String() {
			case "esc", "enter":
				a.inp.Blur()
			case "ctrl+c":
				return a, tea.Quit
			default:
				a.inp, cmd = a.inp.Update(msg)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q", "esc":
				return a, tea.Quit
			case "enter":
				a.inp.Blur()
			case "left", "h":
				a.mainActive = false
			case "right", "l":
				a.mainActive = true
			case "/":
				a.inp.SetValue("")
				a.inp.Focus()
			default:
				switch {
				case a.mainActive:
					a.main, cmd = a.main.Update(msg)
				default:
					a.sidebar, a.main = a.UpdateTextArea(a.sidebar, a.main, msg)
				}
			}
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonLeft:
			a.mainActive = msg.X > a.sidebar.Width()
		case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown:
			switch {
			case a.debug:
				a.debugMenu, cmd = a.debugMenu.Update(cmd)
			case msg.X > a.sidebar.Width():
				a.main, cmd = a.main.Update(msg)
			default:
				a.sidebar, cmd = a.sidebar.Update(msg)
			}
		}
	}
	if a.inp.Focused() && len(a.inp.Value()) > 0 {
		input := strings.ToLower(a.inp.Value())
		rank := fuzzy.RankFind(input, a.sidebarContent)
		target := a.sidebar.Line()
		if len(rank) > 0 {
			target = rank[0].OriginalIndex
		}
		a.sidebar, a.main = a.ScrollTo(a.sidebar, a.main, target-a.sidebar.Line())
	}
	if a.mainActive {
		a.main.Style = a.main.Style.BorderForeground(lipgloss.Color("81"))
		a.sidebar.Blur()
	} else {
		a.main.Style = a.main.Style.BorderForeground(lipgloss.Color("60"))
		a.sidebar.FocusedStyle.Base = a.sidebar.FocusedStyle.Base.BorderForeground(lipgloss.Color("81"))
		a.sidebar.Focus()
	}
	return a, cmd
}

func (a *app) View() string {
	main := a.main.View()
	if a.debug {
		main = lipgloss.JoinVertical(lipgloss.Center, main, a.debugMenu.View())
	}
	content := lipgloss.JoinHorizontal(lipgloss.Center, a.sidebar.View(), main)

	layout := lipgloss.JoinVertical(lipgloss.Left, a.topView(), content, a.bottomView())
	return layout
}

func (a *app) topView() string {
	return titleStyle("  godocs: View Go package documentation")
}

func (a *app) bottomView() string {
	if !a.inp.Focused() {
		return helpStyle("  /: search • ←/↓/↑/→: Navigate  • q: Quit")
	}
	return a.inp.View()
}

func (a *app) UpdateTextArea(m textarea.Model, v viewport.Model, msg tea.Msg) (textarea.Model, viewport.Model) {
	var n int
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, v.KeyMap.PageDown):
			n = m.Height()
		case key.Matches(msg, v.KeyMap.PageUp):
			n = -m.Height()
		case key.Matches(msg, v.KeyMap.HalfPageDown):
			n = m.Height() / 2
		case key.Matches(msg, v.KeyMap.HalfPageUp):
			n = -m.Height() / 2
		case key.Matches(msg, v.KeyMap.Down):
			n = 1
		case key.Matches(msg, v.KeyMap.Up):
			n = -1
		}
	}
	m, v = a.ScrollTo(m, v, n)
	return m, v
}

func (a *app) ScrollTo(m textarea.Model, v viewport.Model, n int) (textarea.Model, viewport.Model) {
	f := m.CursorDown
	if n < 0 {
		f, n = m.CursorUp, -n
	}
	for range n {
		f()
	}
	m, _ = m.Update(tea.KeyMsg{})
	current := a.sidebarContent[m.Line()]
	if current == "Overview" {
		v.SetYOffset(0)
		return m, v
	}

	for i, line := range a.content {
		switch current {
		case "overview":
			v.SetYOffset(0)
			break
		case "constants", "variables", "functions", "types":
			if line == current {
				v.SetYOffset(i - 3)
				break
			}
		default:
			if strings.Contains(line, "#") && strings.Contains(line, current) {
				v.SetYOffset(i - 3)
				break
			}
		}
	}
	return m, v
}

func sidebarItems(raw string) []string {
	return strings.Split(strings.ToLower(strings.NewReplacer(" ", "", "•", "").Replace(raw)), "\n")
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripTrimAnsi(str string) []string {
	stripped := strings.ToLower(re.ReplaceAllString(str, ""))
	split := strings.Split(stripped, "\n")
	for i, s := range split {
		split[i] = strings.TrimSpace(s)
	}
	return split
}
