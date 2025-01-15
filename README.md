# Kubectl-GPT

![Latest GitHub release](https://img.shields.io/github/release/devinjeon/kubectl-gpt.svg)
![GitHub stars](https://img.shields.io/github/stars/devinjeon/kubectl-gpt.svg?label=github%20stars)

Kubectl-GPT is a kubectl plugin to generate `kubectl` commands from natural language input by using GPT model.

| ❗[WARNING] Please verify the generated commands before executing them on your k8s cluster, especially `update` and `patch` commands, as GPT-generated commands may be inaccurate. |
|-|

![demo](demo.gif)

## Installation

### Homebrew

```bash
# Install Homebrew: https://brew.sh/
brew tap devinjeon/kubectl-gpt https://github.com/devinjeon/kubectl-gpt
brew install kubectl-gpt
```

### Krew

```bash
# Install Krew: https://krew.sigs.k8s.io/docs/user-guide/setup/install/
kubectl krew index add devinjeon https://github.com/devinjeon/kubectl-gpt
kubectl krew install devinjeon/gpt
```

## Usage

Run the command line tool with your natural language input to generate a `kubectl` command.

```bash
kubectl gpt "<WHAT-YOU-WANT-TO-DO>"
```

Commands generated by the GPT model may not always be perfect, please verify them before execution.

### Prerequisite

Before you start, make sure to set your OpenAI API key as an environment variable named `OPENAI_API_KEY`.
You can get a key for using the OpenAI API at https://platform.openai.com/account/api-keys

You can add the following line to your `.zshrc` or `.bashrc` file:

```bash
export OPENAI_API_KEY=<your-key>
```

### Examples

It depends on the languages supported by the OpenAI GPT API.

```bash
# English
kubectl gpt "Print the creation time and pod name of all pods in all namespaces."
kubectl gpt "Print the memory limit and request of all pods"
kubectl gpt "Increase the replica count of the coredns deployment to 2"
kubectl gpt "Switch context to the kube-system namespace"

# Korean
kubectl gpt "현재 namespace에서 각 pod 별 생성시간 출력"
kubectl gpt "coredns deployment의 replica를 2로 증가"
```

### Options

- `--yes`: Execute the generated command without asking for confirmation
- `--no`: Print the generated command without executing it
- `--help`: Show usage information
- `--version`: Show the version

you can set the following environment variables for OpenAI API configurations:

- `OPENAI_API_URL`: OpenAI API URL (default is https://api.openai.com/v1/chat/completions)
- `OPENAI_MODEL`: The model to use for GPT-3 (default is `gpt-3.5-turbo`)
- `OPENAI_TEMPERATURE`: Controls the randomness of GPT-3's output (default is `0.2`)
- `OPENAI_MAX_TOKENS`: Maximum number of tokens for GPT-3 to generate (default is `300`)
