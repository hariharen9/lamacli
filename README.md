# LamaCLI 🦙✨

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![GitHub Stars](https://img.shields.io/github/stars/hariharen9/lamacli?style=social)

## 🚀 Your Local LLM Assistant, Right in Your Terminal!

LamaCLI is a powerful and intuitive command-line interface (CLI) tool that brings the magic of Large Language Models (LLMs) directly to your terminal, powered by [Ollama](https://ollama.ai/). 

Engage with your AI assistant, browse files, and share context seamlessly—all without leaving your terminal.

## ✨ Features

*   **Interactive Chat:** Real-time, beautiful chat experience with streaming responses.
*   **Markdown Support:** Enjoy fully rendered markdown, including syntax-highlighted code.
*   **Contextual File Integration:** Easily inject file content into your prompts (`@` command), browse files, and preview them.
*   **Ollama Model Selection:** Seamlessly switch between different Ollama models.
*   **Code Copying:** Quickly copy code snippets from AI responses.
*   **Intuitive Key Bindings:** Effortless navigation and interaction.
*   **Beautiful Theming:** A modern, eye-pleasing color palette.

## ⚡️ Get Started

### Prerequisites

Before you begin, ensure you have [Ollama](https://ollama.com/download) installed and running on your system. You can install Ollama using one of the following methods:

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

LamaCLI is an open-source project. Your support helps keep this project alive and thriving!
<p align="center">
    <a href="https://www.buymeacoffee.com/hariharen">
        <img src="https://img.shields.io/badge/Buy%20Me%20a%20Coffee-FFDD00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Buy Me a Coffee" />
    </a>
    <a href="https://paypal.me/hariharen9">
        <img src="https://img.shields.io/badge/Donate-PayPal-00457C?style=for-the-badge&logo=paypal&logoColor=white" alt="PayPal" />
    </a>
</p>

## 🙏 Credits

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss), [Glamour](https://github.com/charmbracelet/glamour), [Huh](https://github.com/charmbracelet/huh), and powered by [Ollama](https://ollama.ai/).