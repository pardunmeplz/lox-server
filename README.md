# **LoxServer**  

LoxServer is a **hobby implementation** of a **Language Server Protocol (LSP)** for the **Lox programming language**, written in **Go**. The project aims to provide LSP support for diagnostics, formatting, autocompletion, and navigation in lox.  

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
✅ **Basic LSP communication** (via stdin/stdout)  
✅ **Handles `initialize` and `shutdown` requests**  

## **🛠️ Planned Features (TODO)**  
- [x] **Lexical Analysis** – Implement a scanner for Lox.  
- [ ] **AST Parser** – Build a parser to support syntax-aware features.  
    - [x] **Parsing Tokens** - Parse all the lox tokens to a valid AST
    - [x] **Resolution Analysis** - Check for scope issues and resolve variables
    - [ ] **Type Checks** - Do statically viable type checks
    - [ ] **Panic - Recover** - Ignore errors caused by a preceding error to avoid unnecessary error reporting 
- [x] **Diagnostics (`textDocument/publishDiagnostics`)** – Show syntax errors in real-time.  
- [ ] **Hover (`textDocument/hover`)** – Provide information on symbols.  
- [x] **Go-to Definition (`textDocument/definition`)** – Jump to symbol definitions.  
- [x] **References (`textDocument/references`)** – Jump to symbol references.
- [ ] **Formatting (`textDocument/formatting`)** - Auto format code
- [ ] **Auto-Completion (`textDocument/completion`)** – Suggest keywords and variables.  

## **📖 Resources & References**  
- [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/)  
- [Crafting Interpreters](https://craftinginterpreters.com/) – Lox Language Reference  
- [Go Language Documentation](https://go.dev/doc/)  

---

### **💡 Why This Matters?**  
It doesn't, This project is a **learning experience** in both **LSP development** and **Go programming**.
This is both my first LSP and my first project writing go.
It serves as an exploration of how to build a structured Go project from scratch.  

