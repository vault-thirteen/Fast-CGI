<style>
.l1{font-size:2.5em; color: orange; background-color: #EEEEEE; text-align: center; }
.l2{font-size:1.5em; color: grey; background-color: #FFECB3; text-align: center; }
p{font-size: 1.2em; color: black; }
.bq {background: #EEEEEE; border-left: 10px solid #FFECB3; margin: 1.5em 10px; padding: 0.5em 10px; }

</style>
<table>
<tr><td class="l1"><b>Fast CGI</b> Compatibility Layer</td></tr>
<tr><td class="l2"><i>For Go Programming Language</i></td></tr>
</table>

<p>This library provides objects, methods and functions to work with FastCGI 
interface in Go programming language.</p>

The library implements methods for a FastCGI <b>client</b>, mostly.

Implementation of a FastCGI <b>server</b> can be found in a built-in standard library 
of Go language – `net/http/fcgi` package, https://pkg.go.dev/net/http/fcgi.

## Usage

Usage examples can be found in the ['example'](example) folder.

A test application performing an execution of a simple PHP script is here:  
[cmd/RunSimplePhpScript](cmd/RunSimplePhpScript)

The most simple usage example requires only a few lines of code and is as 
simple as the following code:
```go
TODO
```

For more complex tasks, the `Client` object and its methods can be used.

## Why ?

<b>Reason 1</b> 

<p class="bq"><b>FastCGI</b> interface provides a cross-platform and 
cross-language interface for communication between totally different programming 
language. For example, with the help of <b>FastCGI</b> it is possible to call 
a function of a PHP language from the code written in Go language, where Go is 
a compiled language and PHP is an interpreted language with no runtime 
"engine".</p>

<b>Reason 2</b>

<p class="bq"><b>FastCGI</b> interface is still actual for simple tasks and 
applications where client code is not a bottleneck, e.g. for applications using 
a lot of complicated SQL queries, such as bulletin boards and online forums.</p>

<b>Reason 3</b>

<p class="bq">There are a lot of programming languages which still support the 
<b>FastCGI</b> interface. And many of these languages are mature and 
are sufficient for their purposes.</p>

<b>Reason 4</b>

<p class="bq">This library allows to bring a new life into old projects 
written in other programming languages that support the <b>FastCGI</b> 
interface. It is not necessary to re-write PHP projects in Go from scratch. 
This is why it is said in this repository, that this project is a compatibility 
layer.</p>

## Feedback
If you have any feedback, you are free to direct it to this GitHub repository:
* Comments should be written to the `Discussions` section: 
[here](https://github.com/vault-thirteen/Fast-CGI/discussions)
* Bug reports should be directed to the `Issues` section:
[here](https://github.com/vault-thirteen/Fast-CGI/issues)
