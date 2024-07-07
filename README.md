# waiig
Write an interpreter in go

Following along with [interpreterbook.com](https://interpreterbook.com)

Commits at various points:
- Chapter 1: https://github.com/xpcoffee/waiig/tree/e8505dea44d09ea1c7588ff3b45b3727b2a74ac2

## REPL

Currently only outputs tokens.

```bash
go run .
```
```text
Hello, rick! This is the Monkey programming language!
Feel free to type commands
>>let hello = 5;
{Type:LET Literal:let}
{Type:IDENT Literal:hello}
{Type:= Literal:=}
{Type:INT Literal:5}
{Type:; Literal:;}
>>5 + 3;
```
