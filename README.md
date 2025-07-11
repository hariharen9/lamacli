# LamaCLI 🦙✨


![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![GitHub Stars](https://img.shields.io/github/stars/hariharen9/lamacli?style=social)

## 🚀 Unleash the Power of Local LLMs, Right in Your Terminal!

LamaCLI is a sleek, intuitive, and powerful command-line interface (CLI) tool that brings the magic of Large Language Models (LLMs) directly to your terminal, powered by [Ollama](https://ollama.ai/). Forget context switching and privacy concerns; with LamaCLI, your AI assistant is always just a `ctrl+c` away (well, two `ctrl+c`s for safety!).

Designed for developers, writers, and anyone who craves efficiency, LamaCLI integrates seamlessly into your workflow, offering a rich, interactive chat experience, file browsing, and instant context sharing—all without ever leaving your favorite terminal.

## ✨ Features

*   **Interactive Chat Interface:** Engage with your LLM in a beautiful, real-time chat environment.
*   **Streaming Responses:** Watch as your AI assistant types, providing a dynamic and engaging experience.
*   **Markdown Rendering:** Enjoy fully rendered markdown responses, including syntax-highlighted code blocks, making complex information easy to digest.
*   **Contextual File Integration:**
    *   **`@` Command:** Easily inject the content of any local file directly into your chat prompt for instant, context-aware AI interactions.
    *   **Integrated File Explorer:** Browse your project files directly within the application.
    *   **File Viewer:** Preview file contents without leaving the CLI.
*   **Ollama Model Selection:** Seamlessly switch between different Ollama models to suit your task.
*   **Code Block Copying:** Quickly copy code snippets from AI responses to your clipboard with a single keypress.
*   **Safe Exit Confirmation:** Prevent accidental exits with a double `Ctrl+C` confirmation.
*   **Intuitive Key Bindings:** Navigate and interact with the application effortlessly.
*   **Beautiful Theming:** A delightful user experience with a modern, eye-pleasing color palette.

## ⚡️ Get Started in Seconds!

Ready to supercharge your terminal?

### Prerequisites

Before you begin, ensure you have [Ollama](https://ollama.ai/download) installed and running on your system. You can install Ollama using one of the following methods:

**macOS**
- Download the app from [Ollama's website](https://ollama.com/download).
- Or install via Homebrew:
    ```bash
    brew install ollama
    ```

**Windows**
- Download the installer from [Ollama's website](https://ollama.com/download).

**Linux**
- Run the official installation script:
    ```bash
    curl -fsSL https://ollama.com/install.sh | sh
    ```

After installation, make sure Ollama is running, then pull at least one model (e.g., `ollama pull llama3.2:3b`).

```bash
ollama pull llama3.2:3b
```

### Installation

LamaCLI is built with Go, making installation a breeze:

```bash
go install github.com/hariharen9/lamacli@latest
```

### Usage

Simply run LamaCLI from your terminal:

```bash
lamacli
```

### Key Bindings

| Key       | Description                                                               |
| :-------- | :------------------------------------------------------------------------ |
| `Enter`   | Send message (in chat), Open file/folder (in file explorer)               |
| `↑`/`↓`   | Scroll history (in chat), Navigate items (in file tree/model select)      |
| `@`       | Trigger file context selection (in chat input)                            |
| `F`       | Open File Explorer                                                        |
| `M`       | Switch AI Model                                                           |
| `R`       | Reset/Clear Chat History                                                  |
| `C`       | Copy Code Blocks (when available in chat)                                 |
| `H`       | Show detailed Help screen                                                 |
| `Backspace` | Go to parent folder (in file explorer), Back to explorer (in file viewer) |
| `Esc`     | Return to chat from any view (file explorer, model select, help)          |
| `Ctrl+C`  | Exit application (requires two presses for confirmation)                  |

## 🤝 Contributing

We welcome contributions! If you have ideas for new features, bug fixes, or improvements, please feel free to:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature`).
3.  Make your changes.
4.  Commit your changes (`git commit -m 'feat: Add new feature'`).
5.  Push to the branch (`git push origin feature/your-feature`).
6.  Open a Pull Request.

Please ensure your code adheres to the existing style and conventions.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ❤️ Support LamaCLI

LamaCLI is an open-source project built with passion. If you find this tool useful and would like to support its continued development, consider buying me a coffee or contributing via PayPal! Your support helps keep this project alive and thriving.

*   [☕ Buy Me a Coffee](https://www.buymeacoffee.com/hariharen)
*   [💰 PayPal](https://paypal.me/hariharen9)


## 🙏 Credits

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss), [Glamour](https://github.com/charmbracelet/glamour), [Huh](https://github.com/charmbracelet/huh), and powered by [Ollama](https://ollama.ai/).
