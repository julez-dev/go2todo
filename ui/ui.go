// Package ui contains ui elements
package ui

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/julez-dev/go2todo/repo/lists"
	"github.com/julez-dev/go2todo/repo/tasks"
	"github.com/julez-dev/go2todo/service"
)

type mode int

const (
	inputMode mode = 0
	viewMode  mode = 1
)

type page int

const (
	viewListsPage page = 0
	viewTasksPage page = 1
)

type model struct {
	storage *service.Storage
	mode    mode
	page    page

	currentError error

	cursorLists int
	cursorTasks int

	textInput textinput.Model

	lists   []*lists.List
	newList *lists.List

	tasks   []*tasks.Task
	newTask *tasks.Task
}

// Event responses

type getListsResponse struct {
	lists []*lists.List
}

type deleteListResponse struct{}

type errorResponse struct {
	err error
}

func (e *errorResponse) Error() string {
	return e.err.Error()
}

type createListResponse struct{}

type getTasksResponse struct{ tasks []*tasks.Task }

type createTaskResponse struct{}

type delteTaskResponse struct{}

type updateTaskResponse struct {
	task *tasks.Task
}

// Messages

func (m *model) deleteList() tea.Msg {
	if len(m.lists) > 0 && m.cursorLists <= len(m.lists) {
		item := m.lists[m.cursorLists]

		err := m.storage.DeleteList(context.Background(), item.ID)

		if err != nil {
			return &errorResponse{err: err}
		}

		return &deleteListResponse{}
	}

	return nil
}

func (m *model) getLists() tea.Msg {
	lists, err := m.storage.GetLists(context.TODO())

	if err != nil {
		return &errorResponse{err: err}
	}

	return &getListsResponse{
		lists: lists,
	}
}

func (m *model) createList() tea.Msg {
	_, err := m.storage.StoreList(context.Background(), m.newList)

	if err != nil {
		return &errorResponse{err: err}
	}

	return &createListResponse{}
}

func (m *model) getTasks() tea.Msg {
	if len(m.lists) > 0 && m.cursorLists <= len(m.lists) {
		tasks, err := m.storage.GetTasks(context.TODO(), m.lists[m.cursorLists].ID)

		if err != nil {
			return &errorResponse{err: err}
		}

		return &getTasksResponse{tasks: tasks}
	}

	return nil
}

func (m *model) createTask() tea.Msg {
	_, err := m.storage.StoreTask(context.Background(), m.newTask)

	if err != nil {
		return &errorResponse{err: err}
	}

	return &createTaskResponse{}
}

func (m *model) deleteTask() tea.Msg {
	if len(m.tasks) > 0 && m.cursorTasks <= len(m.tasks) {
		err := m.storage.DeleteTask(context.Background(), m.tasks[m.cursorTasks].ID)

		if err != nil {
			return &errorResponse{err: err}
		}

		return &delteTaskResponse{}
	}

	return nil
}

func (m *model) updateTask() tea.Msg {
	task := m.tasks[m.cursorTasks]
	task.Completed = !task.Completed
	task, err := m.storage.UpdateTask(context.Background(), task)

	if err != nil {
		return &errorResponse{err: err}
	}

	return &updateTaskResponse{task: task}
}

func New(storage *service.Storage) *model {
	ti := textinput.NewModel()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return &model{
		storage:   storage,
		textInput: ti,
		mode:      viewMode,
		page:      viewListsPage,
	}
}

func (m *model) Init() tea.Cmd {
	return m.getLists
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.currentError = nil

	switch msg := msg.(type) {
	case *getListsResponse:
		m.lists = msg.lists
		return m, nil

	case *deleteListResponse:
		if m.cursorLists > 0 {
			m.cursorLists--
		}
		return m, m.getLists

	case *errorResponse:
		m.currentError = msg
		return m, nil

	case *createListResponse:
		return m, m.getLists

	case *createTaskResponse:
		return m, m.getTasks

	case *getTasksResponse:
		m.tasks = msg.tasks
		return m, nil

	case *delteTaskResponse:
		if m.cursorTasks > 0 {
			m.cursorTasks--
		}
		return m, m.getTasks

	case *updateTaskResponse:
		m.tasks[m.cursorTasks] = msg.task
		return m, nil

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.mode != viewMode {
				break
			}

			if m.page == viewListsPage {
				if m.cursorLists > 0 {
					m.cursorLists--
				}
			}

			if m.page == viewTasksPage {
				if m.cursorTasks > 0 {
					m.cursorTasks--
				}
			}

		case "down", "j":
			if m.mode != viewMode {
				break
			}

			if m.page == viewListsPage && m.cursorLists < len(m.lists)-1 {
				m.cursorLists++
			}

			if m.page == viewTasksPage && m.cursorTasks < len(m.tasks)-1 {
				m.cursorTasks++
			}

		case "enter":
			if m.mode == inputMode {
				if m.page == viewListsPage {
					m.newList = &lists.List{
						Name: m.textInput.Value(),
					}

					m.textInput.Reset()
					m.mode = viewMode

					return m, m.createList
				}

				if m.page == viewTasksPage {
					m.newTask = &tasks.Task{
						Text:   m.textInput.Value(),
						ListID: m.lists[m.cursorLists].ID,
					}

					m.textInput.Reset()
					m.mode = viewMode

					return m, m.createTask
				}
			}

			if m.page != viewTasksPage && m.mode == viewMode && len(m.lists) > 0 {
				m.page = viewTasksPage
				m.cursorTasks = 0
				return m, m.getTasks
			}

		case tea.KeyEsc.String():
			if m.mode == inputMode {
				m.textInput.Reset()
				m.mode = viewMode

				return m, nil
			}

		case "i":
			if m.mode != inputMode {
				m.mode = inputMode

				if m.page == viewListsPage {
					m.textInput.Placeholder = "New list name"
				}

				if m.page == viewTasksPage {
					m.textInput.Placeholder = "New task name"
				}

				return m, nil
			}

		case tea.KeyDelete.String(), "d":
			if m.page == viewListsPage && m.mode == viewMode {
				return m, m.deleteList
			}

			if m.page == viewTasksPage && m.mode == viewMode {
				return m, m.deleteTask
			}

		case tea.KeyRight.String(), "l":
			if m.page != viewTasksPage && m.mode == viewMode && len(m.lists) > 0 {
				m.page = viewTasksPage
				m.cursorTasks = 0
				return m, m.getTasks
			}

		case tea.KeyLeft.String(), tea.KeyBackspace.String(), "h":
			if m.page == viewTasksPage && m.mode == viewMode {
				m.page = viewListsPage
				return m, m.getLists
			}

		case " ":
			if m.page == viewTasksPage && m.mode == viewMode {
				return m, m.updateTask
			}
		}
	}

	if m.mode == inputMode {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	s := &strings.Builder{}

	if m.currentError != nil {
		s.WriteString("Current error: " + m.currentError.Error() + "\n\n")
	}

	if m.page == viewListsPage {
		longest := 0
		for _, listItem := range m.lists {
			length := utf8.RuneCountInString(listItem.Name)
			if length > longest {
				longest = length
			}
		}

		for i, listItem := range m.lists {
			cursor := " "
			if m.cursorLists == i {
				cursor = ">"
			}

			color.New(color.FgHiGreen).Fprint(s, cursor)
			s.WriteString(fmt.Sprintf(" %-"+fmt.Sprint(longest)+"s\n", listItem.Name))
		}
	}

	if m.page == viewTasksPage {
		list := m.lists[m.cursorLists]

		s.WriteString("  Tasks for " + list.Name + "\n\n")

		longest := 0
		for _, taskItem := range m.tasks {
			length := utf8.RuneCountInString(taskItem.Text)
			if length > longest {
				longest = length
			}
		}

		for i, taskItem := range m.tasks {
			cursor := " "
			if m.cursorTasks == i {
				cursor = ">"
			}

			color.New(color.FgHiGreen).Fprint(s, cursor)
			s.WriteString(fmt.Sprintf(" %-"+fmt.Sprint(longest)+"s [", taskItem.Text))
			if taskItem.Completed {
				color.New(color.FgHiGreen).Fprint(s, "X")
			} else {
				color.New(color.FgHiRed).Fprint(s, "-")
			}
			s.WriteString("]\n")
		}
	}

	if m.mode == inputMode {
		s.WriteString("\n" + m.textInput.View())
	}

	return s.String()
}
