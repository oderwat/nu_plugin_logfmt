# Nu Plugin Logfmt

## Overview

The `nu-plugin-logfmt` is a Nushell plugin that allows conversion between logfmt (a key=value format used by some logging systems) and Nushell
values. The plugin provides two commands:

1. **from logfmt**: Converts a logfmt string to a Nushell record.
2. **to logfmt**: Converts a Nushell record to a logfmt string.

### Steps to Install

Ensure you have Go installed on your system. You can download it from [here](https://golang.org/dl/).

   ```bash
   go install github.com/oderwat/nu_plugin_logfmt@latest
   plugin add ~/go/bin/nu_plugin_logfmt
   plugin use logfmt
   ```

or move it to your nushell plugins directory if you have one for that.

## Usage

### from logfmt Command

The `from logfmt` command converts a logfmt string to a Nushell record.

If the `--typed` flag is given. The function detect simple types.

**Syntax:**

```nu
<logfmt_string> | from logfmt [--typed]
```

**Example:**

```nu
'msg="Test message" level=info esc="» Say \"Hello\""' | from logfmt
```

**Output:**

```nu
{
  "level": "info",
  "msg": "Test message",
  "esc": "» Say \"Hello\""
}
```

### to logfmt Command

The `to logfmt` command converts a Nushell record to a logfmt string.

**Syntax:**

```nu
<record> | to logfmt
```

**Example:**

```nu
{ "msg": "Hello World!", "Lang": { "Go": true, "Rust": false } } | to logfmt
```

**Output:**

```nu
msg="Hello World!" Lang.Go=true Lang.Rust=false
```

## License

Distributed under the MIT License. See `LICENSE` for more information.
