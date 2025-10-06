# **LoxServer**  

LoxServer is a **Language Server Protocol (LSP)** for the **Lox programming language**, written in **GoLang**. The project aims to provide LSP support for diagnostics, formatting, autocompletion, and navigation in lox.  

## **🚀 Getting Started**  

### **1. Clone the Repository**  
```sh
git clone https://github.com/yourusername/lox-server.git
cd lox-server
```

### **2. Run the LSP Server**  
```sh
go run cmd/lsp/main.go
```

## **📌 Current Features**  
- [x] **Basic LSP communication** (via stdin/stdout)  
- [x] **Handles `initialize` and `shutdown` requests**  
- [x] **Lexical Analysis** – Implement a scanner for Lox.  
- [x] **AST Parser** – Build a parser to support syntax-aware features.  
    - [x] **Parsing Tokens** - Parse all the lox tokens to a valid AST
    - [x] **Resolution Analysis** - Check for scope issues and resolve variables
    - [x] **Panic - Recover** - Ignore errors caused by a preceding error to avoid unnecessary error reporting 
- [x] **Diagnostics (`textDocument/publishDiagnostics`)** – Show syntax errors in real-time.  
- [x] **Go-to Definition (`textDocument/definition`)** – Jump to symbol definitions.  
- [x] **References (`textDocument/references`)** – Jump to symbol references.
- [x] **Formatting (`textDocument/formatting`)** - Auto format code
- [x] **Auto-Completion (`textDocument/completion`)** – Suggest keywords and variables.  
- [x] **Semantic-Highlighting (`textDocument/SemanticTokens`)** - code highlighting

### ** Limitations**
- Diagnostics don't highlight the characters associated with the error, only the line number

## **📖 Resources & References**  
- [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/)  
- [Crafting Interpreters](https://craftinginterpreters.com/) – Lox Language Reference  

---

### **💡 Motivation**  
Tooling support for a language or framework is a major make or break in modern software development.
To both understand the inner workings and efforts to create and maintain tooling that are crucial to developing software is the
primary purpose of this project.
