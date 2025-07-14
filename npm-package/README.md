# LamaCLI

A CLI tool for interacting with Llama models via Ollama.

## Installation

Install globally via npm:

```bash
npm install -g lamacli
```

Or use with npx without installation:

```bash
npx lamacli
```

## Usage

### Interactive Mode

Simply run the command without any arguments to start interactive mode:

```bash
lamacli
```

### CLI Mode

You can also use LamaCLI with command line arguments for scripting and automation:

```bash
lamacli [command] [options]
```

## Requirements

- [Ollama](https://ollama.ai) must be installed and running
- At least one model must be pulled in Ollama (e.g., `ollama pull llama2`)

## Platform Support

LamaCLI supports the following platforms:
- Linux (x64, ARM64)
- macOS (x64, ARM64) 
- Windows (x64)

## License

MIT

## Repository

[https://github.com/hariharen9/lamacli](https://github.com/hariharen9/lamacli)
