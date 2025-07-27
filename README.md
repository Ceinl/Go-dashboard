# Go-dashboard

A simple TUI dashboard where you can create workspaces, projects, and configure modules as needed. Each workspace's modules can be configured independently, allowing you to select what you need to see in each one. Each project within a workspace has its own data for the modules, so you can keep data separated between projects.

## Getting Started: A Walkthrough

This guide will walk you through the basic workflow of Go-dashboard.

**1. Run the application:**

Open your terminal and run the following command:
```bash
go run main.go
```

**2. Create a Workspace:**

Workspaces are the top-level containers for your projects. If you don't have any, create one now.
- Press `:` to enter command mode.
- Type `neww` and press Enter.
- Give your workspace a name (e.g., "Hackathon Prep") and press Enter.

**3. Select Your Active Workspace:**

**This is an important step.** Projects are created inside the currently active workspace. Before creating a project, you must select the workspace you want to add it to.
- Enter command mode (`:`), type `swapw`, and press Enter.
- A list of your workspaces will appear. Use the arrow keys to select one and press Enter to make it active.

**4. Create a Project:**

Now that you have an active workspace, you can create a project within it.
- Enter command mode (`:`).
- Type `newp` and press Enter.
- Name your project (e.g., "Dashboard Feature") and press Enter.

**5. Configure Modules:**

Select the tools you want to use in the current workspace.
- Enter command mode (`:`), type `modules`, and press Enter.
- Use the arrow keys to navigate and the spacebar to select/deselect modules.
- Press Enter to save your selection.

**6. Navigate Your Dashboard:**

- Use `Shift+Right` and `Shift+Left` to switch between projects.
- Use `Shift+Up` and `Shift+Down` to cycle through the active modules.

All your data is saved automatically as you work.

## Features

- **Link Saver**: A bookmarking module.
- **Kanban Board**: A task management module.
- **Twitter Drafts**: A module for drafting tweets.

## Commands and Keybindings

| Keybinding        | Action                               |
| ----------------- | ------------------------------------ |
| `Shift+Up`        | Switch to the previous module        |
| `Shift+Down`      | Switch to the next module            |
| `Shift+Left`      | Switch to the previous project       |
| `Shift+Right`     | Switch to the next project           |
| `:q`              | Quit the application                 |

You can also use commands by pressing `:`:

- `:neww`: Create a new workspace.
- `:swapw`: Swap the active workspace.
- `:delw`: Delete a workspace.
- `:newp`: Create a new project in the current workspace.
- `:delp`: Delete the current project.
- `:modules`: Select modules for the current workspace.
- `:help`: Open the help view.

## Installation

To install the necessary dependencies, run the following command:

```bash
go mod tidy
```

## Built With

- [Go](https://go.dev/)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [SQLite](https://www.sqlite.org/)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
