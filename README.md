# Jotted Idea
A plaintext note file that lets you jot ideas down quickly.

## File Format (.ji)
### Head
The `.ji` file has a header at the top of the file starting and ending with
`~`s on their own line (any number you like). 4 required elements are `title`,
`description`, `section`, and `tags`. Other elements can be listed in the header,
they will be included in the `<head>` of any generated html file, but not listed
in the `<body>` of the file (for now).

```
~~~
title: Example Jot
description: Just need to show some stuff offs
section: example/idea
tags: example, notes, quick access
extra: personal meta info
~~~
```

### Body
The body is a very small subset of markdown.
It allows for ordered and unordered lists, code blocks or inline, blockquotes, and 
text that is bold, italic, or underlined as well as links or embeded images.

```
[example](https://example.com)
![alt text](example.png)

- normal list
- some other item
    1. indent list
    2. second item
- last item

here's *just* a _normal_ paragraph with **some** fancy ***formats***

> This is a block quote
```